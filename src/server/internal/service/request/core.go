package request

import (
	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"
	"clusterlizer/pkg/utils"
	"context"
	"encoding/json"
	"fmt"

	"go.uber.org/zap"
)

type Service interface {
	CreateRequest(ctx context.Context, params CreateRequestParams) error
	UpdateRequest(ctx context.Context, params UpdateRequestParams) (entity.Request, error)
	GetAllRequests(ctx context.Context) ([]entity.Request, error)
	GetRequestByID(ctx context.Context, ID entity.RequestID) (entity.Request, error)
	GetRequestByIDDone(ctx context.Context, ID entity.RequestID) (entity.Request, error)
	SaveResult(ctx context.Context, params SaveResultParams) (err error)
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
	res, err := s.storage.Request().GetRequestByID(ctx, ID)
	switch {
	case err != nil:
		return entity.Request{}, fmt.Errorf("get request by id: %w", err)
	case res.Status == entity.StatusError:
		return entity.Request{}, entity.ErrReqFailed
	case res.Status == entity.StatusCreated:
		return entity.Request{}, entity.ErrNoResult
	}

	return res, nil
}

type SaveResultParams struct {
	ID    entity.RequestID
	Graph entity.GraphData
}

func (s *RequestImpl) SaveResult(ctx context.Context, params SaveResultParams) (err error) {
	s.log.Info("update request")

	defer func() {
		if err != nil {
			_, errUpdate := s.storage.Request().UpdateRequest(ctx, storage.UpdateRequestParams{
				ID:     params.ID,
				Result: utils.NewEmptyOptional[*[]byte](),
				Status: utils.NewOptional(entity.StatusError),
			})
			if errUpdate != nil {
				err = errUpdate
			}
		}
	}()

	graphBytes, err := json.Marshal(params.Graph)
	if err != nil {
		return fmt.Errorf("save result: %w", err)
	}

	err = s.storage.DoInTransaction(ctx, func(ctx context.Context) error {
		_, err = s.storage.Request().UpdateRequest(ctx, storage.UpdateRequestParams{
			ID:     params.ID,
			Result: utils.NewOptional(&graphBytes),
			Status: utils.NewOptional(entity.StatusDone),
		})
		if err != nil {
			return fmt.Errorf("update request: %w", err)
		}

		for _, node := range params.Graph.Nodes {
			_, err = s.storage.File().CreateFile(ctx, storage.CreateFileParams{
				Key:   node.ID,
				Type:  entity.FileTypeFromString(node.Type),
				Title: node.Title,
			})
			if err != nil {
				return fmt.Errorf("create file: %w", err)
			}
		}
		return nil
	})
	return err
}
