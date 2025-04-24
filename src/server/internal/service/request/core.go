package request

import (
	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"
	"clusterlizer/pkg/utils"
	"context"
	"fmt"

	"go.uber.org/zap"
)

type Service interface {
	CreateRequest(ctx context.Context, params CreateRequestParams) error
	UpdateRequest(ctx context.Context, params UpdateRequestParams) (entity.Request, error)
	GetAllRequests(ctx context.Context) ([]entity.Request, error)
	GetRequestByID(ctx context.Context, ID entity.RequestID) (entity.Request, error)
	GetRequestByIDDone(ctx context.Context, ID entity.RequestID) (entity.Request, error)
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

	_, err := s.storage.Request().CreateRequest(ctx, storage.CreateRequestParams{
		ID: params.ID,
	})
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	return nil
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

func (s *RequestImpl) GetRequestByID(ctx context.Context, ID entity.RequestID) (entity.Request, error) {
	s.log.Info("get request by id")

	return s.storage.Request().GetRequestByID(ctx, ID)
}

func (s *RequestImpl) GetRequestByIDDone(ctx context.Context, ID entity.RequestID) (entity.Request, error) {
	s.log.Info("get request by id")

	return s.storage.Request().GetRequestByIDDone(ctx, ID)
}
