package goasync_test

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"

	"github.com/mnaufalhilmym/goasync"
)

func TestTryJoin(t *testing.T) {
	t.Log("Running TestTryJoin. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	tasks := make([]goasync.JoinHandle[bool], 10)
	for i := range tasks {
		tasks[i] = goasync.Spawn(fn)
	}

	results, err := goasync.TryJoin(context.Background(), tasks...)
	if err != nil {
		t.Error("TestTryJoin failed. Error should be nil. Error:", err)
	}

	for _, result := range results {
		if result != true {
			t.Error("TestTryJoin failed. All results should be true.", result)
		}
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestTryJoin should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestMultipleTryJoinInGoroutine(t *testing.T) {
	t.Log("Running TestMultipleTryJoinInGoroutine. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	tasks := make([]goasync.JoinHandle[bool], 10)
	for i := range tasks {
		tasks[i] = goasync.Spawn(fn)
	}

	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			results, err := goasync.TryJoin(context.Background(), tasks...)
			if err != nil {
				t.Error("TestMultipleTryJoinInGoroutine failed. Error should be nil. Error:", err)
			}

			for _, result := range results {
				if result != true {
					t.Error("TestMultipleTryJoinInGoroutine failed. All results should be true.", result)
				}
			}
			wg.Done()
		}()
	}
	wg.Wait()

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestMultipleTryJoinInGoroutine should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestTryJoinWithTimeoutCtx(t *testing.T) {
	t.Log("Running TestTryJoinWithTimeoutCtx. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(2 * time.Second)
		return true, nil
	}

	tasks := make([]goasync.JoinHandle[bool], 10)
	for i := range tasks {
		tasks[i] = goasync.Spawn(fn)
	}

	ctx, cancel := context.WithTimeout(context.TODO(), 1*time.Second)
	defer cancel()
	results, err := goasync.TryJoin(ctx, tasks...)
	if err == nil {
		t.Error("TestTryJoinWithTimeoutCtx failed. Error should be not nil.")
	}

	for _, result := range results {
		if result != false {
			t.Error("TestTryJoinWithTimeoutCtx failed. All results should be false.", result)
		}
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestTryJoinWithTimeoutCtx should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestTryJoinWithErrorTask(t *testing.T) {
	t.Log("Running TestTryJoinWithErrorTask. Expected to complete in less than 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}
	fnError := func(context.Context) (bool, error) {
		return false, errors.New("dummy error")
	}

	tasks := make([]goasync.JoinHandle[bool], 10)
	for i := range tasks {
		if i == 5 {
			tasks[i] = goasync.Spawn(fnError)
		} else {
			tasks[i] = goasync.Spawn(fn)
		}
	}

	_, err := goasync.TryJoin(context.Background(), tasks...)
	if err == nil {
		t.Error("TestTryJoinWithErrorTask failed. Error should be not nil.")
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); !(duration < 1) {
		t.Errorf("TestTryJoinWithErrorTask should complete in less than 1 second. Not in %d seconds.", duration)
	}
}
