package workerpool_test

import (
	"context"
	"testing"

	"github.com/fantasy9830/go-workerpool"
)

func BenchmarkWorkerpool(b *testing.B) {
	b.ReportAllocs()
	pool := workerpool.New()
	defer pool.Stop()

	for n := 0; n < b.N; n++ {
		pool.AddTask(func(ctx context.Context) {})
	}
}
