package goasync

import "runtime"

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
