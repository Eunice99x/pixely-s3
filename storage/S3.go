package storage

import (
	"bytes"
	"context"
	"os"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3Bucket struct {
	S3Client *s3.Client
	Bucket   string
}

type S3CFG struct {
	endpoint        string
	region          string
	accessKeyId     string
	secretAccessKey string
}

func NewS3(ctx context.Context) (*S3Bucket, error) {
	c := S3CFG{
		endpoint:        os.Getenv("S3_ENDPOINT"),
		region:          os.Getenv("S3_REGION"),
		accessKeyId:     os.Getenv("S3_ACCESS_KEY"),
		secretAccessKey: os.Getenv("S3_SECRET_KEY"),
	}

	cfg, err := config.LoadDefaultConfig(ctx,
		config.WithRegion(c.region),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(c.accessKeyId, c.secretAccessKey, "")),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(cfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(c.endpoint)
		o.UsePathStyle = true
	})

	return &S3Bucket{S3Client: client, Bucket: os.Getenv("S3_BUCKET")}, nil
}

func (s *S3Bucket) Upload(ctx context.Context, key string, data []byte) error {
	_, err := s.S3Client.PutObject(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.Bucket),
		Key:    aws.String(key),
		Body:   bytes.NewReader(data),
	})
	return err
}
