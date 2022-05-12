package commonerrors

import "errors"

var (
	ErrUserAborted   = errors.New("aborted by user")
	NoActionRequired = errors.New("no action required")
)
