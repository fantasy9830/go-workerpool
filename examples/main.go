package main

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/fantasy9830/go-workerpool"
)

func main() {
	options := []workerpool.OptionFunc{
		workerpool.WithMaxWorker(5),
		workerpool.WithMaxTask(10),
		workerpool.WithTimeout(2 * time.Second),
	}

	pool := workerpool.New(options...)
	defer pool.Stop()

	var wg sync.WaitGroup
	for i := 1; i <= 10; i++ {
		num := i
		wg.Add(1)
		pool.AddTask(func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					return
				default:
					log.Println(num)
					time.Sleep(1 * time.Second)
				}
			}
		})
	}

	wg.Wait()
}
