package pgxclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"clusterlizer/pkg/utils"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	// BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type Tx interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...any) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
type Config struct {
	Username string
	Password string
	Host     string
	Port     string
	Name     string
}

func New(ctx context.Context, maxAttempts int, config Config) (pool *pgxpool.Pool, err error) {
	dsn := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s",
		config.Username,
		config.Password,
		config.Host,
		config.Port,
		config.Name)
	err = utils.DoWithTries(func() error {
		ctxWithTimeout, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()
		pool, err = pgxpool.New(ctxWithTimeout, dsn)
		if err != nil {
			return err
		}
		return nil
	}, maxAttempts, 5)
	if err != nil {
		log.Fatal("error with tries to connect")
	}

	return pool, nil
}
