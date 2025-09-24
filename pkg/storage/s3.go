package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
)

type S3Client interface {
	Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error
	PresignGet(ctx context.Context, key string, expireSec int) (string, error)
}

type MinioClient struct {
	client *minio.Client
	bucket string
}

func NewMinioClient(endpoint, accessKey, secretKey, bucket string, useSSL bool) (*MinioClient, error) {
	if strings.TrimSpace(endpoint) == "" {
		return nil, errors.New("minio endpoint is required")
	}
	if strings.TrimSpace(bucket) == "" {
		return nil, errors.New("minio bucket is required")
	}

	cli, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, fmt.Errorf("create minio client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	exists, err := cli.BucketExists(ctx, bucket)
	if err != nil {
		return nil, fmt.Errorf("check bucket %s: %w", bucket, err)
	}
	if !exists {
		if err := cli.MakeBucket(ctx, bucket, minio.MakeBucketOptions{}); err != nil {
			return nil, fmt.Errorf("create bucket %s: %w", bucket, err)
		}
	}

	return &MinioClient{client: cli, bucket: bucket}, nil
}

func (m *MinioClient) Put(ctx context.Context, key string, r io.Reader, size int64, contentType string) error {
	if m == nil || m.client == nil {
		return errors.New("minio client is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return err
	}
	if strings.TrimSpace(key) == "" {
		return errors.New("object key is required")
	}
	opts := minio.PutObjectOptions{}
	if strings.TrimSpace(contentType) != "" {
		opts.ContentType = contentType
	}
	_, err := m.client.PutObject(ctx, m.bucket, key, r, size, opts)
	if err != nil {
		return fmt.Errorf("put object %s: %w", key, err)
	}
	return nil
}

func (m *MinioClient) PresignGet(ctx context.Context, key string, expireSec int) (string, error) {
	if m == nil || m.client == nil {
		return "", errors.New("minio client is not initialized")
	}
	if err := ctx.Err(); err != nil {
		return "", err
	}
	if strings.TrimSpace(key) == "" {
		return "", errors.New("object key is required")
	}
	if expireSec <= 0 {
		expireSec = 600
	}

	url, err := m.client.PresignedGetObject(ctx, m.bucket, key, time.Duration(expireSec)*time.Second, nil)
	if err != nil {
		return "", fmt.Errorf("presign get for %s: %w", key, err)
	}
	return url.String(), nil
}
