package taskgroup

import (
	"context"
	"fmt"
)

type Interceptor func(ctx context.Context, report func(error), t Task, r Runner)

func (i Interceptor) apply(c *config) {
	if c.interceptor == nil {
		c.interceptor = i
	}
	c.interceptor = ChainInterceptor(c.interceptor, i)
}

func ChainInterceptor(is ...Interceptor) Interceptor {
	switch len(is) {
	case 0:
		panic("ChainInterceptor requires at least one parameter")
	case 1:
		return is[0]
	default: // > 2
		head, tail := is[0], ChainInterceptor(is[1:]...)
		return func(ctx context.Context, report func(error), t Task, r Runner) {
			head(ctx, report, t, applyInterceptor(tail, r))
		}
	}
}

func applyInterceptor(i Interceptor, r Runner) Runner {
	return runnerFunc(func(ctx context.Context, report func(error), t Task) {
		i(ctx, report, t, r)
	})
}

func WithGoroutine() Interceptor {
	return func(ctx context.Context, report func(error), t Task, r Runner) {
		go func() {
			r.Run(ctx, report, t)
		}()
	}
}

type PanicError struct {
	Raw interface{}
}

func (pe *PanicError) Error() string {
	return fmt.Sprintf("panic error: %v", pe.Raw)
}

func (pe *PanicError) Unwrap() error {
	if err, ok := pe.Raw.(error); ok {
		return err
	}
	return nil
}

func WithRecover() Interceptor {
	return func(ctx context.Context, report func(error), t Task, r Runner) {
		defer func() {
			r := recover()
			if r != nil {
				report(&PanicError{r})
			}
		}()
		r.Run(ctx, report, t)
	}
}
