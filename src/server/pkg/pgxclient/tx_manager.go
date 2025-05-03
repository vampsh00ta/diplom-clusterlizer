package pgxclient

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type txManager struct {
	ctxm CtxManager
	db   *pgxpool.Pool
}

type TXManager interface {
	Create(ctx context.Context) (Tx, context.Context, error)
	CreateByKey(ctx context.Context, key interface{}) (Tx, context.Context, error)
}

func NewTxManager(db *pgxpool.Pool) TXManager {
	return &txManager{
		ctxm: NewCtxManager(db),
		db:   db,
	}
}

func (p txManager) Create(ctx context.Context) (Tx, context.Context, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("tx error: %w", err)
	}
	ctx = p.ctxm.Set(ctx, tx)
	return tx, ctx, nil
}

func (p txManager) CreateByKey(ctx context.Context, key interface{}) (Tx, context.Context, error) {
	tx, err := p.db.Begin(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("tx error: %w", err)
	}
	ctx = p.ctxm.SetByKey(ctx, key, tx)
	return tx, ctx, nil
}
