package core

import (
	"errors"
)

var ErrEmptyRule = errors.New("empty rule")
var ErrAbsoluteMultiAddress = errors.New("absolute mode used multi address")
