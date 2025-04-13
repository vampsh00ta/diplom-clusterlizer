package pgxclient

import "context"

type CtxTrKey struct{}

type CtxManager interface {
	Set(ctx context.Context) context.Context
	SetByKey(ctx context.Context, key interface{}) context.Context
	Get(ctx context.Context) interface{}
	GetByKey(ctx context.Context, key interface{}) interface{}
}

type ctxManager struct {
	dataType interface{}
}

func NewCtxManager(dataType interface{}) CtxManager {
	return &ctxManager{
		dataType: dataType,
	}
}
func (m ctxManager) SetByKey(ctx context.Context, key interface{}) context.Context {
	ctx = context.WithValue(ctx, key, m.dataType)
	return ctx
}
func (m ctxManager) Set(ctx context.Context) context.Context {
	ctx = context.WithValue(ctx, CtxTrKey{}, m.dataType)
	return ctx
}

func (m ctxManager) Get(ctx context.Context) interface{} {
	val := ctx.Value(CtxTrKey{})
	return val
}
func (m ctxManager) GetByKey(ctx context.Context, key interface{}) interface{} {
	val := ctx.Value(key)
	return val
}
