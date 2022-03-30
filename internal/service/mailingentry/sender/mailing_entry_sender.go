package sender

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/api"
	"github.com/GeneralKenobi/mailman/internal/db/model"
	"github.com/GeneralKenobi/mailman/pkg/api/apimodel"
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

func New(repository Repository, emailer Emailer) *EntrySender {
	return &EntrySender{
		repository: repository,
		emailer:    emailer,
	}
}

type EntrySender struct {
	repository Repository
	emailer    Emailer
}

// SendMailingRequest sends email for every mailing entry with the given mailing ID and deletes them from the database.
func (sender *EntrySender) SendMailingRequest(ctx context.Context, mailingRequest apimodel.MailingRequest) (err error) {
	entries, err := sender.repository.FindMailingEntriesByMailingId(ctx, mailingRequest.MailingId)
	if err != nil {
		return fmt.Errorf("error listing mailing entries: %w", err)
	}

	mdctx.Debugf(ctx, "Found %d mailing entries with mailing ID %d", len(entries), mailingRequest.MailingId)
	if len(entries) == 0 {
		return api.StatusNotFound.WithMessage("no mailing entries to send")
	}

	for _, entry := range entries {
		err = sender.Send(ctx, entry)
		if err != nil {
			return err
		}
	}

	return nil
}

// Send sends the mailing entry and deletes it from the database.
func (sender *EntrySender) Send(ctx context.Context, mailingEntry model.MailingEntry) error {
	customer, err := sender.repository.FindCustomerById(ctx, mailingEntry.CustomerId)
	if err != nil {
		return fmt.Errorf("error finding customer %d for mailing entry %d: %w", mailingEntry.CustomerId, mailingEntry.Id, err)
	}

	mdctx.Debugf(ctx, "Sending mailing entry with ID %d", mailingEntry.Id)
	err = sender.emailer.Send(ctx, customer.Email, mailingEntry.Title, mailingEntry.Content)
	if err != nil {
		return fmt.Errorf("error sending mailing entry %d to customer %d: %w", mailingEntry.Id, mailingEntry.CustomerId, err)
	}

	mdctx.Infof(ctx, "Deleting mailing entry with ID %d", mailingEntry.Id)
	err = sender.repository.DeleteMailingEntryById(ctx, mailingEntry.Id)
	if err != nil {
		return fmt.Errorf("error deleting mailing entry %d: %w", mailingEntry.Id, err)
	}

	return nil
}
