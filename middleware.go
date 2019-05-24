package taskgroup

import (
	"context"
)

type task func(ctx context.Context, ready func()) error

type middleware func(task) task

func wait(ss []Signal) middleware {
	return func(t task) task {
		return func(ctx context.Context, ready func()) error {
			for _, s := range ss {
				select {
				case <-ctx.Done():
					return ctx.Err()
				case <-s.Ready():
				case <-s.done():
					select {
					case <-s.Ready():
					default:
						// assume deadlock occurred when done is closed but Ready is not closed.
						return ErrDeadlock
					}
				}
			}
			return t(ctx, ready)
		}
	}
}
