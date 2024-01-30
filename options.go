package workerpool

import (
	"context"
	"time"
)

type OptionFunc func(*WorkerPool)

func WithContext(parent context.Context) OptionFunc {
	return func(p *WorkerPool) {
		ctx, cancel := context.WithCancel(parent)

		p.ctx = ctx
		p.cancel = cancel
	}
}

func WithMaxWorker(maxWorkers int) OptionFunc {
	return func(p *WorkerPool) {
		if maxWorkers >= 0 {
			p.maxWorkers = maxWorkers
		}
	}
}

func WithMaxTask(maxTasks int) OptionFunc {
	return func(p *WorkerPool) {
		if maxTasks > 0 {
			p.maxTasks = maxTasks
		}
	}
}

func WithTimeout(timeout time.Duration) OptionFunc {
	return func(p *WorkerPool) {
		p.timeout = timeout
	}
}
