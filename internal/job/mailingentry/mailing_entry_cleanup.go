package mailingentry

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/internal/db"
	"github.com/GeneralKenobi/mailman/internal/service/mailingentry/staleremover"
	"github.com/GeneralKenobi/mailman/pkg/scheduler"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	"time"
)

func NewCleanupJob(transactioner db.Transactioner) *CleanupJob {
	return &CleanupJob{transactioner: transactioner}
}

type CleanupJob struct {
	transactioner db.Transactioner
}

// RunScheduled runs stale entry cleanup periodically until the context is canceled.
func (cleanupJob *CleanupJob) RunScheduled(ctx shutdown.Context) {
	jobScheduler := scheduler.New("stale mailing entry cleanup", cleanupJob.RunCleanup)
	jobScheduler.RunPeriodically(ctx, schedulingPeriod())
}

func (cleanupJob *CleanupJob) RunCleanup(ctx context.Context) error {
	return db.InTransaction(ctx, cleanupJob.transactioner, func(repository db.Repository) error {
		staleMailingEntryRemover := staleremover.New(repository)
		return staleMailingEntryRemover.Remove(ctx)
	})
}

// Hook for mocking in unit tests.
var schedulingPeriod = func() time.Duration {
	return time.Duration(config.Get().MailingEntryCleanupJob.PeriodSeconds) * time.Second
}
