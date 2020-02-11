package taskgroup

import (
	"context"
	"sync"
)

// Group represents a group of tasks running concurrently.
type Group struct {
	wg   sync.WaitGroup
	ctx  context.Context
	fold func(acc, err error) error
	err  error
	mu   sync.Mutex
}

func New(ctx context.Context) *Group {
	return &Group{
		ctx:  ctx,
		fold: appendError,
	}
}

func WithCancel(ctx context.Context) (*Group, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	return New(ctx), cancel
}

func (g *Group) Go(f func(ctx context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		err := f(g.ctx)

		g.mu.Lock()
		defer g.mu.Unlock()

		g.err = g.fold(g.err, err)
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}
