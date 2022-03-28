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

func (repository *Repository) InsertCustomer(ctx context.Context, customer model.Customer) (model.Customer, error) {
	return selectingOne(ctx, "insert customer", repository.sql, customerRowScanSupplier,
		"INSERT INTO mailmandb.customer(email) VALUES($1) RETURNING id, email", customer.Email)
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
