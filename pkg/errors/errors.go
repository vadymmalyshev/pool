package errors

import "errors"

var (
	ErrWalletNotFound = errors.New("wallet not found")
	ErrUserNotFound = errors.New("user not found")
)
