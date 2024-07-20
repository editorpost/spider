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
)

type Bucket struct {
	Name      string
	Endpoint  string
	Access    string
	Secret    string
	Region    string
	PublicURL string
}

type BucketStore struct {
	folder string
	bucket string
	client *s3.Client
}

func NewBucketStore(bucket string, folder string, client *s3.Client) *BucketStore {
	return &BucketStore{
		folder: folder,
		bucket: bucket,
		client: client,
	}
}

func (b *BucketStore) path(filename string) *string {

	if len(b.folder) == 0 {
		return aws.String(filename)
	}

	return aws.String(fmt.Sprintf("%s/%s", b.folder, filename))
}

func (b *BucketStore) Save(data []byte, filename string) error {

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

func (b *BucketStore) Load(filename string) ([]byte, error) {

	obj, err := b.client.GetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    b.path(filename),
	})

	if err != nil {
		return nil, fmt.Errorf("failed to get object, %w", err)
	}

	buf := new(bytes.Buffer)
	_, err = buf.ReadFrom(obj.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read object, %w", err)
	}

	return buf.Bytes(), nil
}

func (b *BucketStore) Delete(filename string) error {
	_, err := b.client.DeleteObject(context.TODO(), &s3.DeleteObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    b.path(filename),
	})
	if err != nil {
		return fmt.Errorf("failed to delete object, %w", err)
	}
	return nil
}

func NewS3Client(bucket Bucket) (*s3.Client, error) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(bucket.Access, bucket.Secret, "")),
		config.WithRegion(bucket.Region),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to load configuration, %w", err)
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(bucket.Endpoint)
		o.UsePathStyle = true
	})

	return client, nil
}
