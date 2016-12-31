package pool

import (
	"sync"
	"math/rand"
)

import (
	"github.com/timtadh/data-structures/errors"
)

type Pool struct {
	workers []*worker
	wg      sync.WaitGroup
	workin  int
	workCond *sync.Cond
	mu      sync.RWMutex
}

func New(n int) *Pool {
	pool := &Pool{
		workers: make([]*worker, 0, n),
		workCond: sync.NewCond(&sync.Mutex{}),
	}
	for i := 0; i < n; i++ {
		w := &worker{
			in: make(chan func(), 100),
			wg: &pool.wg,
			workin: &pool.workin,
			workCond: pool.workCond,
		}
		go w.work()
		pool.workers = append(pool.workers, w)
	}
	return pool
}

func (p *Pool) WaitCount() int {
	p.workCond.L.Lock()
	if p.workin <= 0 {
		n := p.workin
		p.workCond.L.Unlock()
		return n
	}
	p.workCond.Wait()
	n := p.workin
	p.workCond.L.Unlock()
	return n
}

func (p *Pool) WaitLock() {
	p.mu.Lock()
	p.workCond.L.Lock()
	for p.workin > 0 {
		p.workCond.Wait()
	}
}

func (p *Pool) Unlock() {
	p.workCond.L.Unlock()
	p.mu.Unlock()
}

func (p *Pool) Stop() {
	p.mu.Lock()
	if len(p.workers) == 0 {
		p.mu.Unlock()
		return
	}
	workers := p.workers
	p.workers = nil
	p.mu.Unlock()
	for _, wrkr := range workers {
		close(wrkr.in)
	}
	p.wg.Wait()
}

// This will only return an error if the pool is stopped
func (p *Pool) Do(f func()) error {
	p.mu.RLock()
	defer p.mu.RUnlock()
	if len(p.workers) <= 0 {
		return errors.Errorf("The pool was stopped")
	}
	p.workCond.L.Lock()
	p.workin += 1
	p.workCond.L.Unlock()
	offset := rand.Intn(len(p.workers))
	for i := 0; i < len(p.workers); i++ {
		j := (offset + i) % len(p.workers)
		wrkr := p.workers[j].in
		select {
		case wrkr<-f:
			return nil
		default:
		}
	}
	p.workers[offset].in<-f
	return nil
}

type worker struct {
	in chan func()
	wg *sync.WaitGroup
	workin *int
	workCond *sync.Cond
}

func (w *worker) work() {
	w.wg.Add(1)
	for f := range w.in {
		f()
		w.workCond.L.Lock()
		*w.workin -= 1
		w.workCond.L.Unlock()
		w.workCond.Broadcast()
	}
	w.wg.Done()
}
