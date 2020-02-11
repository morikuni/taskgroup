package taskgroup_test

import (
	"errors"
	"testing"

	"github.com/morikuni/taskgroup"
)

func TestMultiError(t *testing.T) {
	errA := errors.New("a")
	errB := errors.New("b")

	var me taskgroup.MultiError

	if want, got := error(nil), me.First(); want != got {
		t.Fatalf("want %v, got %v", want, got)
	}

	me.Append(errA)
	me.Append(errB)

	if want, got := "multi error: [2]{a, b}", me.Error(); want != got {
		t.Fatalf("want %v, got %v", want, got)
	}

	if want, got := errA, me.First(); want != got {
		t.Fatalf("want %v, got %v", want, got)
	}

	errs := me.Errors()
	if want, got := 2, len(errs); want != got {
		t.Fatalf("want %v, got %v", want, got)
	}

	if want, got := errA, errs[0]; want != got {
		t.Fatalf("want %v, got %v", want, got)
	}

	if want, got := errB, errs[1]; want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}
