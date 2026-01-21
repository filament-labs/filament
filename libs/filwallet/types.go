package filwallet

import (
	"errors"
)

var (
	ErrNotFound          = errors.New("wallet not found")
	ErrSessionExpired    = errors.New("session expired")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidSeedPhrase = errors.New("invalid seed phrase")
	ErrInvalidWalletName = errors.New("invalid wallet name")
)
