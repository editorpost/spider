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
	Name   string
	URL    string
	Access string
	Secret string
	Region string
}

type BucketStore struct {
	bucket string
	client *s3.Client
}

func NewBucketStore(bucket string, client *s3.Client) *BucketStore {
	return &BucketStore{
		bucket: bucket,
		client: client,
	}
}

func (b *BucketStore) Save(data []byte, path string) error {

	uploader := manager.NewUploader(b.client)

	_, err := uploader.Upload(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(b.bucket),
		Key:    aws.String(path),
		Body:   bytes.NewReader(data),
	})

	if err != nil {
		return fmt.Errorf("failed to upload object, %w", err)
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
		o.BaseEndpoint = aws.String(bucket.URL)
		o.UsePathStyle = true
	})

	return client, nil
}
