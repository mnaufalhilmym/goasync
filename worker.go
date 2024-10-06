package goasync

import "context"

type Worker struct {
	semaphore chan struct{}
}

// `Spawn` is equivalent to `TypedWorker[any](w).Spawn(fn)`.
func (w *Worker) Spawn(fn func(context.Context) (any, error)) JoinHandle[any] {
	return TypedWorker[any](w).Spawn(fn)
}
