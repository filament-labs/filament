package service

import "github.com/filament-labs/filament/internal/repository"

type TransactionService interface {
	StartTransactionListener(walletAddress string)
}

type transactionService struct {
	transactionRepo repository.TransactionRepo
}

func NewTransactionService(repo *repository.Repository) TransactionService {
	return &transactionService{
		transactionRepo: repo.Transaction,
	}
}

func (s *transactionService) StartTransactionListener(walletAddress string) {

}
