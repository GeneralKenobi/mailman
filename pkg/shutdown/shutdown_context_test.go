package shutdown

import (
	"testing"
	"time"
)

func TestShutdownContext(t *testing.T) {
	parent := NewParentContext(10 * time.Second)
	ctx1 := parent.NewContext("first")
	ctx2 := parent.NewContext("second")
	ctx3 := parent.NewContext("third")
	go simulateWorkWithGracefulShutdown(ctx1, time.Second)
	go simulateWorkWithGracefulShutdown(ctx2, time.Second+300*time.Millisecond)
	go simulateWorkWithGracefulShutdown(ctx3, time.Second+900*time.Millisecond)

	startTime := time.Now()
	parent.Cancel()
	endTime := time.Now()

	// The longest shutdown takes 1.9s, so 3s should be more than enough for each worker to complete shutdown.
	const shouldCompleteIn = 3 * time.Second
	timeTaken := endTime.Sub(startTime)
	if timeTaken > shouldCompleteIn {
		t.Errorf("Expected the wait to take at most %v but it was %v", shouldCompleteIn, timeTaken)
	}
}

func TestShutdownContextShouldTimeOutIfChildDoesNotStop(t *testing.T) {
	parent := NewParentContext(3 * time.Second)
	ctxStopsInTime := parent.NewContext("stops in time")
	ctxDoesNotStopInTime := parent.NewContext("does not stop in time")
	go simulateWorkWithGracefulShutdown(ctxStopsInTime, time.Second)
	go simulateWorkWithGracefulShutdown(ctxDoesNotStopInTime, 10*time.Second)

	startTime := time.Now()
	parent.Cancel()
	endTime := time.Now()

	// Timeout is set to 3s so 4s should be more than enough for timeout to be hit and wait to exit.
	const shouldCompleteIn = 4 * time.Second
	timeTaken := endTime.Sub(startTime)
	if timeTaken > shouldCompleteIn {
		t.Errorf("Expected the wait to take at most %v but it was %v", shouldCompleteIn, timeTaken)
	}
}

func TestParentContextShouldAllowMultipleCancelCalls(t *testing.T) {
	parent := NewParentContext(3 * time.Second)
	ctx := parent.NewContext("test")
	go simulateWorkWithGracefulShutdown(ctx, 0)

	parent.Cancel()
	parent.Cancel()
}

func TestChildContextShouldAllowMultipleNotifyCalls(t *testing.T) {
	parent := NewParentContext(3 * time.Second)
	ctx := parent.NewContext("test")

	ctx.Notify()
	ctx.Notify()
}

func simulateWorkWithGracefulShutdown(ctx Context, simulatedShutdownDuration time.Duration) {
	<-ctx.Done()
	time.Sleep(simulatedShutdownDuration)
	ctx.Notify()
}
