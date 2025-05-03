package postgresrep

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"
	"clusterlizer/pkg/utils"

	"github.com/google/uuid"

	sq "github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5"
)

var requestFields = []string{
	fieldID,
	fieldResult,
	fieldStatus,
	fieldCreatedAt,
	fieldUpdatedAt,
}

type request struct {
	ID     uuid.UUID        `db:"id"`
	Result *json.RawMessage `db:"result"`
	Status string           `db:"status"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (s *Storage) GetRequestByIDDone(ctx context.Context, ID entity.RequestID) (entity.Request, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return entity.Request{}, fmt.Errorf("client: %w", err)
	}

	query := sq.Select(
		requestFields...,
	).
		From(tableRequest).
		Where(sq.Eq{
			fieldID:     ID,
			fieldStatus: entity.StatusDone.String(),
		}).
		PlaceholderFormat(sq.Dollar)

	q, args, err := query.ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("query builder: %w", err)
	}
	rows, err := client.Query(ctx, q, args...)
	if err != nil {
		return entity.Request{}, pgError(fmt.Errorf("row query: %w", err))
	}
	defer rows.Close()

	rowModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[request])
	if err != nil {
		return entity.Request{}, storage.DbError(fmt.Errorf("collect rows: %w", err))
	}

	return requestToEntity(rowModel)
}

func (s *Storage) GetRequestByID(ctx context.Context, ID entity.RequestID) (entity.Request, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return entity.Request{}, fmt.Errorf("client: %w", err)
	}

	query := sq.Select(
		requestFields...,
	).
		From(tableRequest).
		Where(sq.Eq{
			fieldID: ID,
		}).
		PlaceholderFormat(sq.Dollar)

	q, args, err := query.ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("query builder: %w", err)
	}
	rows, err := client.Query(ctx, q, args...)
	if err != nil {
		return entity.Request{}, pgError(fmt.Errorf("row query: %w", err))
	}
	defer rows.Close()

	rowModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[request])
	if err != nil {
		return entity.Request{}, fmt.Errorf("collect rows: %w", err)
	}

	return requestToEntity(rowModel)
}

func (s *Storage) GetAllRequests(ctx context.Context) ([]entity.Request, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return nil, fmt.Errorf("client: %w", err)
	}

	query := sq.Select(
		requestFields...,
	).
		From(tableRequest).
		PlaceholderFormat(sq.Dollar)

	q, _, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("query builder: %w", err)
	}
	rows, err := client.Query(ctx, q)
	if err != nil {
		return nil, pgError(fmt.Errorf("row query: %w", err))
	}
	defer rows.Close()

	rowModels, err := pgx.CollectRows(rows, pgx.RowToStructByName[request])
	if err != nil {
		return nil, fmt.Errorf("collect rows: %w", err)
	}

	return requestToEntities(rowModels)
}

func (s *Storage) CreateRequest(ctx context.Context, params storage.CreateRequestParams) (entity.Request, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return entity.Request{}, fmt.Errorf("client: %w", err)
	}

	q, args, err := sq.Insert(tableRequest).
		Columns(fieldID, fieldStatus).
		Values(params.ID, entity.StatusCreated.String()).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("query builder: %w", err)
	}
	row, err := client.Query(ctx, q, args...)
	if err != nil {
		return entity.Request{}, pgError(fmt.Errorf("row query: %w", err))
	}
	defer row.Close()

	rowModel, err := pgx.CollectOneRow(row, pgx.RowToStructByName[request])
	if err != nil {
		return entity.Request{}, fmt.Errorf("collect rows: %w", err)
	}

	return requestToEntity(rowModel)
}

func (s *Storage) UpdateRequest(ctx context.Context, params storage.UpdateRequestParams) (entity.Request, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return entity.Request{}, fmt.Errorf("client: %w", err)
	}

	query := sq.Update(tableRequest).
		SetMap(map[string]interface{}{
			fieldUpdatedAt: time.Now(),
		}).
		Where(sq.Eq{
			fieldID: params.ID,
		}).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar)
	query = patchQueryOptional(query, fieldResult, params.Result)
	query = patchQueryOptional(query, fieldStatus, params.Status)

	q, args, err := query.ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("query builder: %w", err)
	}
	row, err := client.Query(ctx, q, args...)
	if err != nil {
		return entity.Request{}, pgError(fmt.Errorf("row query: %w", err))
	}
	defer row.Close()

	rowModel, err := pgx.CollectOneRow(row, pgx.RowToStructByName[request])
	if err != nil {
		return entity.Request{}, fmt.Errorf("collect rows: %w", err)
	}
	return requestToEntity(rowModel)
}

func requestToEntity(r request) (entity.Request, error) {
	var result entity.GraphData
	if r.Result != nil {
		if err := json.Unmarshal(utils.SafeNil(r.Result), &result); err != nil {
			return entity.Request{}, err
		}
	}

	return entity.Request{
		ID:        entity.RequestID(r.ID.String()),
		Result:    result,
		Status:    entity.StatusFromString(r.Status),
		CreatedAt: r.CreatedAt,
		UpdatedAt: r.UpdatedAt,
	}, nil
}

func requestToEntities(requests []request) ([]entity.Request, error) {
	res := make([]entity.Request, 0, len(requests))
	for _, r := range requests {
		var result entity.GraphData

		if r.Result != nil {
			if err := json.Unmarshal(utils.SafeNil(r.Result), &result); err != nil {
				return nil, err
			}
		}
		res = append(res, entity.Request{
			ID:        entity.RequestID(r.ID.String()),
			Result:    result,
			Status:    entity.StatusFromString(r.Status),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	return res, nil
}
