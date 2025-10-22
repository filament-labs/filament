package repository

import "github.com/dgraph-io/badger/v4"

type Repository struct {
	Wallet      WalletRepo
	Setting     SettingRepo
	Transaction TransactionRepo
}

func New(db *badger.DB) *Repository {
	return &Repository{
		Wallet:      NewWalletRepo(db),
		Setting:     NewSettingRepo(db),
		Transaction: NewTransactionRepo(db),
	}
}
