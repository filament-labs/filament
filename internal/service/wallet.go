package service

import (
	"github.com/filament-labs/filament/internal/repository"
	"github.com/filament-labs/filament/pkg/wallet"
)

type WalletService interface {
}

type walletService struct {
	walletManager *wallet.Manager
	walletRepo    repository.WalletRepo
}

func NewWalletService(repo *repository.Repository, walletManager *wallet.Manager) WalletService {
	return &walletService{
		walletRepo:    repo.Wallet,
		walletManager: walletManager,
	}
}

func (s *walletService) GetWalletCount() {

}
