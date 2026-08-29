package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"strings"

	"spendsnap-backend/internal/config"
	"spendsnap-backend/internal/domain/storage"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

type r2Storage struct {
	client *s3.Client
	bucket string
	pubURL string
}

func NewR2Provider(cfg config.R2Config) (storage.StorageProvider, error) {
	if cfg.AccountID == "" || cfg.AccessKeyID == "" || cfg.SecretAccessKey == "" {
		return nil, errors.New("thiếu cấu hình R2")
	}

	endpoint := "https://" + cfg.AccountID + ".r2.cloudflarestorage.com"
	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithCredentialsProvider(
			credentials.NewStaticCredentialsProvider(cfg.AccessKeyID, cfg.SecretAccessKey, ""),
		),
		awsconfig.WithRegion("auto"),
	)
	if err != nil {
		return nil, err
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.BaseEndpoint = aws.String(endpoint)
		o.UsePathStyle = true
	})

	return &r2Storage{
		client: client,
		bucket: cfg.BucketName,
		pubURL: strings.TrimRight(cfg.PublicURL, "/"), // Xóa dấu / ở cuối nếu có
	}, nil
}

func (s *r2Storage) Upload(ctx context.Context, key string, data []byte, contentType string) (string, error) {
	cleanKey := strings.TrimLeft(key, "/") // Xóa dấu / ở đầu key để tránh URL bị dúp //

	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(cleanKey),
		Body:        bytes.NewReader(data),
		ContentType: aws.String(contentType),
	})
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("%s/%s", s.pubURL, cleanKey), nil
}