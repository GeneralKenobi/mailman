package creator

import (
	"context"
	"errors"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	"github.com/GeneralKenobi/mailman/internal/persistence/model"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
)

type Repository interface {
	FindCustomerByEmail(ctx context.Context, email string) (model.Customer, error)
	InsertCustomer(ctx context.Context, customer model.Customer) (model.Customer, error)
}

func New(repository Repository) *Creator {
	return &Creator{repository: repository}
}

type Creator struct {
	repository Repository
}

// CreateFromEmail creates a customer with the given email.
// Returns api.StatusBadInput error if the email is already assigned to a customer.
func (creator *Creator) CreateFromEmail(ctx context.Context, email string) (model.Customer, error) {
	customer := model.Customer{Email: email}
	return creator.Create(ctx, customer)
}

// Create saves the given customer.
// Returns api.StatusBadInput error if the customer's email is already assigned to a customer.
func (creator *Creator) Create(ctx context.Context, customer model.Customer) (model.Customer, error) {
	if err := creator.assertEmailIsNotUsed(ctx, customer.Email); err != nil {
		return model.Customer{}, err
	}

	mdctx.Debugf(ctx, "Creating a new customer")
	customer, err := creator.repository.InsertCustomer(ctx, customer)
	if err != nil {
		return model.Customer{}, fmt.Errorf("error creating customer: %w", err)
	}
	mdctx.Infof(ctx, "Created customer %d", customer.Id)
	return customer, nil
}

func (creator *Creator) assertEmailIsNotUsed(ctx context.Context, email string) error {
	customer, err := creator.repository.FindCustomerByEmail(ctx, email)
	if err == nil {
		mdctx.Debugf(ctx, "Email is already used by customer %d", customer.Id)
		return api.StatusBadInput.WithMessage("customer with this email already exists")
	}
	if errors.Is(err, persistence.ErrNoRows) {
		mdctx.Debugf(ctx, "No customer found - email is available")
		return nil
	}
	return fmt.Errorf("error checking if customer with an email already exists: %w", err)
}
