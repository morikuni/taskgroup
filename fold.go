package taskgroup

import (
	"strconv"
	"strings"
)

type MultiError struct {
	errors []error
}

func (me *MultiError) Error() string {
	var b strings.Builder
	b.WriteString("multi error: [")
	b.WriteString(strconv.Itoa(len(me.errors)))
	b.WriteString("]{")

	first := true
	for _, err := range me.errors {
		if !first {
			b.WriteString(", ")
		}
		b.WriteString(err.Error())
		first = false
	}

	b.WriteString("}")

	return b.String()
}

func appendError(acc, err error) error {
	const defaultSize = 8

	if err == nil {
		return acc
	}

	if acc == nil {
		return &MultiError{append(make([]error, 0, defaultSize), err)}
	}

	me := acc.(*MultiError)
	me.errors = append(me.errors, err)
	return me
}
