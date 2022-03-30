package scheduler

import (
	"context"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"github.com/GeneralKenobi/mailman/pkg/shutdown"
	"time"
)

// New creates a scheduler for a function.
// operationName is used for logging and enhancing MDC in contexts passed to the scheduled function.
//
// Scheduler recovers from scheduled function panics and logs errors from its execution.
func New(operationName string, todo func(ctx context.Context) error) *Scheduler {
	return &Scheduler{
		operationName: operationName,
		todo:          todo,
	}
}

type Scheduler struct {
	operationName string
	todo          func(ctx context.Context) error
}

// RunPeriodically runs this scheduler's function every period until the context is canceled. The first execution occurs after period has
// elapsed.
func (scheduler *Scheduler) RunPeriodically(ctx shutdown.Context, period time.Duration) {
	defer ctx.Notify()
	mdctx.Infof(nil, "Starting scheduled execution of %q every %v", scheduler.operationName, period)

	ticker := time.NewTicker(period)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			scheduler.execute()
		case <-ctx.Done():
			mdctx.Infof(nil, "Context canceled - stopping scheduled execution of %q", scheduler.operationName)
			return
		}
	}
}

func (scheduler *Scheduler) execute() {
	ctx := mdctx.New()
	ctx = mdctx.WithOperationName(ctx, scheduler.operationName)
	mdctx.Debugf(ctx, "Starting scheduled execution")

	err := scheduler.todo(ctx)
	if panicErr := recover(); panicErr != nil {
		mdctx.Errorf(ctx, "Recovered from panic in scheduled execution: %v", panicErr)
		return
	}
	if err != nil {
		mdctx.Errorf(ctx, "Scheduled execution ended with error: %v", err)
		return
	}
	mdctx.Debugf(ctx, "Scheduled execution completed with success")
}
