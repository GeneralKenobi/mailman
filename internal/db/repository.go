package db

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/db/model"
	"time"
)

// Context is the top-level interface implemented by db providers with support for both transactional and no-transaction-guarantee
// repositories.
type Context interface {
	Querent
	Transactioner
}

// Querent is a db manager that creates repositories without any transaction guarantees. Queries may be executed without a
// transaction or within a default one, depending on the underlying DB technology.
type Querent interface {
	// Repository creates a Repository that executes queries without a transaction or within a default one, depending on the underlying
	// DB technology.
	Repository(ctx context.Context) (Repository, error)
}

// Transactioner is a db manager that creates transaction-scoped repositories.
type Transactioner interface {
	// TransactionalRepository creates a Repository that runs all queries within the returned Transaction. It's not allowed to use the
	// repository after committing or rolling back the transaction.
	TransactionalRepository(ctx context.Context) (Repository, Transaction, error)
}

type Transaction interface {
	Commit() error
	Rollback() error
}

// Repository aggregates all queries implemented by db providers.
type Repository interface {
	CustomerRepository
	MailingEntryRepository
}

type CustomerRepository interface {
	FindCustomerById(ctx context.Context, id int) (model.Customer, error)
	FindCustomerByEmail(ctx context.Context, email string) (model.Customer, error)

	InsertCustomer(ctx context.Context, customer model.Customer) (model.Customer, error)

	DeleteCustomerById(ctx context.Context, id int) error
}

type MailingEntryRepository interface {
	FindMailingEntriesByMailingId(ctx context.Context, mailingId int) ([]model.MailingEntry, error)
	FindMailingEntriesOlderThan(ctx context.Context, olderThan time.Duration) ([]model.MailingEntry, error)
	FindMailingEntriesByMailingIdOlderThan(ctx context.Context, mailingId int, olderThan time.Duration) ([]model.MailingEntry, error)
	FindMailingEntriesByCustomerId(ctx context.Context, id int) ([]model.MailingEntry, error)
	FindMailingEntriesByCustomerIdMailingIdTitleContentInsertTime(
		ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error)

	InsertMailingEntry(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error)

	DeleteMailingEntryById(ctx context.Context, id int) error
}

var (
	// ErrNoRows is returned from queries that returned/affected 0 rows but at least 1 was expected (e.g. select one found no rows).
	ErrNoRows = fmt.Errorf("no row matched the query")
	// ErrTooManyRows is returned from queries that returned/affected more rows than expected (e.g. select one found 2 rows).
	ErrTooManyRows = fmt.Errorf("more rows than expected matched the query")
)
