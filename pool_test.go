package pooli_test

import (
	"context"
	"sync/atomic"
	"testing"

	"github.com/0x9n0p/pooli"
	"github.com/go-playground/assert/v2"
)

func TestTaskWait(t *testing.T) {

	var doneCount int32

	p := pooli.Open(context.Background(), pooli.Config{
		Goroutines: 1,
		Pipe:       make(chan pooli.Task),
	})

	p.Start()

	for i := 0; i < 1000; i++ {
		go func(n int) {
			p.SendTask(pooli.NewTask(func(ctx context.Context) error {
				atomic.AddInt32(&doneCount, 1)

				return nil
			}))
		}(i)

	}
	p.TaskWait()

	assert.Equal(t, int32(1000), doneCount)

}
