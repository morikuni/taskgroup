package taskgroup

import (
	"errors"
)

var (
	ErrDeadlock = errors.New("deadlock detected")
)
