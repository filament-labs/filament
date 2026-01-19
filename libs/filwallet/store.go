package filwallet

import (
	"context"

	pbv1 "github.com/codemaestro64/filament/libs/proto/gen/go/v1"
)

type SaveWalletParams struct {
	KeyJSON       []byte
	EncryptedSeed []byte
	Addresses     []*pbv1.Address
	Name          string
	Salt          []byte
	Password      string
}

type Store interface {
	CountWallets(ctx context.Context) (int, error)
	GetWallets(ctx context.Context) ([]*Wallet, error)
	FindWallet(ctx context.Context, walletID int) (*Wallet, error)
	SaveWallet(ctx context.Context, p SaveWalletParams) (*Wallet, error)
	DeleteWallet(ctx context.Context, walletID int) error
}
