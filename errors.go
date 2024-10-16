package fastorm

import "errors"

var (
	ErrInvalidIsolationLevel = errors.New("invalid isolation level")
	ErrTxnInactive           = errors.New("transaction is inactive")
	ErrNotFound              = errors.New("not found")
)
