package goasync

import (
	"context"
	"sync"

	"github.com/mnaufalhilmym/goresult"
)

type JoinHandle[T any] struct {
	sync.Mutex
	fn         func(context.Context) (T, error)
	resultCh   chan goresult.Result[T]
	result     goresult.Result[T]
	isFinished *bool
	cancel     context.CancelFunc
}

func (h *JoinHandle[T]) Abort() {
	h.cancel()
}

func (h *JoinHandle[T]) Await(ctx context.Context) (T, error) {
	for {
		if h.result != nil {
			if h.result.IsOk() {
				return h.result.Unwrap(), nil
			} else {
				var empty T
				return empty, h.result.UnwrapErr()
			}
		}

		select {
		case result, ok := <-h.resultCh:
			if !ok && h.result == nil {
				continue
			}
			if ok {
				h.Lock()
				h.result = result
				h.Unlock()
			}
		case <-ctx.Done():
			var empty T
			return empty, ctx.Err()
		}
	}
}

func (h *JoinHandle[T]) IsFinished() bool {
	return *h.isFinished
}

func Spawn[T any](fn func(context.Context) (T, error)) *JoinHandle[T] {
	ctx, cancel := context.WithCancel(context.Background())
	resultCh := make(chan goresult.Result[T], 1)
	isFinished := false

	go func() {
		res, err := fn(ctx)
		cancel()
		if err != nil {
			resultCh <- goresult.NewErr[T](err)
		} else {
			resultCh <- goresult.NewOk(res)
		}
		close(resultCh)
		isFinished = true
	}()

	return &JoinHandle[T]{
		fn:         fn,
		resultCh:   resultCh,
		isFinished: &isFinished,
		cancel:     cancel,
	}
}

func TryJoin[T any](ctx context.Context, handles ...*JoinHandle[T]) ([]T, error) {
	countHandles := len(handles)

	resultsCh := make(chan []any, countHandles)
	errCh := make(chan error, 1)

	for i := range handles {
		idx := i
		go func() {
			res, err := handles[idx].Await(ctx)
			if err != nil {
				errCh <- err
				return
			}
			resultsCh <- []any{idx, res}
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
