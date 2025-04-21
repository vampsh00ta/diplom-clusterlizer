package storage

import (
	"errors"
	"github.com/jackc/pgx/v5"
)

func dbError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil
	}

	return err
}
