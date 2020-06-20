package taskgroup_test

import (
	"context"
	"errors"
	"reflect"
	"sync/atomic"
	"testing"
	"time"

	"github.com/morikuni/taskgroup"
)

func equal(t *testing.T, want, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %T(%#v) but got %T(%#v)", want, want, got, got)
	}
}
func noError(t *testing.T, err error) {
	t.Helper()

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGroup(t *testing.T) {
	g := taskgroup.New()

	count := 0
	g.AddFunc(func(ctx context.Context) error {
		count++
		equal(t, 1, count)
		return nil
	})
	g.AddFunc(func(ctx context.Context) error {
		count++
		equal(t, 2, count)
		return errors.New("hello world")
	})
	g.AddFunc(func(ctx context.Context) error {
		count++
		equal(t, 3, count)
		return nil
	})

	equal(t, 0, count)

	err := g.Process(context.Background())
	equal(t, errors.New("hello world"), err)
	equal(t, 3, count)
}

func TestFailFast(t *testing.T) {
	g, ctx, cancel := taskgroup.FailFast(context.Background())
	defer cancel()

	count := int64(0)
	g.AddFunc(func(ctx context.Context) error {
		equal(t, int64(1), atomic.AddInt64(&count, 1))
		return nil
	})
	g.AddFunc(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(time.Second):
		}
		t.Fatal("never come")
		return nil
	})
	g.AddFunc(func(ctx context.Context) error {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-time.After(100 * time.Millisecond):
		}
		equal(t, int64(2), atomic.AddInt64(&count, 1))
		panic("hello world")
	})

	equal(t, int64(0), count)

	err := g.Process(ctx)
	equal(t, &taskgroup.PanicError{Raw: "hello world"}, err)
	equal(t, int64(2), count)
}
