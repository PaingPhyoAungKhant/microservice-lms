package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/paingphyoaungkhant/asto-microservice/shared/config"
	"github.com/paingphyoaungkhant/asto-microservice/shared/logger"
	"go.uber.org/zap"
)

type MinIOClient struct {
	client *minio.Client
	logger *logger.Logger
	buckets []string
}

func NewMinIOClient(cfg config.MinIOConfig, log *logger.Logger) (*MinIOClient, error) {
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
		Region: cfg.Region,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	mc := &MinIOClient{
		client: client,
		logger: log,
		buckets: []string{"course-thumbnails", "course-videos", "zoom-recordings", "general-files"},
	}

	
	if err := mc.ensureBuckets(context.Background()); err != nil {
		return nil, fmt.Errorf("failed to ensure buckets: %w", err)
	}

	return mc, nil
}

func (m *MinIOClient) ensureBuckets(ctx context.Context) error {
	for _, bucketName := range m.buckets {
		exists, err := m.client.BucketExists(ctx, bucketName)
		if err != nil {
			return fmt.Errorf("failed to check bucket existence: %w", err)
		}

		if !exists {
			err = m.client.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
				Region: "us-east-1",
			})
			if err != nil {
				return fmt.Errorf("failed to create bucket %s: %w", bucketName, err)
			}
			m.logger.Info("created bucket", zap.String("bucket", bucketName))
		}
	}
	return nil
}

func (m *MinIOClient) UploadFile(ctx context.Context, bucketName, objectName string, reader io.Reader, size int64, contentType string) error {
	_, err := m.client.PutObject(ctx, bucketName, objectName, reader, size, minio.PutObjectOptions{
		ContentType: contentType,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file: %w", err)
	}

	m.logger.Info("file uploaded",
		zap.String("bucket", bucketName),
		zap.String("object", objectName),
		zap.Int64("size", size),
	)
	return nil
}

func (m *MinIOClient) DownloadFile(ctx context.Context, bucketName, objectName string) (io.ReadCloser, error) {
	obj, err := m.client.GetObject(ctx, bucketName, objectName, minio.GetObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get object: %w", err)
	}

	return obj, nil
}

func (m *MinIOClient) DeleteFile(ctx context.Context, bucketName, objectName string) error {
	err := m.client.RemoveObject(ctx, bucketName, objectName, minio.RemoveObjectOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete object: %w", err)
	}

	m.logger.Info("file deleted",
		zap.String("bucket", bucketName),
		zap.String("object", objectName),
	)
	return nil
}

func (m *MinIOClient) GetPresignedURL(ctx context.Context, bucketName, objectName string, expiry time.Duration) (string, error) {
	url, err := m.client.PresignedGetObject(ctx, bucketName, objectName, expiry, nil)
	if err != nil {
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}
	return url.String(), nil
}

func (m *MinIOClient) GetObjectInfo(ctx context.Context, bucketName, objectName string) (*minio.ObjectInfo, error) {
	info, err := m.client.StatObject(ctx, bucketName, objectName, minio.StatObjectOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to stat object: %w", err)
	}
	return &info, nil
}

