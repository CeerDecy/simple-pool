package pool

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/CeerDecy/simple-pool/tools"
)

type Worker struct {
	WorkerId  string
	pool      *Pool
	done      chan struct{}
	wg        *sync.WaitGroup
	freeCount int
	stopCount int
	waitTime  time.Duration
}

func NewWorker(pool *Pool) *Worker {
	return &Worker{
		WorkerId:  tools.RandString(8),
		pool:      pool,
		done:      make(chan struct{}, 1),
		wg:        pool.wg,
		freeCount: 0,
		waitTime:  pool.options.WorkerWaitTime,
		stopCount: pool.options.WorkerFreeCount,
	}
}

func (w *Worker) Run() {
	defer w.wg.Done()
	defer w.stop()
	for {
		select {
		case job := <-w.pool.tasks:
			if job == nil {
				return
			}
			job()
			atomic.AddInt64(&w.pool.done, 1)
		case <-w.done:
			return
		default:
			if w.freeCount >= w.stopCount {
				return
			}
			w.freeCount++
			time.Sleep(w.waitTime)
		}
	}
}

func (w *Worker) stop() {
	w.pool.lock.Lock()
	defer w.pool.lock.Unlock()
	delete(w.pool.workers, w.WorkerId)
}
