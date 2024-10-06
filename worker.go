package goasync

import (
	"context"
	"runtime"

	"github.com/mnaufalhilmym/goresult"
)

type WorkerBuilder struct {
	maxWorkers uint
}

// `NewWorkerBuilder` creates a new `WorkerBuilder` instance with the default
// maximum number of workers set to the number of available CPU
// cores.
func NewWorkerBuilder() *WorkerBuilder {
	return &WorkerBuilder{
		maxWorkers: uint(runtime.NumCPU()),
	}
}

// `SetMaxWorkers` allows the user to specify a custom maximum number
// of concurrent workers.
func (w *WorkerBuilder) SetMaxWorkers(maxWorkers uint) *WorkerBuilder {
	w.maxWorkers = maxWorkers
	return w
}

// `Build` creates a new `Worker` instance with the configured maximum
// number of workers.
func (w *WorkerBuilder) Build() *Worker {
	return &Worker{
		semaphore: make(chan struct{}, w.maxWorkers),
	}
}

type Worker struct {
	semaphore chan struct{}
}

// `Spawn` is equivalent to `TypedWorker[any](w).Spawn(fn)`.
func (w *Worker) Spawn(fn func(context.Context) (any, error)) JoinHandle[any] {
	return TypedWorker[any](w).Spawn(fn)
}

// `TypedWorker` creates a type-safe worker for spawning tasks that return
// results of type T.
func TypedWorker[T any](w *Worker) *WorkerTyped[T] {
	return &WorkerTyped[T]{w}
}

type WorkerTyped[T any] struct {
	w *Worker
}

// Spawns a new asynchronous task, returning a
// `JoinHandle[T any]` for it.
//
// `Spawn` will block if all workers are currently in use. If any worker
// becomes available, the call will be unblocked, and the function will
// be executed in a new goroutine.
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
func (wt *WorkerTyped[T]) Spawn(fn func(context.Context) (T, error)) JoinHandle[T] {
	wt.w.semaphore <- struct{}{}

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
		<-wt.w.semaphore
	}()

	return JoinHandle[T]{
		doneCh: doneCh,
		result: &result,
		cancel: cancel,
	}
}
