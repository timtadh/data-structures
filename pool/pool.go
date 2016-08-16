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
	workin  sync.WaitGroup
	mu      sync.RWMutex
}

func New(n int) *Pool {
	pool := &Pool{
		workers: make([]*worker, 0, n),
	}
	for i := 0; i < n; i++ {
		w := &worker{
			in: make(chan func()),
			wg: &pool.wg,
			workin: &pool.workin,
		}
		go w.work()
		pool.workers = append(pool.workers, w)
	}
	return pool
}

func (p *Pool) WaitLock() {
	p.mu.Lock()
	p.workin.Wait()
}

func (p *Pool) Unlock() {
	p.mu.Unlock()
}

func (p *Pool) Stop() {
	p.mu.Lock()
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
	p.workin.Add(1)
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
	workin *sync.WaitGroup
}

func (w *worker) work() {
	w.wg.Add(1)
	for f := range w.in {
		f()
		w.workin.Done()
	}
	w.wg.Done()
}
