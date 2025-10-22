package repository

import "github.com/dgraph-io/badger/v4"

type TransactionRepo interface {
}

type transactionRepo struct {
	db *badger.DB
}

func NewTransactionRepo(db *badger.DB) TransactionRepo {
	return &transactionRepo{
		db: db,
	}
}
