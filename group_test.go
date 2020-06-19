package taskgroup_test

import (
	"context"
	"errors"
	"reflect"
	"testing"

	"github.com/morikuni/taskgroup"
)

func equal(t *testing.T, want, got interface{}) {
	t.Helper()

	if !reflect.DeepEqual(want, got) {
		t.Errorf("want %#v but got %#v", want, got)
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

	err := g.Run(context.Background())
	equal(t, errors.New("hello world"), err)
	equal(t, 3, count)
}
