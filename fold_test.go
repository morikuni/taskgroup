package taskgroup_test

import (
	"errors"
	"testing"
)

func TestMultiError_Error(t *testing.T) {
	me := &MultiError{errors: []error{
		errors.New("a"),
		errors.New("b"),
	}}

	want := "multi error: [2]{a, b}"
	got := me.Error()

	if want != got {
		t.Fatalf("want %v, got %v", want, got)
	}
}
