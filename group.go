package taskgroup

import (
	"context"
	"sync"
)

type Group struct {
	opts groupOption

	wg   sync.WaitGroup
	once sync.Once
	err  error
}

func newGroup(opts ...GroupOption) *Group {
	var gopt groupOption
	gopt.apply(opts)
	return &Group{opts: gopt}
}

func WithContext(ctx context.Context, opts ...GroupOption) (*Group, context.Context, func()) {
	ctx, cancel := context.WithCancel(ctx)
	return newGroup(onFirstError(cancel)), ctx, cancel
}

func (g *Group) Go(ctx context.Context, f func(ctx context.Context) error, opts ...CallOption) Signal {
	return g.GoReady(ctx, func(ctx context.Context, ready func()) error {
		err := f(ctx)
		if err == nil {
			ready()
		}
		return err
	}, opts...)
}

func (g *Group) GoReady(ctx context.Context, t func(ctx context.Context, ready func()) error, opts ...CallOption) Signal {
	var co callOption
	co.apply(opts)
	for _, m := range co.middlewares {
		t = m(t)
	}

	var (
		ready = make(chan struct{})
		done  = make(chan struct{})
	)
	g.wg.Add(1)
	go func() {
		defer close(done)
		defer g.wg.Done()

		err := t(ctx, func() { close(ready) })
		g.done(err)
	}()

	return &signal{ready, done}
}

func (g *Group) done(err error) {
	if err == nil {
		return
	}

	g.once.Do(func() {
		g.err = err
		if f := g.opts.onFirstError; f != nil {
			f(err)
		}
	})
}

func (g *Group) Wait() error {
	g.wg.Wait()
	return g.err
}
