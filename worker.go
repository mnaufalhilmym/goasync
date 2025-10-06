package goasync

import (
	"context"
	"runtime"
)

type Worker struct {
	semaphore chan struct{}
}

type WorkerConfig struct {
	MaxWorkers uint
}

// `NewWorker` creates a new `Worker` instance with an optional configuration.
//
// The function initializes the internal semaphore based on the configured maximum
// number of workers. If zero or one `WorkerConfig` is provided, the settings
// are applied. If no configuration is provided, the `MaxWorkers` defaults to
// the number of available CPU cores (`runtime.NumCPU()`). If a configuration
// is provided, `MaxWorkers` overrides this default.
func NewWorker(cfgs ...WorkerConfig) *Worker {
	maxWorkers := uint(runtime.NumCPU())

	if len(cfgs) > 0 {
		cfg := cfgs[0]
		maxWorkers = cfg.MaxWorkers
	}

	return &Worker{
		semaphore: make(chan struct{}, maxWorkers),
	}
}

// `Spawn` is equivalent to `TypedWorker[any](w).Spawn(fn)`.
func (w *Worker) Spawn(fn func(context.Context) (any, error)) JoinHandle[any] {
	return TypedWorker[any](w).Spawn(fn)
}
