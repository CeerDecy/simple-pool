package pool

import (
	"sync"
	"sync/atomic"
	"time"
)

type JobFunc func() any

type Options struct {
	PoolLimit       uint64
	WorkerStopCount int
	WorkerWaitTime  time.Duration
}

type WorkerConfig struct {
	WorkerStopCount int
	WorkerWaitTime  time.Duration
}

type Pool struct {
	limit        uint64
	jobs         chan JobFunc
	jobCount     int64
	jobDone      int64
	workers      map[string]*Worker
	wg           *sync.WaitGroup
	lock         *sync.Mutex
	workerConfig WorkerConfig
}

func New(opt Options) *Pool {
	if opt.WorkerWaitTime == 0 {
		opt.WorkerWaitTime = time.Second
	}
	pool := &Pool{
		limit:   opt.PoolLimit,
		jobs:    make(chan JobFunc, opt.PoolLimit),
		workers: make(map[string]*Worker),
		wg:      &sync.WaitGroup{},
		lock:    &sync.Mutex{},
		workerConfig: WorkerConfig{
			WorkerStopCount: opt.WorkerStopCount,
			WorkerWaitTime:  opt.WorkerWaitTime,
		},
	}
	return pool
}

func (p *Pool) Submit(job JobFunc) {
	p.jobs <- job
	atomic.AddInt64(&p.jobCount, 1)
	p.lock.Lock()
	defer p.lock.Unlock()
	if uint64(len(p.workers)) == p.limit {
		return
	}
	if len(p.jobs) < len(p.workers) {
		return
	}
	p.wg.Add(1)
	w := NewWorker(p)
	if p.workers[w.WorkerId] != nil {
		return
	}
	go w.Run()
	p.workers[w.WorkerId] = w
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
	close(p.jobs)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
