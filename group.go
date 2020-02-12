package taskgroup

import (
	"context"
	"sync"
)

// Group represents a group of tasks running concurrently.
type Group struct {
	conf *config

	wg  sync.WaitGroup
	ctx context.Context
	err error
	mu  sync.Mutex
}

func New(ctx context.Context, opts ...Option) *Group {
	return &Group{
		conf: newConfig(opts),
		ctx:  ctx,
	}
}

func FailFast(ctx context.Context, opts ...Option) (*Group, context.CancelFunc) {
	ctx, cancel := context.WithCancel(ctx)
	foldFunc := FoldFunc(func(acc, err error) error {
		if err != nil {
			cancel()
		}

		// TakeFirst because other goroutines would be context.Canceled.
		return TakeFirst(acc, err)
	})
	return New(ctx, append(opts, foldFunc)...), cancel
}

func (g *Group) Go(f func(ctx context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()

		err := f(g.ctx)

		g.mu.Lock()
		defer g.mu.Unlock()

		g.err = g.conf.foldFunc(g.err, err)
	}()
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}
