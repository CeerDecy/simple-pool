package pool

import (
	"sync"
	"sync/atomic"
	"time"
)

type JobFunc func() any

type Pool struct {
	limit uint64
	//current uint64
	jobs     chan JobFunc
	jobCount int64
	jobDone  int64
	workers  []*Worker
	wg       *sync.WaitGroup
	lock     *sync.Mutex
}

func New(limit uint64) *Pool {
	pool := &Pool{
		limit: limit,
		jobs:  make(chan JobFunc, limit),
		wg:    &sync.WaitGroup{},
	}
	for i := 0; i < int(limit); i++ {
		pool.wg.Add(1)
		w := Worker{
			i,
			pool,
			make(chan struct{}),
			pool.wg,
		}
		go w.Run()
		pool.workers = append(pool.workers, &w)
	}
	return pool
}

func (p *Pool) Submit(job JobFunc) {
	p.jobs <- job
	atomic.AddInt64(&p.jobCount, 1)
}

func (p *Pool) Close() {
	for {
		if atomic.LoadInt64(&p.jobCount) == atomic.LoadInt64(&p.jobDone) {
			break
		}
		time.Sleep(10 * time.Millisecond) // 避免 busy loop
	}
	for _, w := range p.workers {
		w.done <- struct{}{}
	}

	p.Wait()
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
