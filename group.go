package taskgroup

import (
	"context"
	"sync"
)

type Group struct {
	tasks  []Task
	config *config
}

func New(opts ...Option) *Group {
	return &Group{
		config: evaluateOptions(opts),
	}
}

func FailFast(ctx context.Context, opts ...Option) (*Group, context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	defaultOpts := []Option{
		WithGoroutine(),
		WithRecover(),
		WithFoldFunc(func(acc, err error) error {
			if acc != nil {
				return acc
			}

			if err != nil {
				cancel()
			}

			return err
		}),
	}
	g := New(append(defaultOpts, opts...)...)

	return g, ctx, cancel
}

func (g *Group) Add(t Task) {
	g.tasks = append(g.tasks, t)
}

func (g *Group) AddFunc(f func(ctx context.Context) error) {
	g.tasks = append(g.tasks, TaskFunc(f))
}

func (g *Group) Process(ctx context.Context) error {
	var (
		result error
		mu     sync.Mutex
		count  = len(g.tasks)
		done   = make(chan struct{})
	)

	report := func() func(err error) {
		called := false
		return func(err error) {
			mu.Lock()
			defer mu.Unlock()

			// ignore except first call.
			if called {
				return
			}
			called = true

			result = g.config.fold(result, err)

			count--
			if count == 0 {
				close(done)
			}
		}
	}

	for _, t := range g.tasks {
		if g.config.interceptor == nil {
			g.config.runner.Run(ctx, report(), t)
		} else {
			g.config.interceptor(ctx, report(), t, g.config.runner)
		}
	}

	<-done

	return result
}

type Task interface {
	Process(ctx context.Context) error
}

type TaskFunc func(ctx context.Context) error

func (f TaskFunc) Process(ctx context.Context) error {
	return f(ctx)
}

type Runner interface {
	Run(ctx context.Context, report func(error), t Task)
}

var defaultRunner = runnerFunc(func(ctx context.Context, report func(error), t Task) {
	report(t.Process(ctx))
})
var DefaultRunner Runner = defaultRunner

type runnerFunc func(ctx context.Context, report func(error), t Task)

func (f runnerFunc) Run(ctx context.Context, report func(error), t Task) {
	f(ctx, report, t)
}
