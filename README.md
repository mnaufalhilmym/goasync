# goasync

Thread-safe asynchronous library for Go

## Usage

- Using asynchronous task

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/mnaufalhilmym/goasync"
)

task := goasync.Spawn(func(ctx context.Context) (bool, error) {
    select {
    case <-ctx.Done():
        return false, ctx.Err()
    case <-time.After(2 * time.Second):
        return true, nil
    }
})

result, err := task.Await(context.Background()) // Wait for the task to complete (2 secs)
if err != nil { // false
    fmt.Println(err) // Never executed
}

fmt.Println(result) // Print: true
```

- Aborting asynchronous task

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/mnaufalhilmym/goasync"
)

task := goasync.Spawn(func(ctx context.Context) ([]any, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-time.After(2 * time.Second):
        return []any{1, "2-two", 3.3}, nil
    }
})

task.Abort() // Aborting async task

result, err := task.Await(context.Background()) // Immediately return the result
if err != nil { // true
    fmt.Println(err) // Print: context canceled
}

fmt.Println(result) // Print: []
```

- Waiting for many asynchronous tasks

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/mnaufalhilmym/goasync"
)

fn := func(ctx context.Context) ([]any, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-time.After(2 * time.Second):
        return []any{1, "2-two", 3.3}, nil
    }
}

task1 := goasync.Spawn(fn)
task2 := goasync.Spawn(fn)
task3 := goasync.Spawn(fn)

results, err := goasync.TryJoin(context.Background(), task1, task2, task3) // Wait for all tasks to complete (2 secs) or return immediately if an error occurs
if err != nil { // false
    fmt.Println(err) // Never executed
}

for _, result := range results {
    fmt.Println(result) // Print: [1 2-two 3.3]
}
```

- Await asynchronous tasks with timeout

```go
import (
    "context"
    "fmt"
    "time"

    "github.com/mnaufalhilmym/goasync"
)

fn := func(ctx context.Context) ([]any, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-time.After(2 * time.Second):
        return []any{1, "2-two", 3.3}, nil
    }
}

task1 := goasync.Spawn(fn)
task2 := goasync.Spawn(fn)
task3 := goasync.Spawn(fn)

ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
defer cancel()
results, err := goasync.TryJoin(ctx, task1, task2, task3) // Wait for all tasks to complete (2 secs) or return immediately if an error occurs
if err != nil { // true
    fmt.Println(err) // Print: context deadline exceeded
}

fmt.Println(results) // Print: []
```

- Limit the number of goroutines

```go
import (
	"context"
	"time"

	"github.com/mnaufalhilmym/goasync"
)

worker := goasync.
    NewWorker().
    SetMaxWorkers(3). // Limit to a maximum of three goroutines
    Build()

fn := func(context.Context) (any, error) {
    time.Sleep(1 * time.Second)
    return true, nil
}

tasks := make([]goasync.JoinHandle[any], 6)
for i := 0; i < len(tasks); i++ {
    tasks[i] = worker.Spawn(fn) // Will block execution when all workers are used. If there are workers available, it will unblock
}

for _, task := range tasks {
    res, err := task.Await(context.Background()) // Wait for the task to complete
    if err != nil {
        fmt.Println(err) // Never executed
    }
    fmt.Println(result) // Print: true
}
```

> [!TIP]
> Further usage methods can be seen in the *_test.go files.
