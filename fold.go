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

func (me *MultiError) Errors() []error {
	return me.errors
}

func (me *MultiError) First() error {
	if len(me.errors) == 0 {
		return nil
	}

	return me.errors[0]
}

func (me *MultiError) Append(err error) {
	const defaultSize = 8

	if me.errors == nil {
		me.errors = make([]error, 0, defaultSize)
	}

	me.errors = append(me.errors, err)
}

func appendError(acc, err error) error {
	if err == nil {
		return acc
	}

	if acc == nil {
		var me MultiError
		me.Append(err)
		return &me
	}

	me := acc.(*MultiError)
	me.Append(err)
	return me
}
