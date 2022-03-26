package staleremover

import (
	"context"
	"fmt"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/internal/persistence/model"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"time"
)

type Repository interface {
	FindMailingEntriesOlderThan(ctx context.Context, olderThan time.Duration) ([]model.MailingEntry, error)
	FindMailingEntriesByMailingIdOlderThan(ctx context.Context, mailingId int, olderThan time.Duration) ([]model.MailingEntry, error)
	FindMailingEntriesByCustomerId(ctx context.Context, customerId int) ([]model.MailingEntry, error)
	DeleteMailingEntryById(ctx context.Context, id int) error
	DeleteCustomerById(ctx context.Context, id int) error
}

type StaleEntryRemover struct {
	repository Repository
}

func New(repository Repository) *StaleEntryRemover {
	return &StaleEntryRemover{repository: repository}
}

// RemoveByMailingId finds and removes all mailing entries with the given mailing ID that are older than the configured threshold.
func (remover *StaleEntryRemover) RemoveByMailingId(ctx context.Context, mailingId int) error {
	staleEntries, err := remover.repository.FindMailingEntriesByMailingIdOlderThan(ctx, mailingId, stalenessThreshold())
	if err != nil {
		return err
	}
	return remover.removeStaleEntries(ctx, staleEntries)
}

// Remove finds and removes all mailing entries that are older than the configured threshold.
func (remover *StaleEntryRemover) Remove(ctx context.Context) error {
	staleEntries, err := remover.repository.FindMailingEntriesOlderThan(ctx, stalenessThreshold())
	if err != nil {
		return err
	}
	return remover.removeStaleEntries(ctx, staleEntries)
}

func (remover *StaleEntryRemover) removeStaleEntries(ctx context.Context, staleEntries []model.MailingEntry) error {
	mdctx.Infof(ctx, "Removing %d stale mailing entries", len(staleEntries))

	for _, entry := range staleEntries {
		mdctx.Infof(ctx, "Removing stale mailing entry %d", entry.Id)
		err := remover.repository.DeleteMailingEntryById(ctx, entry.Id)
		if err != nil {
			return fmt.Errorf("error removing stale entry %d: %w", entry.Id, err)
		}
	}

	err := remover.removeNoLongerReferencedCustomers(ctx, staleEntries)
	if err != nil {
		return fmt.Errorf("error cleaning up customers after removing stale mailing entries: %w", err)
	}
	return err
}

func (remover *StaleEntryRemover) removeNoLongerReferencedCustomers(ctx context.Context, removedEntries []model.MailingEntry) error {
	uniqueCustomerIds := make(map[int]bool)
	for _, entry := range removedEntries {
		uniqueCustomerIds[entry.MailingId] = true
	}

	for customerId := range uniqueCustomerIds {
		entriesForCustomer, err := remover.repository.FindMailingEntriesByCustomerId(ctx, customerId)
		if err != nil {
			return fmt.Errorf("error listing mailing entries for customer %d: %err", customerId, err)
		}
		if len(entriesForCustomer) > 0 {
			mdctx.Debugf(ctx, "Not cleaning up customer %d - they still have %d mailing entries", customerId, len(entriesForCustomer))
			continue
		}

		mdctx.Infof(ctx, "Removing customer %d - no more mailing entries are associated with them", customerId)
		err = remover.repository.DeleteCustomerById(ctx, customerId)
		if err != nil {
			return fmt.Errorf("error cleaning up customer %d: %w", customerId, err)
		}
	}

	return nil
}

// Hook for mocking in unit tests.
var stalenessThreshold = func() time.Duration {
	return time.Duration(config.Get().StaleMailingEntryRemover.StalenessThresholdSeconds) * time.Second
}
