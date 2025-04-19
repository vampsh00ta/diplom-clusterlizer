package postgresrep

import (
	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"

	"context"
	"fmt"
	"time"

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
	ID     entity.RequestID `db:"id"`
	Result []byte           `db:"result"`
	Status string           `db:"status"`

	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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

	return requestToEntity(rowModels), nil
}

func (s *Storage) CreateRequest(ctx context.Context, params storage.CreateRequestParams) error {
	client, err := s.db.Client(ctx)
	if err != nil {
		return fmt.Errorf("client: %w", err)
	}

	q, args, err := sq.Insert(tableRequest).
		Columns(fieldID).
		PlaceholderFormat(sq.Dollar).
		Values(params.ID).
		ToSql()

	if err != nil {
		return fmt.Errorf("query builder: %w", err)
	}
	if err := client.QueryRow(ctx, q, args...).Scan(); err != nil {
		return pgError(fmt.Errorf("row query: %w", err))
	}

	return nil
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
		PlaceholderFormat(sq.Dollar)
	query = patchQueryOptional(query, fieldResult, params.Result)
	query = patchQueryOptional(query, fieldStatus, params.Status)

	q, args, err := query.ToSql()
	if err != nil {
		return entity.Request{}, fmt.Errorf("query builder: %w", err)
	}
	
	var req request
	if err := client.QueryRow(ctx, q, args...).Scan(&req); err != nil {
		return entity.Request{}, pgError(fmt.Errorf("row query: %w", err))
	}

	return requestToEntity([]request{req})[0], nil
}

func requestToEntity(requests []request) []entity.Request {
	res := make([]entity.Request, 0, len(requests))
	for _, r := range requests {
		res = append(res, entity.Request{
			ID:        r.ID,
			Result:    r.Result,
			Status:    entity.StatusFromString(r.Status),
			CreatedAt: r.CreatedAt,
			UpdatedAt: r.UpdatedAt,
		})
	}
	return res
}
