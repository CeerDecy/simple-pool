package pool

import (
	"fmt"
	"sync"
	"sync/atomic"
)

type Worker struct {
	WID  int
	pool *Pool
	done chan struct{}
	wg   *sync.WaitGroup
}

func (w *Worker) Run() {
	defer w.wg.Done()
	for {
		select {
		case job := <-w.pool.jobs:
			fmt.Printf("worker %d received job\n", w.WID)
			job()
			atomic.AddInt64(&w.pool.jobDone, 1)
		case <-w.done:
			return
		}
	}
}
