package file

import (
	"context"
	"fmt"

	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"

	"go.uber.org/zap"
)

type Service interface {
	CreateFile(ctx context.Context, params CreateFileParams) error
	GetRequestByKey(ctx context.Context, key string) (entity.File, error)
}

func New(
	storage storage.Storage,
	log *zap.SugaredLogger,
) *FileImpl {
	return &FileImpl{
		storage: storage,
		log:     log,
	}
}

type FileImpl struct {
	storage storage.Storage
	log     *zap.SugaredLogger
}
type CreateFileParams struct {
	Key  string
	Type entity.FileType
}

func (s *FileImpl) CreateFile(ctx context.Context, params CreateFileParams) error {
	s.log.Info("create file")

	_, err := s.storage.File().CreateFile(ctx, storage.CreateFileParams{
		Key:  params.Key,
		Type: params.Type,
	})
	if err != nil {
		return fmt.Errorf("create file: %w", err)
	}
	return nil
}

func (s *FileImpl) GetRequestByKey(ctx context.Context, key string) (entity.File, error) {
	s.log.Info("get file by id")

	return s.storage.File().GetFileByKey(ctx, key)
}
