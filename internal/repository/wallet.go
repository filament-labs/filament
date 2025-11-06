package repository

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/filament-labs/filament/pkg/util"
	"github.com/filament-labs/filament/pkg/wallet"
)

type WalletRepo interface {
	GetWallets() ([]*wallet.Wallet, error)
	Find(id string) (*wallet.Wallet, error)
	Save(wallet *wallet.Wallet) error
}

type walletRepo struct {
	db *badger.DB
}

const walletKeyPrefix = "wallet:"

func NewWalletRepo(db *badger.DB) WalletRepo {
	return &walletRepo{
		db: db,
	}
}

// walletKey constructs the full DB key for a wallet ID.
func walletKey(id string) []byte {
	return append([]byte(walletKeyPrefix), []byte(id)...)
}

func (r *walletRepo) GetWallets() ([]*wallet.Wallet, error) {
	var wallets []*wallet.Wallet

	err := r.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(walletKeyPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			var w *wallet.Wallet
			if err := item.Value(func(val []byte) error {
				return util.Decode(val, &w)
			}); err != nil {
				return fmt.Errorf("failed to decode wallet: %w", err)
			}

			wallets = append(wallets, w)
		}
		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("failed to list wallets: %w", err)
	}

	return wallets, nil
}

// Find retrieves a wallet by its ID.
func (r *walletRepo) Find(id string) (*wallet.Wallet, error) {
	var wallet wallet.Wallet

	err := r.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(walletKey(id))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			return util.Decode(val, &wallet)
		})
	})
	if err != nil {
		if err == badger.ErrKeyNotFound {
			return nil, nil // not an error, wallet doesnt exist
		}
		return nil, fmt.Errorf("failed to get wallet %s: %w", id, err)
	}

	return &wallet, nil
}

// Save persists a wallet to the database.
func (r *walletRepo) Save(wallet *wallet.Wallet) error {
	encoded, err := util.Encode(wallet)
	if err != nil {
		return fmt.Errorf("failed to encode wallet: %w", err)
	}

	err = r.db.Update(func(txn *badger.Txn) error {
		return txn.Set(walletKey(wallet.ID), encoded)
	})
	if err != nil {
		return fmt.Errorf("failed to save wallet %s: %w", wallet.ID, err)
	}

	return nil
}
