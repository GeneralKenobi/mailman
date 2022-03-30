package creator

import (
	"context"
	"errors"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/internal/db"
	"github.com/GeneralKenobi/mailman/internal/db/model"
	"github.com/GeneralKenobi/mailman/pkg/api/apimodel"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"time"
)

type Repository interface {
	FindCustomerByEmail(ctx context.Context, email string) (model.Customer, error)
	FindMailingEntriesByCustomerIdMailingIdTitleContentInsertTime(
		ctx context.Context, customerId, mailingId int, title, content string, insertTime time.Time) ([]model.MailingEntry, error)
	InsertMailingEntry(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error)
}

type CustomerCreator interface {
	CreateFromEmail(ctx context.Context, email string) (model.Customer, error)
}

func New(repository Repository, customerCreator CustomerCreator) *Creator {
	return &Creator{
		repository:      repository,
		customerCreator: customerCreator,
	}
}

type Creator struct {
	repository      Repository
	customerCreator CustomerCreator
}

// CreateFromDto creates a new mailing entry. It finds or creates a new user based on the email in the DTO.
// This operation is idempotent - same mailing entry can't be created twice. In that case api.StatusBadInput is returned.
func (creator *Creator) CreateFromDto(ctx context.Context, mailingEntryDto apimodel.MailingEntry) (model.MailingEntry, error) {
	customer, err := creator.getOrCreateCustomer(ctx, mailingEntryDto.Email)
	if err != nil {
		return model.MailingEntry{}, fmt.Errorf("error resolving customer for new mailing entry: %w", err)
	}

	mailingEntry := model.MailingEntry{
		CustomerId: customer.Id,
		MailingId:  mailingEntryDto.MailingId,
		Title:      mailingEntryDto.Title,
		Content:    mailingEntryDto.Content,
		InsertTime: mailingEntryDto.InsertTime,
	}

	if err = creator.assertMailingEntryDoesNotExist(ctx, mailingEntry); err != nil {
		return model.MailingEntry{}, err
	}

	return creator.Create(ctx, mailingEntry)
}

func (creator *Creator) getOrCreateCustomer(ctx context.Context, email string) (customer model.Customer, err error) {
	customer, err = creator.repository.FindCustomerByEmail(ctx, email)
	if err == nil {
		mdctx.Debugf(ctx, "Customer already exists (ID %d)", customer.Id)
		return customer, nil
	}
	if !errors.Is(err, db.ErrNoRows) {
		return model.Customer{}, fmt.Errorf("error finding customer by email: %w", err)
	}

	mdctx.Debugf(ctx, "Customer doesn't exist - creating")
	customer, err = creator.customerCreator.CreateFromEmail(ctx, email)
	if err != nil {
		return model.Customer{}, fmt.Errorf("error creating customer for new mailing entry: %w", err)
	}
	return customer, nil
}

func (creator *Creator) assertMailingEntryDoesNotExist(ctx context.Context, mailingEntry model.MailingEntry) error {
	foundEntries, err := creator.repository.FindMailingEntriesByCustomerIdMailingIdTitleContentInsertTime(
		ctx, mailingEntry.CustomerId, mailingEntry.MailingId, mailingEntry.Title, mailingEntry.Content, mailingEntry.InsertTime)
	if err != nil {
		return fmt.Errorf("error checking if mailing entry already exists: %w", err)
	}

	if len(foundEntries) == 0 {
		mdctx.Debugf(ctx, "Mailing entry doesn't exist")
		return nil
	}

	// There shouldn't be more than 1 already existing mailing entry, so logging the first one's ID should be enough information.
	mdctx.Debugf(ctx, "Found %d existing mailing entries, the first one is %d", len(foundEntries), foundEntries[0].Id)
	return api.StatusBadInput.WithMessage("this mailing entry already exists")
}

func (creator *Creator) Create(ctx context.Context, mailingEntry model.MailingEntry) (model.MailingEntry, error) {
	mdctx.Debugf(ctx, "Creating mailing entry with mailing ID %d and insert time %v for customer for customer %d",
		mailingEntry.MailingId, mailingEntry.InsertTime, mailingEntry.CustomerId)

	mailingEntry, err := creator.repository.InsertMailingEntry(ctx, mailingEntry)
	if err != nil {
		return model.MailingEntry{}, fmt.Errorf("error creating mailing entry: %w", err)
	}

	mdctx.Infof(ctx, "Created mailing entry %d", mailingEntry.Id)
	return mailingEntry, nil
}
