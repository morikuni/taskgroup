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

func (g *Group) Add(t Task) {
	g.tasks = append(g.tasks, t)
}

func (g *Group) AddFunc(f func(ctx context.Context) error) {
	g.tasks = append(g.tasks, TaskFunc(f))
}

func (g *Group) Run(ctx context.Context) error {
	var (
		result error
		mu     sync.Mutex
		count  int
		done   chan struct{}
	)

	report := func(err error) {
		mu.Lock()
		defer mu.Unlock()
		result = g.config.fold(result, err)

		count--
		if done != nil && count == 0 {
			close(done)
		}
	}

	for _, t := range g.tasks {
		count++
		if g.config.interceptor == nil {
			g.config.runner.Run(ctx, report, t)
		} else {
			g.config.interceptor(ctx, report, t, g.config.runner)
		}
	}

	mu.Lock()

	// check if there is tasks still running
	// in case all process has completed
	// before done channel is created.
	if count == 0 {
		return result
	}

	done = make(chan struct{})
	mu.Unlock()

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
