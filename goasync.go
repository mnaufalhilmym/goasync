package goasync

import (
	"context"

	"github.com/mnaufalhilmym/goresult"
)

type JoinHandle[T any] struct {
	doneCh chan byte
	result *goresult.Result[T]
	cancel context.CancelFunc
}

// Abort the task associated with the handle.
//
// Awaiting a cancelled task might complete as usual if the task was
// already completed at the time it was cancelled, but most likely it
// will fail with a `context canceled`.
func (h *JoinHandle[T]) Abort() {
	h.cancel()
}

// Suspend execution until the result of a `JoinHandle` is ready.
//
// `.Await`ing a future will suspend the current function's execution
// until the executor has run the future to completion.
func (h *JoinHandle[T]) Await(ctx context.Context) (T, error) {
	for {
		if *h.result != nil {
			if (*h.result).IsOk() {
				return (*h.result).Unwrap(), nil
			} else {
				var empty T
				return empty, (*h.result).UnwrapErr()
			}
		}

		select {
		case <-h.doneCh:
			continue
		case <-ctx.Done():
			var empty T
			return empty, ctx.Err()
		}
	}
}

// Checks if the task associated with this `JoinHandle` has finished.
//
// Please note that this method can return `false` even if `Abort` has been
// called on the task. This is because the cancellation process may take
// some time, and this method does not return `true` until it has
// completed.
func (h *JoinHandle[T]) IsFinished() bool {
	return (*h.result) != nil
}

// Spawns a new asynchronous task, returning a
// `JoinHandle[T any]` for it.
//
// The provided future will start running in the background immediately
// when `Spawn` is called, even if you don't await the returned
// `JoinHandle`.
//
// Spawning a task enables the task to execute concurrently to other tasks. The
// spawned task may execute on the current thread, or it may be sent to a
// different thread to be executed. The specifics depend on how the Go runtime
// schedules it.
//
// There is no guarantee that a spawned task will execute to completion.
// When a runtime is shutdown, all outstanding tasks are dropped,
// regardless of the lifecycle of that task.
func Spawn[T any](fn func(context.Context) (T, error)) JoinHandle[T] {
	ctx, cancel := context.WithCancel(context.Background())
	doneCh := make(chan byte, 1)
	var result goresult.Result[T]

	go func() {
		res, err := fn(ctx)
		cancel()
		if err != nil {
			result = goresult.Err[T](err)
		} else {
			result = goresult.Ok(res)
		}
		doneCh <- 1
		close(doneCh)
	}()

	return JoinHandle[T]{
		doneCh: doneCh,
		result: &result,
		cancel: cancel,
	}
}

// Waits on multiple concurrent branches, returning when all branches
// complete with `success` value or on the first `error`.
//
// The `TryJoin` function takes a list of `JoinHandle` and evaluates them
// concurrently.`.
func TryJoin[T any](ctx context.Context, handles ...JoinHandle[T]) ([]T, error) {
	countHandles := len(handles)

	resultsCh := make(chan []any, countHandles)
	errCh := make(chan error, 1)

	for i := range handles {
		idx := i
		go func() {
			res, err := handles[idx].Await(ctx)
			if err != nil {
				errCh <- err
			} else {
				resultsCh <- []any{idx, res}
			}
		}()
	}

	results := make([]T, countHandles)

	for range handles {
		select {
		case result := <-resultsCh:
			res, ok := result[1].(T)
			if ok {
				results[result[0].(int)] = res
			}
		case err := <-errCh:
			for _, handle := range handles {
				handle.cancel()
			}
			return nil, err
		case <-ctx.Done():
			for _, handle := range handles {
				handle.cancel()
			}
			return nil, ctx.Err()
		}
	}

	return results, nil
}
