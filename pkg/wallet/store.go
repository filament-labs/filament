package wallet

import (
	"fmt"

	"github.com/dgraph-io/badger/v4"
	"github.com/filament-labs/filament/pkg/util"
)

// store manages persistence of Wallet objects in BadgerDB using a key prefix.
type store struct {
	db *badger.DB
}

const walletKeyPrefix = "wallet:"

// newStore creates a new store instance bound to the given Badger DB.
func newStore(db *badger.DB) *store {
	return &store{db: db}
}

// walletKey constructs the full DB key for a wallet ID.
func walletKey(id string) []byte {
	return append([]byte(walletKeyPrefix), []byte(id)...)
}

// listWallets loads all wallets from the database.
func (s *store) listWallets() ([]Wallet, error) {
	var wallets []Wallet

	err := s.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		prefix := []byte(walletKeyPrefix)
		for it.Seek(prefix); it.ValidForPrefix(prefix); it.Next() {
			item := it.Item()

			var w Wallet
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

// getWalletByID retrieves a wallet by its ID.
func (s *store) getWalletByID(id string) (*Wallet, error) {
	var wallet Wallet

	err := s.db.View(func(txn *badger.Txn) error {
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

// saveWallet persists a wallet to the database.
func (s *store) saveWallet(wallet Wallet) error {
	encoded, err := util.Encode(wallet)
	if err != nil {
		return fmt.Errorf("failed to encode wallet: %w", err)
	}

	err = s.db.Update(func(txn *badger.Txn) error {
		return txn.Set(walletKey(wallet.ID), encoded)
	})
	if err != nil {
		return fmt.Errorf("failed to save wallet %s: %w", wallet.ID, err)
	}

	return nil
}
