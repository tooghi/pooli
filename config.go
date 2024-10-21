package pooli

type Config struct {
	Goroutines int
	Pipe       chan Task
}

func setupConfig(config Config, p *Pool) {
	p.pipe = config.Pipe

	var goroutines []*Goroutine
	for i := 0; i < config.Goroutines; i++ {
		g := NewGoroutine(p, p.pipe)

		goroutines = append(goroutines, g)
	}

	p.goroutines = goroutines
}
