package storage

import (
	"context"

	"clusterlizer/internal/entity"
	"clusterlizer/pkg/pgxclient"
)

type Storage interface {
	DoInTransaction(ctx context.Context, f pgxclient.TxFunc) error
	Request() Request
	File() File
}

type Request interface {
	CreateRequest(ctx context.Context, params CreateRequestParams) (entity.Request, error)
	UpdateRequest(ctx context.Context, params UpdateRequestParams) (entity.Request, error)
	GetAllRequests(ctx context.Context) ([]entity.Request, error)
	GetRequestByID(ctx context.Context, ID entity.RequestID) (entity.Request, error)
}

type File interface {
	CreateFile(ctx context.Context, params CreateFileParams) (entity.File, error)
	GetFileByKey(ctx context.Context, key string) (entity.File, error)
}
