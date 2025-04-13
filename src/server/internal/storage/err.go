package storage

import (
	"errors"
	"github.com/jackc/pgx/v5"
)

var (
	NoCityError         = errors.New("no such city")
	NoSuchKeywordError  = errors.New("no such keyword")
	NullCustomerIDError = errors.New("null id")
)

func dbError(err error) error {
	switch {
	case errors.Is(err, pgx.ErrNoRows):
		return nil
	}

	return err
}
