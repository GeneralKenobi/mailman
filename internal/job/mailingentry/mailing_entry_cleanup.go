package mailingentry

import (
	"context"
	"github.com/GeneralKenobi/mailman/internal/config"
	"github.com/GeneralKenobi/mailman/internal/job"
	"github.com/GeneralKenobi/mailman/internal/persistence"
	"github.com/GeneralKenobi/mailman/internal/service/mailingentry/staleremover"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	"time"
)

func NewCleanupJob(transactioner persistence.Transactioner) *CleanupJob {
	return &CleanupJob{transactioner: transactioner}
}

type CleanupJob struct {
	transactioner persistence.Transactioner
}

// RunScheduled runs stale entry cleanup periodically until the context is canceled.
func (cleanupJob *CleanupJob) RunScheduled(ctx shutdown.Context) {
	scheduler := job.NewScheduler("stale mailing entry cleanup", cleanupJob.RunCleanup)
	scheduler.RunPeriodically(ctx, schedulingPeriod())
}

func (cleanupJob *CleanupJob) RunCleanup(ctx context.Context) error {
	return persistence.WithinTransaction(ctx, cleanupJob.transactioner, func(transactionalRepository persistence.Repository) error {
		staleMailingEntryRemover := staleremover.New(transactionalRepository)
		return staleMailingEntryRemover.Remove(ctx)
	})
}

// Hook for mocking in unit tests.
var schedulingPeriod = func() time.Duration {
	return time.Duration(config.Get().MailingEntryCleanupJob.PeriodSeconds) * time.Second
}
