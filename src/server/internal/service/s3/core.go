package s3

import (
	"context"
	"fmt"
	"strconv"

	"clusterlizer/pkg/s3"

	"go.uber.org/zap"
)

type Service interface {
	Upload(ctx context.Context, key string, fileBytes []byte) error
	UploadMultiple(ctx context.Context, key string, fileBytes [][]byte) error
	Download(ctx context.Context, key string) ([]byte, error)
}

type ServiceImpl struct {
	log    *zap.SugaredLogger
	client s3.Client
}

func New(
	log *zap.SugaredLogger,
	client s3.Client,
) *ServiceImpl {
	return &ServiceImpl{
		client: client,
		log:    log,
	}
}

func (s *ServiceImpl) Upload(ctx context.Context, key string, fileBytes []byte) error {
	s.log.Infof("upload key: %s", key)

	return s.client.Upload(ctx, key, fileBytes)
}

func (s *ServiceImpl) Download(ctx context.Context, key string) ([]byte, error) {
	s.log.Infof("upload key: %s", key)

	return s.client.Get(ctx, key)
}

func (s *ServiceImpl) UploadMultiple(ctx context.Context, key string, fileBytes [][]byte) error {
	for i, bytes := range fileBytes {
		fileKey := fmt.Sprintf("%s_%s", key, strconv.Itoa(i))
		if err := s.client.Upload(ctx, fileKey, bytes); err != nil {
			return fmt.Errorf("multiple upload: %w", err)
		}
	}
	return nil
}
