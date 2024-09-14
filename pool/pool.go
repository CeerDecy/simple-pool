package pool

import (
	"sync"
	"sync/atomic"
	"time"
)

const DefaultWaitTime = time.Millisecond * 500
const DefaultFreeCount = 3

type TaskFunc func() any

type Options struct {
	PoolLimit       uint64
	WorkerFreeCount int
	WorkerWaitTime  time.Duration
}

type Pool struct {
	tasks   chan TaskFunc
	count   int64
	done    int64
	workers map[string]*Worker
	wg      *sync.WaitGroup
	lock    *sync.Mutex
	options Options
}

func NewWithOptions(opt Options) *Pool {
	return &Pool{
		tasks:   make(chan TaskFunc, opt.PoolLimit),
		workers: make(map[string]*Worker),
		wg:      &sync.WaitGroup{},
		lock:    &sync.Mutex{},
		options: opt,
	}
}

func New(limit uint64) *Pool {
	return NewWithOptions(Options{
		PoolLimit:       limit,
		WorkerWaitTime:  DefaultWaitTime,
		WorkerFreeCount: DefaultFreeCount,
	})
}

func (p *Pool) Submit(job TaskFunc) {
	p.tasks <- job
	atomic.AddInt64(&p.count, 1)
	p.lock.Lock()
	defer p.lock.Unlock()
	if uint64(len(p.workers)) == p.options.PoolLimit {
		return
	}
	if len(p.tasks) < len(p.workers)/4 {
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
		if atomic.LoadInt64(&p.count) == atomic.LoadInt64(&p.done) {
			break
		}
		time.Sleep(10 * time.Millisecond) // 避免 busy loop
	}
	p.lock.Lock()
	var workers []*Worker
	for _, w := range p.workers {
		workers = append(workers, w)
	}
	p.lock.Unlock()
	for _, w := range workers {
		w.done <- struct{}{}
		close(w.done)
	}

	p.Wait()
	close(p.tasks)
}

func (p *Pool) Wait() {
	p.wg.Wait()
}
