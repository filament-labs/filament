package repository

import (
	"fmt"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/filament-labs/filament/internal/model"
	"github.com/filament-labs/filament/pkg/util"
)

type WalletRepo interface {
	GetWallets() ([]model.Wallet, error)
}

type walletRepo struct {
	db *badger.DB
}

const (
	walletsKeyPrefix = "wallets_"
)

func NewWalletRepo(db *badger.DB) WalletRepo {
	return &walletRepo{
		db: db,
	}
}

func (r *walletRepo) GetWallets() ([]model.Wallet, error) {
	var wallets []model.Wallet
	err := r.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := string(item.Key())
			if !strings.HasPrefix(key, walletsKeyPrefix) {
				continue
			}

			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			var wallet model.Wallet
			err = util.Decode(val, &wallet)
			if err != nil {
				return err
			}

			wallets = append(wallets, wallet)
		}

		return nil
	})
	if err != nil {
		return nil, fmt.Errorf("error fetching wallets: %w", err)
	}

	return wallets, err
}
