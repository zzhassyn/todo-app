package core_pgx_pool

import (
	"context"
	"errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	core_postgres_pool "github.com/zzhassyn/todo-app/internal/core/repository/postgres/pool"
)

const (
	pgErrCodeUniqueViolation     = "23505"
	pgErrCodeForeignKeyViolation = "23503"
)

type pgxRows struct {
	pgx.Rows
}

type pgxRow struct {
	pgx.Row
}

func (r pgxRow) Scan(dest ...any) error {
	err := r.Row.Scan(dest...)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return core_postgres_pool.ErrNoRows
		}

		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			switch pgErr.Code {
			case pgErrCodeUniqueViolation:
				return core_postgres_pool.ErrUniqueViolation
			case pgErrCodeForeignKeyViolation:
				return core_postgres_pool.ErrForeignKeyViolation
			}
		}

		return err
	}
	return nil
}

type pgxCommandTag struct {
	pgconn.CommandTag
}

// pgxTx adapts a pgx.Tx to the core_postgres_pool.Tx interface, reusing
// the same Scan error-translation as the pool-level pgxRow.
type pgxTx struct {
	pgx.Tx
}

func (t pgxTx) Query(ctx context.Context, sql string, args ...any) (core_postgres_pool.Rows, error) {
	rows, err := t.Tx.Query(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxRows{rows}, nil
}

func (t pgxTx) QueryRow(ctx context.Context, sql string, args ...any) core_postgres_pool.Row {
	row := t.Tx.QueryRow(ctx, sql, args...)
	return pgxRow{row}
}

func (t pgxTx) Exec(ctx context.Context, sql string, args ...any) (core_postgres_pool.CommandTag, error) {
	tag, err := t.Tx.Exec(ctx, sql, args...)
	if err != nil {
		return nil, err
	}
	return pgxCommandTag{tag}, nil
}

func (t pgxTx) Commit(ctx context.Context) error {
	return t.Tx.Commit(ctx)
}

func (t pgxTx) Rollback(ctx context.Context) error {
	return t.Tx.Rollback(ctx)
}
