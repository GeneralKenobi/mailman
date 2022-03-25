package repository

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/persistence/model"
)

func (repository *Repository) FindCustomerById(ctx context.Context, id int) (model.Customer, error) {
	return selectingOne(ctx, "find customer by ID", repository.sql, customerRowScanSupplier,
		"SELECT id, email FROM mailmandb.customer WHERE id = $1", id)
}

func (repository *Repository) FindCustomerByEmail(ctx context.Context, email string) (model.Customer, error) {
	return selectingOne(ctx, "find customer by email", repository.sql, customerRowScanSupplier,
		"SELECT id, email FROM mailmandb.customer WHERE email = $1", email)
}

func (repository *Repository) DeleteCustomerById(ctx context.Context, id int) error {
	return affectingOne(ctx, "delete customer by ID", repository.sql,
		"DELETE FROM mailmandb.customer WHERE id = $1", id)
}

func customerRowScanSupplier() (*model.Customer, []any) {
	var customer model.Customer
	return &customer, []any{
		&customer.Id,
		&customer.Email,
	}
}
