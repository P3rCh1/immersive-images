package object

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log/slog"

	"github.com/P3rCh1/immersive-images/backend/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type S3 struct {
	log    *slog.Logger
	client *s3.Client
}

func New(log *slog.Logger) *S3 {
	return &S3{
		client: s3.New(s3.Options{
			Region: config.Config.S3.Region,
			Credentials: credentials.NewStaticCredentialsProvider(
				config.Config.S3.AccessKey,
				config.Config.S3.SecretKey,
				"",
			),
			BaseEndpoint: aws.String(config.Config.S3.Endpoint),
			UsePathStyle: true,
		}),
		log: log,
	}
}

func (s *S3) Upload(ctx context.Context, key string, data []byte) error {
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(config.Config.S3.Bucket),
		Key:         aws.String(key),
		Body:        bytes.NewReader(data),
		ContentType: aws.String("image/png"),
	})
	if err != nil {
		s.log.Error(
			"failed to upload object to S3",
			"error", err,
		)

		return fmt.Errorf("failed to upload object to S3: %w", err)
	}

	return nil
}

func (s *S3) Get(ctx context.Context, key string) ([]byte, error) {
	obj, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(config.Config.S3.Bucket),
		Key:    aws.String(key),
	})
	if err != nil {
		s.log.Error(
			"failed to get object from S3",
			"error", err,
		)

		return nil, fmt.Errorf("failed to get object from S3: %w", err)
	}

	defer obj.Body.Close()

	data, err := io.ReadAll(obj.Body)
	if err != nil {
		s.log.Error(
			"failed to read object data",
			"error", err,
		)

		return nil, fmt.Errorf("failed to read object data: %w", err)
	}

	return data, nil
}
