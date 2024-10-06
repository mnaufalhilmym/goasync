package goasync_test

import (
	"context"
	"testing"
	"time"

	"github.com/mnaufalhilmym/goasync"
)

func TestWorker(t *testing.T) {
	t.Log("Running TestWorker. Expected to complete in about 2 seconds.")

	start := time.Now()

	worker := goasync.
		NewWorkerBuilder().
		SetMaxWorkers(3).
		Build()

	fn := func(context.Context) (any, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}

	tasks := make([]goasync.JoinHandle[any], 6)
	for i := 0; i < len(tasks); i++ {
		tasks[i] = worker.Spawn(fn)
	}

	for _, task := range tasks {
		res, err := task.Await(context.Background())
		if err != nil {
			t.Error("TestWorker failed. Error should be nil. Error:", err)
		}
		if res != true {
			t.Error("TestWorker failed. Result must be true. Result:", res)
		}
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 2 {
		t.Errorf("TestWorker should complete in 2 seconds. Not in %d seconds.", duration)
	}
}

func TestTypedWorker(t *testing.T) {
	t.Log("Running TestTypedWorker. Expected to complete in about 2 seconds.")

	start := time.Now()

	worker := goasync.
		NewWorkerBuilder().
		SetMaxWorkers(3).
		Build()

	workerBool := goasync.TypedWorker[bool](worker)
	workerInt := goasync.TypedWorker[int](worker)

	fn1 := func(context.Context) (bool, error) {
		time.Sleep(1 * time.Second)
		return true, nil
	}
	fn2 := func(context.Context) (int, error) {
		time.Sleep(1 * time.Second)
		return 1, nil
	}

	tasks1 := make([]goasync.JoinHandle[bool], 3)
	tasks2 := make([]goasync.JoinHandle[int], 3)
	for i := 0; i < 6; i++ {
		if i%2 == 0 {
			tasks1[i/2] = workerBool.Spawn(fn1)
		} else {
			tasks2[i/2] = workerInt.Spawn(fn2)
		}
	}

	for _, task := range tasks1 {
		res, err := task.Await(context.Background())
		if err != nil {
			t.Error("TestTypedWorker failed. Error should be nil. Error:", err)
		}
		if res != true {
			t.Error("TestTypedWorker failed. Result must be true. Result:", res)
		}
	}
	for _, task := range tasks2 {
		res, err := task.Await(context.Background())
		if err != nil {
			t.Error("TestTypedWorker failed. Error should be nil. Error:", err)
		}
		if res != 1 {
			t.Error("TestTypedWorker failed. Result must be 1. Result:", res)
		}
	}

	finish := time.Now()

	if duration := int(finish.Sub(start).Seconds()); duration != 2 {
		t.Errorf("TestTypedWorker should complete in 2 seconds. Not in %d seconds.", duration)
	}
}
