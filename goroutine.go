package pooli

import (
	"context"
	"sync"
	"sync/atomic"
)

type Goroutine struct {
	status Status
	Pipe   chan Task
	p      *Pool
	ctx    context.Context
	cnl    context.CancelFunc

	m *sync.RWMutex
}

func NewGoroutine(p *Pool) *Goroutine {
	ctx, cnl := context.WithCancel(p.ctx)
	return &Goroutine{
		status: Idle,
		Pipe:   p.pipe,
		p:      p,
		ctx:    ctx,
		cnl:    cnl,

		m: new(sync.RWMutex),
	}
}

func (g *Goroutine) Start() {
	wg := new(sync.WaitGroup)
	wg.Add(1)

	go func() {
		wg.Done()
		for {
			select {
			case <-g.ctx.Done():
				return
			case t := <-g.Pipe:
				g.SetStatus(Progress)
				ExecuteTask(g.ctx, t)
				if atomic.AddInt32(&g.p.sentTaskNum, -1) == 0 && len(g.p.pipe) == 0 { // all task done
					g.p.done <- struct{}{}
				}
				g.SetStatus(Idle)

			}
		}
	}()

	wg.Wait()
}

func (g *Goroutine) SetStatus(status Status) {
	g.m.Lock()
	defer g.m.Unlock()

	g.status = status
}

func (g Goroutine) Status() Status {
	return g.status
}

func (g *Goroutine) Kill() {
	g.cnl()
}
