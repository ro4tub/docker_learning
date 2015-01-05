package common

import (
	"errors"
)

const (
	// common
	// ok
	OK = 0
	// param error
	ParamErr = 1
	// internal error
	InternalErr = 2
)

var (
	ErrParam = errors.New("parameter error")
)
