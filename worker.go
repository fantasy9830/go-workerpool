package workerpool

import (
	"context"
)

type Worker struct {
	pool *WorkerPool
	done chan struct{}
}

func NewWorker(pool *WorkerPool) *Worker {
	return &Worker{
		pool: pool,
		done: make(chan struct{}),
	}
}

func (w *Worker) Start() {
	go func() {
		for {
			select {
			case <-w.pool.ctx.Done():
				return
			case task := <-w.pool.Tasks:
				w.execute(task)
			case <-w.done:
				// we have received a signal to stop
				return
			}
		}
	}()
}

// Stop signals the worker to stop listening for work requests.
func (w *Worker) Stop() {
	go func() {
		w.done <- struct{}{}
	}()
}

func (w *Worker) execute(task Task) {
	ctx := w.pool.ctx
	if w.pool.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(w.pool.ctx, w.pool.timeout)
		defer cancel()
	}

	done := make(chan struct{})

	go func() {
		task(ctx)
		close(done)
	}()

	select {
	case <-done:
		// The task finished successfully
		return
	case <-ctx.Done():
		// The context timed out, the task took too long
		return
	}
}
