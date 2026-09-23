package object

import (
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
