package emit_msg

import (
	"errors"
)

var (
	ErrCountMustBeGTZero = errors.New("cmd.Count must be greater than 0")
)

type Handler struct {
}

type Command struct {
	StartOrderIndex int64
	Count           int
}

type Result struct {
}

func Handle(cmd Command) (*Result, error) {
	if cmd.Count <= 0 {
		return nil, ErrCountMustBeGTZero
	}

	// ...

	return &Result{}, nil
}
