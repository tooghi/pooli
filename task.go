package pooli

import "context"

type Task struct {
	execute func(ctx context.Context) error
	success func()
	fail    func(error)
	final   func()
}

func NewTask(execute func(ctx context.Context) error) Task {
	return Task{execute: execute}
}

func (t Task) Success(s func()) Task {
	t.success = s
	return t
}

func (t Task) Fail(f func(err error)) Task {
	t.fail = f
	return t
}

func (t Task) Final(f func()) Task {
	t.final = f
	return t
}

// func ExecuteTask(ctx context.Context, task Task) {
// 	err := task.execute(ctx)
// 	if err == nil {
// 		if task.success != nil {
// 			task.success()
// 		}
// 	} else {
// 		if task.fail != nil {
// 			task.fail(err)
// 		}
// 	}
// 	if task.final != nil {
// 		task.final()
// 	}
// }

func ExecuteTask(g *Goroutine, task Task) {
	defer g.Pool.wg.Done()
	err := task.execute(g.ctx)
	if err == nil {
		if task.success != nil {
			task.success()
		}
	} else {
		if task.fail != nil {
			task.fail(err)
		}
	}
	if task.final != nil {
		task.final()
	}
}
