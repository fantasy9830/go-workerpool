package workerpool_test

import (
	"context"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/fantasy9830/go-workerpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type Suit struct {
	suite.Suite
}

func TestWorkerpool(t *testing.T) {
	suite.Run(t, new(Suit))
}

func (s *Suit) TestAddTask() {
	parentCtx := context.Background()

	options := []workerpool.OptionFunc{
		workerpool.WithContext(parentCtx),
		workerpool.WithMaxWorker(5),
		workerpool.WithMaxTask(10),
	}

	pool := workerpool.New(options...)
	defer pool.Stop()

	var count atomic.Int32
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		pool.AddTask(func(ctx context.Context) {
			defer wg.Done()
			count.Add(1)
		})
	}

	wg.Wait()

	assert.Equal(s.T(), int32(10), count.Load())
}

func (s *Suit) TestWithTimeout() {
	options := []workerpool.OptionFunc{
		workerpool.WithTimeout(100 * time.Millisecond),
	}

	pool := workerpool.New(options...)
	defer pool.Stop()

	var count atomic.Int32
	var wg sync.WaitGroup

	for i := 1; i <= 10; i++ {
		wg.Add(1)
		pool.AddTask(func(ctx context.Context) {
			defer wg.Done()
			for {
				select {
				case <-ctx.Done():
					count.Add(1)
					return
				default:
					time.Sleep(200 * time.Millisecond)
				}
			}
		})
	}

	wg.Wait()

	assert.Equal(s.T(), int32(10), count.Load())
}

func (s *Suit) TestAdjustWorkers() {
	options := []workerpool.OptionFunc{
		workerpool.WithMaxWorker(1),
	}

	pool := workerpool.New(options...)
	defer pool.Stop()

	pool.AdjustWorkers(2)
	assert.Equal(s.T(), 2, pool.CurrentWorkers())

	pool.AdjustWorkers(5)
	assert.Equal(s.T(), 5, pool.CurrentWorkers())

	pool.AdjustWorkers(3)
	assert.Equal(s.T(), 3, pool.CurrentWorkers())

	pool.AdjustWorkers(0)
	assert.Equal(s.T(), 0, pool.CurrentWorkers())

	pool.AdjustWorkers(-1)
	assert.Equal(s.T(), 0, pool.CurrentWorkers())
}

func (s *Suit) TestCurrentTasks() {
	options := []workerpool.OptionFunc{
		workerpool.WithMaxWorker(0),
	}

	pool := workerpool.New(options...)
	defer pool.Stop()

	for i := 0; i < 3; i++ {
		pool.AddTask(func(ctx context.Context) {
			<-ctx.Done()
		})
	}

	assert.Equal(s.T(), 3, pool.CurrentTasks())
}
