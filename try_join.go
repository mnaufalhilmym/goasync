package goasync

import "context"

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
