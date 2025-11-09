package service

import (
	"github.com/filament-labs/filament/internal/repository"
	"github.com/filament-labs/filament/pkg/wallet"
)

type Service struct {
	Wallet WalletService
}

func New(repo *repository.Repository, walletManager *wallet.Manager) *Service {
	return &Service{
		Wallet: NewWalletService(repo, walletManager),
	}
}
