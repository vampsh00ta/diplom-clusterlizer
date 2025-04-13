package postgresrep

import (
	"clusterlizer/internal/storage"
	"clusterlizer/pkg/pgxclient"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Storage struct {
	db pgxclient.DBManager
	tx pgxclient.TXManager
}

func New(db *pgxpool.Pool) *Storage {
	return &Storage{
		db: pgxclient.NewPgxManager(db),
		tx: pgxclient.NewTxManager(db),
	}
}

func (s *Storage) City() storage.City {
	return s
}

func (s *Storage) Tx() pgxclient.Manager {
	return s
}
