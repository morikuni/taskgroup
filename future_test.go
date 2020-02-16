package taskgroup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"

	"github.com/morikuni/taskgroup"
)

func TestFuture(t *testing.T) {
	g := taskgroup.New(context.Background())

	errA := errors.New("a")

	var counter int64
	f1 := taskgroup.GoFuture(g, func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return nil
	})
	f2 := taskgroup.GoFuture(g, func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return errA
	})

	f3 := taskgroup.GoFuture(g, func(ctx context.Context) error {
		err := taskgroup.WaitSuccess(ctx, f1)
		if err != nil {
			return err
		}

		atomic.AddInt64(&counter, 1)
		return errors.New("b")
	})

	g.Go(func(ctx context.Context) error {
		err := taskgroup.WaitSuccess(ctx, f2, f3)
		if err != nil {
			return err
		}

		atomic.AddInt64(&counter, 1)
		return nil
	})

	err := g.Wait()
	if err == nil {
		t.Fatal("want error")
	}

	me := err.(*taskgroup.MultiError)
	if want, got := 3, len(me.Errors()); want != got {
		t.Fatalf("want %v, got %v", want, got)
	}

	_, ok := err.(*taskgroup.MultiError)
	if !ok {
		t.Fatal("want MultiError")
	}

	if want, got := int64(3), counter; want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}
