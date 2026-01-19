package filwallet

import (
	"errors"
	"time"

	pbv1 "github.com/codemaestro64/filament/libs/proto/gen/go/v1"
)

type Wallet struct {
	ID               int
	IsDefault        bool
	Name             string
	SeedPhrase       string
	Addresses        []*pbv1.Address
	Salt             []byte
	EncryptedKeyJSON []byte
	EncryptedSeed    []byte
	CreatedAt        time.Time
}

type Transaction struct {
}

var (
	ErrNotFound          = errors.New("wallet not found")
	ErrSessionExpired    = errors.New("session expired")
	ErrWalletExists      = errors.New("wallet exists")
	ErrInvalidPassword   = errors.New("invalid password")
	ErrInvalidSeedPhrase = errors.New("invalid seed phrase")
	ErrInvalidWalletName = errors.New("invalid wallet name")
)
