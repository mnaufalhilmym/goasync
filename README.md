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

task := goasync.Spawn(func(ctx context.Context) ([]any, error) {
    select {
    case <-ctx.Done():
        return nil, ctx.Err()
    case <-time.After(2 * time.Second):
        return []any{1, "2-dua", 3.3}, nil
    }
})

result, err := task.Await(context.Background()) // Wait for the task to complete (2 secs)
if err != nil { // false
    fmt.Println(err) // never executed
}

fmt.Println(result) // print: [1 2-dua 3.3]
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
        return []any{1, "2-dua", 3.3}, nil
    }
})

task.Abort() // Aborting async task

result, err := task.Await(context.Background()) // Immediately return the result
if err != nil { // true
    fmt.Println(err) // print: context canceled
}

fmt.Println(result) // print: []
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
        return []any{1, "2-dua", 3.3}, nil
    }
}

task1 := goasync.Spawn(fn)
task2 := goasync.Spawn(fn)
task3 := goasync.Spawn(fn)

results, err := goasync.TryJoin(context.Background(), task1, task2, task3) // Wait for all tasks to complete (2 secs) or return immediately if an error occurs
if err != nil { // false
    fmt.Println(err) // never executed
}

for _, result := range results {
    fmt.Println(result) // print: [1 2-dua 3.3]
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
        return []any{1, "2-dua", 3.3}, nil
    }
}

task1 := goasync.Spawn(fn)
task2 := goasync.Spawn(fn)
task3 := goasync.Spawn(fn)

ctx, cancel := context.WithTimeout(context.Background(), 1 * time.Second)
defer cancel()
results, err := goasync.TryJoin(ctx, task1, task2, task3) // Wait for all tasks to complete (2 secs) or return immediately if an error occurs
if err != nil { // true
    fmt.Println(err) // print: context deadline exceeded
}

fmt.Println(results) // print: []
```

> [!TIP]
> Further usage methods can be seen in the goasync_test.go file.
