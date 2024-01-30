package workerpool

import (
	"context"
	"runtime"
	"sync"
	"time"
)

type WorkerPool struct {
	mux        sync.Mutex
	ctx        context.Context
	cancel     context.CancelFunc
	maxWorkers int
	maxTasks   int
	timeout    time.Duration
	workers    []*Worker
	Tasks      chan Task
}

func New(optFuncs ...OptionFunc) *WorkerPool {
	ctx, cancel := context.WithCancel(context.Background())

	wp := &WorkerPool{
		ctx:        ctx,
		cancel:     cancel,
		maxWorkers: runtime.NumCPU(),
		maxTasks:   4096,
	}

	for _, f := range optFuncs {
		f(wp)
	}

	wp.Tasks = make(chan Task, wp.maxTasks)

	wp.start()

	return wp
}

func (wp *WorkerPool) start() {
	wp.workers = make([]*Worker, 0, wp.maxWorkers)
	for i := 0; i < wp.maxWorkers; i++ {
		worker := NewWorker(wp)
		wp.workers = append(wp.workers, worker)
		worker.Start()
	}
}

func (wp *WorkerPool) Stop() {
	wp.cancel()
}

func (wp *WorkerPool) AddTask(task Task) {
	wp.Tasks <- task
}

func (wp *WorkerPool) AdjustWorkers(maxWorkers int) {
	if maxWorkers < 0 {
		return
	}

	wp.mux.Lock()
	defer wp.mux.Unlock()

	currentWorkers := len(wp.workers)

	// Add workers
	if maxWorkers > currentWorkers {
		addWorkers := maxWorkers - currentWorkers
		for i := 0; i < addWorkers; i++ {
			worker := NewWorker(wp)
			wp.workers = append(wp.workers, worker)
			worker.Start()
		}

		// Remove workers
	} else if maxWorkers < currentWorkers {
		for i := currentWorkers - 1; i >= maxWorkers; i-- {
			wp.workers[i].Stop()
			wp.workers[i] = nil
		}

		removeWorkers := currentWorkers - maxWorkers
		wp.workers = wp.workers[:len(wp.workers)-removeWorkers]
	}

	wp.maxWorkers = maxWorkers
}

func (wp *WorkerPool) CurrentWorkers() int {
	return len(wp.workers)
}

func (wp *WorkerPool) CurrentTasks() int {
	return len(wp.Tasks)
}
