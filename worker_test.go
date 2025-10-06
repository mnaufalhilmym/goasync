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

	worker := goasync.NewWorker(
		goasync.WorkerConfig{
			MaxWorkers: 3,
		},
	)

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
