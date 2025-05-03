package pgxclient

import (
	"context"
)

type DBManager interface {
	Client(ctx context.Context) (Client, error)
	ClientByKey(ctx context.Context, key interface{}) (Client, error)
}

type pgxManager struct {
	db   Client
	ctxm CtxManager
}

func NewPgxManager(db Client) DBManager {
	return &pgxManager{
		db:   db,
		ctxm: NewCtxManager(db),
	}
}

func (p pgxManager) Client(ctx context.Context) (Client, error) {
	val := p.ctxm.Get(ctx)
	ctxClient, ok := val.(Client)
	if ok {
		return ctxClient, nil
	}

	return p.db, nil
}

func (p pgxManager) ClientByKey(ctx context.Context, key interface{}) (Client, error) {
	val := p.ctxm.GetByKey(ctx, key)
	ctxClient, ok := val.(Client)
	if ok {
		return ctxClient, nil
	}
	return p.db, nil
}
