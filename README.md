# go-workerpool

## Installation
```shell
go get github.com/fantasy9830/go-workerpool
```


## Usage
### func WithContext
```go
func WithContext(parent context.Context)
```

### func WithMaxWorker
```go
func WithMaxWorker(maxWorkers int)
```

### func WithMaxTask
```go
func WithMaxTask(maxTasks int)
```

### func WithTimeout
```go
func WithTimeout(timeout time.Duration)
```

## examples
```go
options := []workerpool.OptionFunc{
    workerpool.WithMaxWorker(5),
    workerpool.WithMaxTask(10),
}

pool := workerpool.New(options...)
defer pool.Stop()

var wg sync.WaitGroup
for i := 1; i <= 10; i++ {
    num := i
    wg.Add(1)
    pool.AddTask(func(ctx context.Context) {
        defer wg.Done()
        log.Println(num)
        time.Sleep(5 * time.Second)
    })
}

wg.Wait()
```