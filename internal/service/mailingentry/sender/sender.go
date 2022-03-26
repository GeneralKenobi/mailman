package sender

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/internal/persistence/model"
	apimodel "github.com/GeneralKenobi/mailman/pkg/api/model"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
)

type Repository interface {
	FindMailingEntriesByMailingId(ctx context.Context, mailingId int) ([]model.MailingEntry, error)
	FindCustomerById(ctx context.Context, id int) (model.Customer, error)
	DeleteMailingEntryById(ctx context.Context, id int) error
}

type Emailer interface {
	Send(ctx context.Context, emailAddress, title, content string) error
}

type EntrySender struct {
	repository Repository
	emailer    Emailer
}

func New(repository Repository, emailer Emailer) *EntrySender {
	return &EntrySender{
		repository: repository,
		emailer:    emailer,
	}
}

// Send sends email for every mailing entry with the given mailing ID and deletes them from the database.
func (sender *EntrySender) Send(ctx context.Context, mailingRequest apimodel.MailingRequestDto) (err error) {
	entries, err := sender.repository.FindMailingEntriesByMailingId(ctx, mailingRequest.MailingId)
	if err != nil {
		return fmt.Errorf("error listing mailing entries: %w", err)
	}

	mdctx.Debugf(ctx, "Found %d mailing entries with mailing ID %d", len(entries), mailingRequest.MailingId)
	if len(entries) == 0 {
		return api.StatusNotFound.Error("no mailing entries to send")
	}

	for _, entry := range entries {
		customer, err := sender.repository.FindCustomerById(ctx, entry.CustomerId)
		if err != nil {
			return fmt.Errorf("can't find customer with ID %d: %w", entry.CustomerId, err)
		}

		mdctx.Debugf(ctx, "Sending mailing entry with ID %d", entry.Id)
		err = sender.emailer.Send(ctx, customer.Email, entry.Title, entry.Content)
		if err != nil {
			return fmt.Errorf("error sending mailing entry %d to customer %d: %w", entry.Id, entry.CustomerId, err)
		}

		mdctx.Infof(ctx, "Deleting mailing entry with ID %d", entry.Id)
		err = sender.repository.DeleteMailingEntryById(ctx, entry.Id)
		if err != nil {
			return fmt.Errorf("error deleting mailing entry with ID %d: %w", entry.Id, err)
		}
	}

	return nil
}
