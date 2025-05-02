package storage

import (
	"clusterlizer/internal/entity"
	"errors"
	"github.com/jackc/pgx/v5"
)

func DbError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return entity.ErrNoResult
	}

	return err
}
