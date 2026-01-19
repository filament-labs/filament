package filwallet

import (
	"context"

	"github.com/codemaestro64/filament/libs/filwallet/address"
	"github.com/codemaestro64/filament/libs/filwallet/wallet"
)

type SaveWalletParams struct {
	KeyJSON       []byte
	EncryptedSeed []byte
	Addresses     []address.Address
	Name          string
	Salt          []byte
	Password      string
}

type Store interface {
	CountWallets(ctx context.Context) (int, error)
	GetWallets(ctx context.Context) ([]*wallet.Wallet, error)
	FindWallet(ctx context.Context, walletID int) (*wallet.Wallet, error)
	SaveWallet(ctx context.Context, p SaveWalletParams) (*wallet.Wallet, error)
	DeleteWallet(ctx context.Context, walletID int) error
}
