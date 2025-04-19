package request

import (
	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"
	"clusterlizer/pkg/utils"
	"context"

	"go.uber.org/zap"
)

type Service interface {
	CreateRequest(ctx context.Context, params CreateRequestParams) error
	UpdateRequest(ctx context.Context, params UpdateRequestParams) (entity.Request, error)
	GetAllRequests(ctx context.Context) ([]entity.Request, error)
}

func NewRequest(
	storage storage.Storage,
	log *zap.SugaredLogger) *RequestImpl {
	return &RequestImpl{
		storage: storage,
		log:     log,
	}
}

type RequestImpl struct {
	storage storage.Storage
	log     *zap.SugaredLogger
}
type CreateRequestParams struct {
	ID entity.RequestID `db:"id"`
}

func (s *RequestImpl) CreateRequest(ctx context.Context, params CreateRequestParams) error {
	s.log.Info("create request")

	return s.storage.Request().CreateRequest(ctx, storage.CreateRequestParams{
		ID: params.ID,
	})
}

type UpdateRequestParams struct {
	ID entity.RequestID

	Result utils.Optional[*[]byte]
	Status utils.Optional[entity.Status]
}

func (s *RequestImpl) UpdateRequest(ctx context.Context, params UpdateRequestParams) (entity.Request, error) {
	s.log.Info("update request")

	return s.storage.Request().UpdateRequest(ctx, storage.UpdateRequestParams{
		ID:     params.ID,
		Result: params.Result,
		Status: params.Status,
	})
}

func (s *RequestImpl) GetAllRequests(ctx context.Context) ([]entity.Request, error) {
	s.log.Info("get all requests")

	return s.storage.Request().GetAllRequests(ctx)
}
