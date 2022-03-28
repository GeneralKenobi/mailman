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
	DeleteMailingEntryById(ctx context.Context, id int) error
}

func New(repository Repository) *StaleEntryRemover {
	return &StaleEntryRemover{repository: repository}
}

type StaleEntryRemover struct {
	repository Repository
}

// RemoveByMailingId finds and removes all mailing entries with the given mailing ID that are older than the configured threshold.
func (remover *StaleEntryRemover) RemoveByMailingId(ctx context.Context, mailingId int) error {
	staleEntries, err := remover.repository.FindMailingEntriesByMailingIdOlderThan(ctx, mailingId, stalenessThreshold())
	if err != nil {
		return fmt.Errorf("error listing stale mailing entries with mailing ID %d: %w", mailingId, err)
	}

	return remover.removeStaleEntries(ctx, staleEntries)
}

// Remove finds and removes all mailing entries that are older than the configured threshold.
func (remover *StaleEntryRemover) Remove(ctx context.Context) error {
	staleEntries, err := remover.repository.FindMailingEntriesOlderThan(ctx, stalenessThreshold())
	if err != nil {
		return fmt.Errorf("error listing stale mailing entries: %w", err)
	}

	return remover.removeStaleEntries(ctx, staleEntries)
}

func (remover *StaleEntryRemover) removeStaleEntries(ctx context.Context, staleEntries []model.MailingEntry) error {
	mdctx.Infof(ctx, "Removing %d stale mailing entries", len(staleEntries))

	for _, entry := range staleEntries {
		mdctx.Infof(ctx, "Removing stale mailing entry %d", entry.Id)
		err := remover.repository.DeleteMailingEntryById(ctx, entry.Id)
		if err != nil {
			return fmt.Errorf("error removing stale mailing entry %d: %w", entry.Id, err)
		}
	}

	return nil
}

// Hook for mocking in unit tests.
var stalenessThreshold = func() time.Duration {
	return time.Duration(config.Get().StaleMailingEntryRemover.StalenessThresholdSeconds) * time.Second
}
