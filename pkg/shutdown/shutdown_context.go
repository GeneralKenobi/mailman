package shutdown

import (
	"context"
	"github.com/GeneralKenobi/mailman/pkg/mdctx"
	"sync"
	"time"
)

// ParentContext is used for gracefully stopping worker goroutines. It is usually created by the main goroutine. It can create child
// contexts that are then passed to worker goroutines. Canceling the parent cancels every child context created by it and waits for them to
// signal that they've finished, or until the configured timeout elapses.
type ParentContext interface {
	// NewContext creates a child context. The name is only needed for logging purposes.
	NewContext(name string) Context
	// Cancel cancels each child context created by this parent. It then waits for every child to complete its graceful shutdown, or until
	// the configured timeout elapses. It can be called multiple times with no side effects.
	Cancel()
}

// Context is passed to worker goroutines. Channel obtained from Done is closed when the context is canceled and the goroutine should stop.
// The goroutines should notify about finishing graceful shutdown by calling Notify.
type Context interface {
	// Timeout returns the time allotted for gracefully shutting down. It's counted from the moment the Done channel is closed. After that
	// time there's no guarantee the parent will keep waiting for the goroutine to stop gracefully.
	Timeout() time.Duration
	// Done returns a channel that's closed when this context is canceled.
	Done() <-chan struct{}
	// Notify should be called to inform the parent that graceful shutdown has been completed. It can be called multiple times with no side
	// effects.
	Notify()
}

// NewParentContext creates a parent context which waits for all children to complete shutdown for at most shutdownTimeout.
// Child contexts are informed that they should complete their shutdown in at most shutdownTimeout.
func NewParentContext(shutdownTimeout time.Duration) ParentContext {
	return &parentContext{
		shutdownTimeout: shutdownTimeout,
	}
}

// parentContext implements ParentContext.
type parentContext struct {
	shutdownTimeout time.Duration
	children        []*childContext
	hasBeenCanceled bool // Set to true when cancel has been called for the first time
}

// Interface guard.
var _ ParentContext = (*parentContext)(nil)

func (parent *parentContext) NewContext(name string) Context {
	ctx := childContext{
		name:            name,
		shutdownTimeout: parent.shutdownTimeout,
		done:            make(chan struct{}),
		notify:          make(chan struct{}),
	}
	parent.children = append(parent.children, &ctx)
	return &ctx
}

func (parent *parentContext) Cancel() {
	parent.cancelChildren()
	parent.waitForChildren()
}

func (parent *parentContext) cancelChildren() {
	if parent.hasBeenCanceled {
		return
	}
	parent.hasBeenCanceled = true

	for _, child := range parent.children {
		mdctx.Debugf(nil, "Signaling shutdown to context %s", child.name)
		close(child.done)
	}
}

func (parent *parentContext) waitForChildren() {
	var waitGroup sync.WaitGroup
	waitGroup.Add(len(parent.children))

	shutdownTimeoutContext, cancel := context.WithTimeout(context.Background(), parent.shutdownTimeout)
	defer cancel()

	for _, child := range parent.children {
		go parent.waitForChild(shutdownTimeoutContext, &waitGroup, child)
	}
	waitGroup.Wait()
}

// waitForChild waits for child to notify that shutdown has completed or for ctx to be canceled, whichever occurs first, and subtracts
// itself from the wait group.
func (parent *parentContext) waitForChild(ctx context.Context, waitGroup *sync.WaitGroup, child *childContext) {
	mdctx.Debugf(nil, "Waiting for child %s", child.name)
	select {
	case <-child.notify:
		mdctx.Debugf(nil, "Child %s completed shutdown", child.name)
	case <-ctx.Done():
		mdctx.Debugf(nil, "Child %s didn't complete shutdown in time", child.name)
	}
	waitGroup.Done()
}

// childContext implements Context.
type childContext struct {
	name            string
	shutdownTimeout time.Duration
	done            chan struct{}
	notify          chan struct{}
	hasNotified     bool
}

// Interface guard.
var _ Context = (*childContext)(nil)

func (ctx *childContext) Timeout() time.Duration {
	return ctx.shutdownTimeout
}

func (ctx *childContext) Done() <-chan struct{} {
	return ctx.done
}

func (ctx *childContext) Notify() {
	if ctx.hasNotified {
		return
	}
	ctx.hasNotified = true

	close(ctx.notify)
}
