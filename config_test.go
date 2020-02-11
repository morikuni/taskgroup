package taskgroup_test

import (
	"context"
	"errors"
	"testing"

	"github.com/morikuni/taskgroup"
)

func TestFoldFunc(t *testing.T) {
	errA := errors.New("a")
	errB := errors.New("b")

	g := taskgroup.New(context.Background(),
		taskgroup.FoldFunc(func(acc, err error) error {
			return errA
		}),
	)

	g.Go(func(ctx context.Context) error {
		return errB
	})

	err := g.Wait()
	if want, got := errA, err; want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}
