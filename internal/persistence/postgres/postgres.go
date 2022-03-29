package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	"github.com/GeneralKenobi/mailman/internal/persistence/postgres/repository"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	_ "github.com/lib/pq" // Postgres driver registration by import.
)

// Context implements DB integration for postgres.
type Context struct {
	db *sql.DB
}

var _ persistence.Context = (*Context)(nil) // Interface guard

// NewContext creates a postgres DB context. The DB client is closed when context is canceled.
func NewContext(ctx shutdown.Context) (*Context, error) {
	cfg := config.Get().Postgres
	// TODO: Add configuration for ssl
	connectionString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, fmt.Errorf("connection configuration is invalid: %w", err)
	}

	dbCtx := Context{db: db}
	go shutdownDbOnContextCancellation(ctx, db)
	return &dbCtx, nil
}

func shutdownDbOnContextCancellation(ctx shutdown.Context, db *sql.DB) {
	defer ctx.Notify()

	<-ctx.Done()
	mdctx.Infof(nil, "DB context canceled")
	shutdownDb(db)
}

func shutdownDb(db *sql.DB) {
	mdctx.Infof(nil, "Shutting down DB connection")
	err := db.Close()
	if err != nil {
		mdctx.Errorf(nil, "Error closing DB connection: %v", err)
	}
	mdctx.Infof(nil, "DB connection closed")
}

func (postgresCtx *Context) Repository(ctx context.Context) (persistence.Repository, error) {
	return repository.New(postgresCtx.db), nil
}

func (postgresCtx *Context) TransactionalRepository(ctx context.Context) (persistence.Repository, persistence.Transaction, error) {
	transaction, err := postgresCtx.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, nil, fmt.Errorf("error beginning a transaction: %w", err)
	}

	transactionalRepository := repository.New(transaction)
	transactionCtx := &transactionContext{transaction: transaction}
	return transactionalRepository, transactionCtx, nil
}

// transactionContext implements the persistence.Transaction interface.
type transactionContext struct {
	transaction *sql.Tx
}

var _ persistence.Transaction = (*transactionContext)(nil) // Interface guard

func (transactionCtx *transactionContext) Commit() error {
	return transactionCtx.transaction.Commit()
}

func (transactionCtx *transactionContext) Rollback() error {
	return transactionCtx.transaction.Rollback()
}
