package pooli

import (
	"context"
	"sync"
)

type Pool struct {
	ctx        context.Context
	wg         sync.WaitGroup
	goroutines []*Goroutine
	pipe       chan Task
}

func Open(ctx context.Context, config Config) *Pool {
	p := &Pool{
		ctx: ctx,
	}

	setupConfig(config, p)

	return p
}

func (p *Pool) SendTask(task Task) {
	p.wg.Add(1)
	p.pipe <- task
}

func (p *Pool) TaskWait() {
	p.wg.Wait()
}

func (p *Pool) Start() {
	go func() {
		for _, goroutine := range p.goroutines {
			goroutine.Start()
		}
	}()
}

func (p *Pool) SetGoroutines(n int) {
	if len(p.goroutines) == n {
		return
	}

	n = len(p.goroutines) - n
	if n > 0 {
		for i := 0; i < n; i++ {
			if len(p.goroutines) > 0 {
				g := p.goroutines[0]
				p.RemoveGoroutine(g)
				go g.Kill()
			}
		}
	} else {
		for i := n; i < 0; i++ {
			g := NewGoroutine(p, p.pipe)
			p.AddGoroutine(g)
		}
	}
}

func (p *Pool) Len() int {
	return len(p.goroutines)
}

func (p *Pool) Goroutines() []*Goroutine {
	return p.goroutines
}

func (p *Pool) AddGoroutine(g *Goroutine) {
	g.Start()
	p.goroutines = append(p.goroutines, g)
}

func (p *Pool) RemoveGoroutine(g *Goroutine) {
	for i, gr := range p.goroutines {
		if gr != g {
			continue
		}

		gr.cnl()
		p.goroutines = append(p.goroutines[:i], p.goroutines[i+1:]...)
	}
}

func (p *Pool) Close() {
	for _, g := range p.goroutines {
		p.RemoveGoroutine(g)
		go g.Kill()
	}
}
