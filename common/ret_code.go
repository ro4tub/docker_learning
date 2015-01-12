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
	// 找不到玩家
	NotFoundErr = 3
)

var (
	ErrParam = errors.New("parameter error")
	ErrNotFoundPlayer = errors.New("not found player")
)
