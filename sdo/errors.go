package sdo

import "errors"

var (
	ErrIncompletData = errors.New("incomplete data")
	ErrInvalidData   = errors.New("invalid data")
)
