package storage

import (
	"clusterlizer/internal/entity"
	"clusterlizer/pkg/pgxclient"
	"context"
)

type Storage interface {
	DoInTransaction(ctx context.Context, f pgxclient.TxFunc) error
	Request() Request
}

type Request interface {
	CreateRequest(ctx context.Context, params CreateRequestParams) (entity.Request, error)
	UpdateRequest(ctx context.Context, params UpdateRequestParams) (entity.Request, error)
	GetAllRequests(ctx context.Context) ([]entity.Request, error)
	GetRequestByID(ctx context.Context, ID entity.RequestID) (entity.Request, error)
	GetRequestByIDDone(ctx context.Context, ID entity.RequestID) (entity.Request, error)
}
