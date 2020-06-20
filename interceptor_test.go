package taskgroup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/morikuni/taskgroup"
)

func TestChainInterceptor(t *testing.T) {
	count := 0
	i1 := taskgroup.Interceptor(func(ctx context.Context, report func(error), task taskgroup.Task, r taskgroup.Runner) {
		count++
		equal(t, 1, count)
		r.Run(ctx, report, task)
	})
	i2 := taskgroup.Interceptor(func(ctx context.Context, report func(error), task taskgroup.Task, r taskgroup.Runner) {
		r.Run(ctx, report, task)
		count++
		equal(t, 4, count)
	})
	i3 := taskgroup.Interceptor(func(ctx context.Context, report func(error), task taskgroup.Task, r taskgroup.Runner) {
		count++
		equal(t, 2, count)
		r.Run(ctx, report, task)
	})

	i := taskgroup.ChainInterceptor(i1, i2, i3)

	var err error
	i(context.Background(), func(e error) { err = e }, taskgroup.TaskFunc(func(ctx context.Context) error {
		count++
		equal(t, 3, count)
		return errors.New("hello world")
	}), taskgroup.DefaultRunner)

	equal(t, errors.New("hello world"), err)
	equal(t, 4, count)
}

func TestWithGoroutine(t *testing.T) {
	g := taskgroup.New(
		taskgroup.WithGoroutine(),
	)

	count := int64(0)
	g.AddFunc(func(ctx context.Context) error {
		time.Sleep(100 * time.Millisecond)
		equal(t, int64(2), atomic.AddInt64(&count, 1))
		return nil
	})
	g.AddFunc(func(ctx context.Context) error {
		time.Sleep(200 * time.Millisecond)
		equal(t, int64(3), atomic.AddInt64(&count, 1))
		return errors.New("hello world")
	})
	g.AddFunc(func(ctx context.Context) error {
		equal(t, int64(1), atomic.AddInt64(&count, 1))
		return nil
	})

	equal(t, int64(0), count)

	err := g.Process(context.Background())
	equal(t, errors.New("hello world"), err)
	equal(t, int64(3), count)
}

func TestWithRecover(t *testing.T) {
	g := taskgroup.New(
		taskgroup.WithRecover(),
	)

	count := 0
	g.AddFunc(func(ctx context.Context) error {
		count++
		equal(t, 1, count)
		return nil
	})
	g.AddFunc(func(ctx context.Context) error {
		count++
		equal(t, 2, count)
		panic("hello")
	})
	g.AddFunc(func(ctx context.Context) error {
		count++
		equal(t, 3, count)
		return nil
	})

	equal(t, 0, count)

	err := g.Process(context.Background())
	equal(t, &taskgroup.PanicError{Raw: "hello"}, err)
	equal(t, 3, count)
}
