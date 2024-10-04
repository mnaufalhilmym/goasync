package goasync

import (
	"context"

	"github.com/mnaufalhilmym/goresult"
)

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
	doneCh := make(chan struct{}, 1)
	var result goresult.Result[T]

	go func() {
		res, err := fn(ctx)
		cancel()
		if err != nil {
			result = goresult.Err[T](err)
		} else {
			result = goresult.Ok(res)
		}
		close(doneCh)
	}()

	return JoinHandle[T]{
		doneCh: doneCh,
		result: &result,
		cancel: cancel,
	}
}
