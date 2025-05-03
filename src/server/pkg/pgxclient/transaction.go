package pgxclient

import "context"

type (
	TxFunc  func(ctx context.Context) error
	Manager interface {
		DoInTransaction(ctx context.Context, f TxFunc) error
	}
)
