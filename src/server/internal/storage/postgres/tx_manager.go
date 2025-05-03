package postgresrep

import (
	"context"
	"fmt"

	"clusterlizer/pkg/pgxclient"
)

func (s *Storage) DoInTransaction(ctx context.Context, f pgxclient.TxFunc) error {
	tx, ctx, err := s.tx.Create(ctx)
	if err != nil {
		return fmt.Errorf("tx create: %w", err)
	}
	defer func() {
		_ = tx.Commit(ctx)
	}()
	// вызываем фукнцию с новым контекстом
	err = f(ctx)
	// закрывыем транзанцию
	if err != nil {
		_ = tx.Rollback(ctx)
		return err
	}
	return nil
}
