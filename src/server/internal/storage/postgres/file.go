package postgresrep

import (
	"context"
	"fmt"

	"clusterlizer/internal/entity"
	"clusterlizer/internal/storage"

	sq "github.com/Masterminds/squirrel"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
)

var fileFields = []string{
	fieldID,
	fileTitle,
	fieldKey,
	fieldType,
}

type file struct {
	ID    uuid.UUID `json:"id"`
	Key   string    `json:"key"`
	Type  string    `json:"type"`
	Title string    `json:"title"`
}

func (s *Storage) GetFileByKey(ctx context.Context, key string) (entity.File, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return entity.File{}, fmt.Errorf("client: %w", err)
	}

	query := sq.Select(
		fileFields...,
	).
		From(tableFile).
		Where(sq.Eq{
			fieldKey: key,
		}).
		PlaceholderFormat(sq.Dollar)

	q, args, err := query.ToSql()
	if err != nil {
		return entity.File{}, fmt.Errorf("query builder: %w", err)
	}
	rows, err := client.Query(ctx, q, args...)
	if err != nil {
		return entity.File{}, pgError(fmt.Errorf("row query: %w", err))
	}
	defer rows.Close()

	rowModel, err := pgx.CollectOneRow(rows, pgx.RowToStructByName[file])
	if err != nil {
		return entity.File{}, fmt.Errorf("collect rows: %w", err)
	}

	return fileToEntity(rowModel)
}

func (s *Storage) CreateFile(ctx context.Context, params storage.CreateFileParams) (entity.File, error) {
	client, err := s.db.Client(ctx)
	if err != nil {
		return entity.File{}, fmt.Errorf("client: %w", err)
	}
	id := uuid.New()
	q, args, err := sq.Insert(tableFile).
		Columns(fieldID, fieldKey, fieldType, fileTitle).
		Values(id, params.Key, params.Type, params.Title).
		Suffix("RETURNING *").
		PlaceholderFormat(sq.Dollar).
		ToSql()
	if err != nil {
		return entity.File{}, fmt.Errorf("query builder: %w", err)
	}
	row, err := client.Query(ctx, q, args...)
	if err != nil {
		return entity.File{}, pgError(fmt.Errorf("row query: %w", err))
	}
	defer row.Close()

	rowModel, err := pgx.CollectOneRow(row, pgx.RowToStructByName[file])
	if err != nil {
		return entity.File{}, fmt.Errorf("collect rows: %w", err)
	}

	return fileToEntity(rowModel)
}

func fileToEntity(f file) (entity.File, error) {
	return entity.File{
		ID:    entity.FileID(f.ID.String()),
		Key:   f.Key,
		Type:  entity.FileTypeFromString(f.Type),
		Title: f.Title,
	}, nil
}
