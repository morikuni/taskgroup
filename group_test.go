package taskgroup_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/morikuni/taskgroup"
)

func TestGroup(t *testing.T) {
	g := taskgroup.New(context.Background())

	var counter int64
	g.Go(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return errors.New("a")
	})
	g.Go(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return nil
	})
	g.Go(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return errors.New("b")
	})
	g.Go(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return nil
	})

	err := g.Wait()
	if err == nil {
		t.Fatal("want error")
	}

	_, ok := err.(*taskgroup.MultiError)
	if !ok {
		t.Fatal("want MultiError")
	}

	if want, got := int64(4), counter; want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}

func TestWithCancel(t *testing.T) {
	g, cancel := taskgroup.WithCancel(context.Background())
	defer cancel()

	var counter int64
	g.Go(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return nil
	})
	g.Go(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return nil
	})
	g.Go(func(ctx context.Context) error {
		atomic.AddInt64(&counter, 1)
		return nil
	})
	g.Go(func(ctx context.Context) error {
		select {
		case <-time.After(time.Second):
		case <-ctx.Done():
			return nil
		}
		atomic.AddInt64(&counter, 1)
		return nil
	})

	cancel()

	err := g.Wait()
	if err != nil {
		t.Fatalf("error: %v", err)
	}

	if want, got := int64(3), counter; want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}
