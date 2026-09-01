package core_postgres_pool

import (
	"context"
	"time"
)

type Pool interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)
	Close()

	BeginTx(ctx context.Context) (Tx, error)

	Ping(ctx context.Context) error
	OpTimeout() time.Duration
}

// Tx is a database transaction with the same query surface as Pool. Callers
// must call Commit or Rollback exactly once; deferring Rollback and calling
// Commit on the success path is the expected pattern (Rollback after a
// successful Commit is a documented no-op in pgx, so the defer is safe).
type Tx interface {
	Query(ctx context.Context, sql string, args ...any) (Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) Row
	Exec(ctx context.Context, sql string, arguments ...any) (CommandTag, error)

	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}

type Rows interface {
	Close()
	Err() error
	Next() bool
	Scan(dest ...any) error
}

type Row interface {
	Scan(dest ...any) error
}

type CommandTag interface {
	RowsAffected() int64
}
