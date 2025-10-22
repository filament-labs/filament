package service

import "github.com/filament-labs/filament/internal/repository"

type TransactionService interface {
}

type transactionService struct {
	transactionRepo repository.TransactionRepo
}

func NewTransactionService(repo *repository.Repository) TransactionService {
	return &transactionService{
		transactionRepo: repo.Transaction,
	}
}
