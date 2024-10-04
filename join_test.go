package goasync_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/mnaufalhilmym/goasync"
)

func TestAwait(t *testing.T) {
	t.Log("Running TestAwait. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	task := goasync.Spawn(fn)
	res, err := task.Await(context.Background())
	if err != nil {
		t.Error("TestAwait failed. Error should be nil. Error:", err)
	}
	if res != true {
		t.Error("TestAwait failed. Result must be true")
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestAwait should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestMultipleAwait(t *testing.T) {
	t.Log("Running TestMultipleAwait. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	task := goasync.Spawn(fn)
	for i := 0; i < 10; i++ {
		res, err := task.Await(context.Background())
		if err != nil {
			t.Error("TestMultipleAwait failed. Error should be nil. Error:", err)
		}
		if res != true {
			t.Error("TestMultipleAwait failed. Result must be true")
		}
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestMultipleAwait should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestMultipleAwaitInGoroutine(t *testing.T) {
	t.Log("Running TestMultipleAwaitInGoroutine. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	task := goasync.Spawn(fn)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			res, err := task.Await(context.Background())
			if err != nil {
				t.Error("TestMultipleAwaitInGoroutine failed. Error should be nil. Error:", err)
			}
			if res != true {
				t.Error("TestMultipleAwaitInGoroutine failed. Result must be true")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestMultipleAwaitInGoroutine should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestAwaitAfterSleep(t *testing.T) {
	t.Log("Running TestAwaitAfterSleep. Expected to complete in about 2 seconds.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	task := goasync.Spawn(fn)
	time.Sleep(2 * time.Second)
	res, err := task.Await(context.Background())
	if err != nil {
		t.Error("TestAwaitAfterSleep failed. Error should be nil. Error:", err)
	}
	if res != true {
		t.Error("TestAwaitAfterSleep failed. Result must be true")
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 2 {
		t.Errorf("TestAwaitAfterSleep should complete in 2 second. Not in %d seconds.", duration)
	}
}

func TestAwaitWithTimeoutCtx(t *testing.T) {
	t.Log("Running TestAwaitWithTimeoutCtx. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(ctx context.Context) (bool, error) {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(2 * time.Second):
			return true, nil
		}
	}

	task := goasync.Spawn(fn)
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	result, err := task.Await(ctx)
	if err == nil {
		t.Error("TestAwaitWithTimeoutCtx failed. Error should be not nil")
	}
	if result != false {
		t.Error("TestAwaitWithTimeoutCtx failed. Result should be false")
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestAwaitWithTimeoutCtx should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestAbort(t *testing.T) {
	t.Log("Running TestAbort. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(ctx context.Context) (bool, error) {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(2 * time.Second):
			return true, nil
		}
	}

	task := goasync.Spawn(fn)
	time.Sleep(1 * time.Second)
	task.Abort()
	result, err := task.Await(context.Background())
	if err == nil {
		t.Error("TestAbort failed. Error should be not nil")
	}
	if result != false {
		t.Error("TestAbort failed. Result should be false")
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestAbort should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestMultipleAbortInGoroutine(t *testing.T) {
	t.Log("Running TestMultipleAbortInGoroutine. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(ctx context.Context) (bool, error) {
		select {
		case <-ctx.Done():
			return false, ctx.Err()
		case <-time.After(2 * time.Second):
			return true, nil
		}
	}

	task := goasync.Spawn(fn)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			time.Sleep(1 * time.Second)
			task.Abort()
			result, err := task.Await(context.Background())
			if err == nil {
				t.Error("TestMultipleAbortInGoroutine failed. Error should be not nil")
			}
			if result != false {
				t.Error("TestMultipleAbortInGoroutine failed. Result should be false")
			}
			wg.Done()
		}()
	}
	wg.Wait()

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestMultipleAbortInGoroutine should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestIsFinished(t *testing.T) {
	t.Log("Running TestIsFinished. Expected to complete in about 1 second.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	task := goasync.Spawn(fn)
	if task.IsFinished() {
		t.Error("TestIsFinished failed. task.IsFinished() should return false")
	}
	task.Await(context.Background())
	if !task.IsFinished() {
		t.Error("TestIsFinished failed. task.IsFinished() should return true")
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 1 {
		t.Errorf("TestIsFinished should complete in 1 second. Not in %d seconds.", duration)
	}
}

func TestIsFinishedAfterSleepWithoutAwait(t *testing.T) {
	t.Log("Running TestIsFinishedAfterSleepWithoutAwait. Expected to complete in about 2 seconds.")

	start := time.Now()

	fn := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	task := goasync.Spawn(fn)
	if task.IsFinished() {
		t.Error("TestIsFinishedAfterSleepWithoutAwait failed. task.IsFinished() should return false")
	}
	time.Sleep(2 * time.Second)
	if !task.IsFinished() {
		t.Error("TestIsFinishedAfterSleepWithoutAwait failed. task.IsFinished() should return true")
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 2 {
		t.Errorf("TestIsFinishedAfterSleepWithoutAwait should complete in 2 seconds. Not in %d seconds.", duration)
	}
}
