package wallet

import "context"

type Wallet struct {
	ID      string            `json:"id"`
	Name    string            `json:"name"`
	KeyJSON []byte            `json:"key"`   // encrypted key from keystore
	Addrs   map[string]string `json:"addrs"` // type => address
	Meta    map[string]string `json:"meta"`  // optional meta data
}

// WalletStore is the DB-agnostic interface the caller must implement.
type WalletStore interface {
	SaveWallet(ctx context.Context, w Wallet) error
	GetWallet(ctx context.Context, id string) (*Wallet, error)
	ListWallets(ctx context.Context) ([]Wallet, error)
	DeleteWallet(ctx context.Context, id string) error
}
