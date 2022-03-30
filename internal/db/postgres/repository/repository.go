package repository

import (
	"context"
	"database/sql"
	"github.com/GeneralKenobi/mailman/internal/db"
)

// SqlExecutor is an interface for common query execution functions from sql.DB and sql.Tx.
type SqlExecutor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
}

func New(sql SqlExecutor) *Repository {
	return &Repository{sql: sql}
}

type Repository struct {
	sql SqlExecutor
}

var _ db.Repository = (*Repository)(nil) // Interface guard
