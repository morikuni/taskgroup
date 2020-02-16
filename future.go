package taskgroup

import (
	"context"
)

type Future struct {
	done <-chan struct{}
	err  error
}

func NewFuture() (*Future, func(error)) {
	c := make(chan struct{})
	f := &Future{
		done: c,
		err:  nil,
	}

	return f, func(err error) {
		f.err = err
		close(c)
	}
}

func (f *Future) Wait(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-f.done:
		return f.err
	}
}

func GoFuture(g *Group, f func(ctx context.Context) error) *Future {
	fut, tell := NewFuture()

	g.Go(func(ctx context.Context) error {
		err := f(ctx)
		tell(err)
		return err
	})

	return fut
}

func WaitSuccess(ctx context.Context, fs ...*Future) error {
	for _, ft := range fs {
		err := ft.Wait(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}
