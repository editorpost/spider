package store

import (
	"bytes"
	"context"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/editorpost/donq/res"
	"log/slog"
)

const (
	PayloadFile     = "payload.json"
	VisitedFile     = "visited.json"
	HTMLSourceFile  = "index.html"
	ChunkTimeFormat = "06-01"
)

type (
	BucketStorage struct {
		folder string
		bucket string
		client *s3.Client
	}
	Storage interface {
		Save(data []byte, filename string) error
		Load(filename string) ([]byte, error)
		Delete(filename string) error
		Reset() error
	}
)

func NewS3Client(bucket res.S3) (*s3.Client, error) {

	access := credentials.NewStaticCredentialsProvider(bucket.AccessKey, bucket.SecretKey, "")

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(access),
		config.WithRegion(bucket.Region),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(bucket.EndPoint)
		o.UsePathStyle = true
	})

	return client, nil
}

func NewBucketStorage(bucket, folder string, client *s3.Client) *BucketStorage {
	return &BucketStorage{
		folder: folder,
		bucket: bucket,
		client: client,
	}
}

func NewStorage(bucket res.S3, folder string) (Storage, error) {

	if IsLocalBucket(bucket) {
		return NewFolderStorage(bucket, folder)
	}

	client, err := NewS3Client(bucket)
	if err != nil {
		return nil, fmt.Errorf("failed to create s3 client, %w", err)
	}

	return NewBucketStorage(bucket.Bucket, folder, client), nil
}

func (b *BucketStorage) path(filename string) *string {

	if len(b.folder) == 0 {
		return aws.String(filename)
	}

	return aws.String(fmt.Sprintf("%s/%s", b.folder, filename))
}

func (b *BucketStorage) Save(data []byte, filename string) error {

	uploader := manager.NewUploader(b.client)

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    b.path(filename),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("failed to upload object, %w", err)
	}

	return nil
}

func (b *BucketStorage) Load(filename string) ([]byte, error) {

	obj, err := b.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    b.path(filename),
	})

	if err != nil {
		return nil, err
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj.Body)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (b *BucketStorage) Delete(filename string) error {
	_, err := b.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    b.path(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object, %w", err)
	}
	return nil
}

func (b *BucketStorage) list() ([]types.Object, error) {

	list, err := b.client.ListObjectsV2(context.TODO(), &s3.ListObjectsV2Input{
		Bucket: aws.String(b.bucket),
		Prefix: aws.String(b.folder),
	})

	if err != nil {
		return []types.Object{}, fmt.Errorf("failed to list objects, %w", err)
	}

	return list.Contents, nil
}

// Reset objects recursively with prefix b.folder
func (b *BucketStorage) Reset() error {

	if len(b.folder) == 0 {
		return fmt.Errorf("folder is empty, use DropBucket method to explicitly drop bucket")
	}

	return CleanupBucket(b.client, b.bucket, b.folder)
}

// CleanupBucket deletes the contents of a S3 bucket
// code from: github.com/aws/aws-sdk-go-v2/feature/s3/manager@v1.17.2/internal/integration/integration.go
func CleanupBucket(client *s3.Client, bucketName, prefix string) error {
	var errs []error

	{
		slog.Info("TearDown: Deleting objects from bucket", "bucket", bucketName)
		input := &s3.ListObjectsV2Input{
			Bucket: &bucketName,
			Prefix: aws.String(prefix),
		}
		for {
			listObjectsV2, err := client.ListObjectsV2(context.Background(), input)
			if err != nil {
				return fmt.Errorf("failed to list objects, %w", err)
			}

			var del types.Delete
			for _, content := range listObjectsV2.Contents {
				obj := content
				del.Objects = append(del.Objects, types.ObjectIdentifier{Key: obj.Key})
			}

			deleteObjects, err := client.DeleteObjects(context.Background(), &s3.DeleteObjectsInput{
				Bucket: &bucketName,
				Delete: &del,
			})
			if err != nil {
				errs = append(errs, err)
				break
			}
			for _, deleteError := range deleteObjects.Errors {
				errs = append(errs, fmt.Errorf("failed to delete %s, %s", aws.ToString(deleteError.Key), aws.ToString(deleteError.Message)))
			}

			if aws.ToBool(listObjectsV2.IsTruncated) {
				input.ContinuationToken = listObjectsV2.NextContinuationToken
			} else {
				break
			}
		}
	}

	{
		slog.Info("TearDown: Deleting partial uploads from bucket", "bucket", bucketName)

		input := &s3.ListMultipartUploadsInput{
			Bucket: &bucketName,
			Prefix: aws.String(prefix),
		}
		for {
			uploads, err := client.ListMultipartUploads(context.Background(), input)
			if err != nil {
				return fmt.Errorf("failed to list multipart objects, %w", err)
			}

			for _, upload := range uploads.Uploads {
				client.AbortMultipartUpload(context.Background(), &s3.AbortMultipartUploadInput{
					Bucket:   &bucketName,
					Key:      upload.Key,
					UploadId: upload.UploadId,
				})
			}

			if aws.ToBool(uploads.IsTruncated) {
				input.KeyMarker = uploads.NextKeyMarker
				input.UploadIdMarker = uploads.NextUploadIdMarker
			} else {
				break
			}
		}
	}

	if len(errs) != 0 {
		return fmt.Errorf("failed to delete objects, %s", errs)
	}

	fmt.Println("TearDown: Deleting bucket,", bucketName)
	if _, err := client.DeleteBucket(context.Background(), &s3.DeleteBucketInput{Bucket: &bucketName}); err != nil {
		return fmt.Errorf("failed to delete bucket %s, %w", bucketName, err)
	}

	return nil
}
