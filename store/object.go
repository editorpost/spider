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
	"path/filepath"
)

type Bucket struct {
	Name   string
	URL    string
	Access string
	Secret string
	Region string
}

type BucketFolder struct {
	bucket Bucket
	prefix string
	client *s3.Client
}

func NewBucketFolder(path string, bucket Bucket, client *s3.Client) (*BucketFolder, error) {
	return &BucketFolder{
		bucket: bucket,
		prefix: path,
		client: client,
	}, nil
}

func (b *BucketFolder) Save(data []byte, filename string) (string, error) {

	uploader := manager.NewUploader(b.client)

	result, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.bucket.Name),
		Key:    aws.String(b.Path(filename)),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return "", fmt.Errorf("failed to upload object, %w", err)
	}

	return result.Location, nil
}

// Path returns the full path of the file in the bucket.
func (b *BucketFolder) Path(filename string) string {
	return filepath.Join(b.prefix, filename)
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
		o.BaseEndpoint = aws.String(bucket.URL)
		o.UsePathStyle = true
	})

	return client, nil
}
