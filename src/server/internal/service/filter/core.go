package service

import (
	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"
	"context"
)

type AllFilter interface {
	GetAll(ctx context.Context) (entity.AllFilter, error)
}

type FilterImpl struct {
	storage storage.Storage
}

func New(
	storage storage.Storage,
) *FilterImpl {
	return &FilterImpl{
		storage: storage,
	}
}

func (a FilterImpl) GetAll(ctx context.Context) (entity.AllFilter, error) {

	return entity.AllFilter{}, nil
}
