package storage

import (
	"clusterlizer/internal/entity"
	"clusterlizer/pkg/pgxclient"
	"context"
)

type Storage interface {
	DoInTransaction(ctx context.Context, f pgxclient.TxFunc) error
	City() City
}

type City interface {
	GetAllWithVacancyCount(ctx context.Context) ([]entity.CityWithVacancyCount, error)
}
