package goasync

import (
	"context"

	"github.com/mnaufalhilmym/goresult"
)

type JoinHandle[T any] struct {
	doneCh chan struct{}
	result *goresult.Result[T]
	cancel context.CancelFunc
}

// Abort the task associated with the handle.
//
// Awaiting a cancelled task might complete as usual if the task was
// already completed at the time it was cancelled, but most likely it
// will fail with a `context canceled`.
func (h JoinHandle[T]) Abort() {
	h.cancel()
}

// Suspend execution until the result of a `JoinHandle` is ready.
//
// `.Await`ing a future will suspend the current function's execution
// until the executor has run the future to completion.
func (h JoinHandle[T]) Await(ctx context.Context) (T, error) {
	select {
	case <-ctx.Done():
		h.cancel()
		var empty T
		return empty, ctx.Err()
	case <-h.doneCh:
	}

	if (*h.result).IsOk() {
		return (*h.result).Unwrap(), nil
	} else {
		var empty T
		return empty, (*h.result).UnwrapErr()
	}
}

// Checks if the task associated with this `JoinHandle` has finished.
//
// Please note that this method can return `false` even if `Abort` has been
// called on the task. This is because the cancellation process may take
// some time, and this method does not return `true` until it has
// completed.
func (h JoinHandle[T]) IsFinished() bool {
	return (*h.result) != nil
}
