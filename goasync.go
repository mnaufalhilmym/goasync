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

func (h *JoinHandle[T]) Abort() {
	h.cancel()
}

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

func (h *JoinHandle[T]) IsFinished() bool {
	return (*h.result) != nil
}

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
