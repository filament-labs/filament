package service

import (
	"context"

	"github.com/filament-labs/filament/internal/dto"
	"github.com/filament-labs/filament/internal/repository"
	"github.com/filament-labs/filament/pkg/wallet"
)

type WalletService interface {
	GetWallets(ctx context.Context) (*dto.WalletsResponse, error)
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

func (s *walletService) GetWallets(ctx context.Context) (*dto.WalletsResponse, error) {
	wallets, err := s.walletRepo.GetWallets()
	if err != nil {
		return nil, err
	}

	respWallets := make([]dto.WalletResponse, 0, len(wallets))
	for _, w := range wallets {
		respWallets = append(respWallets, dto.WalletResponse{
			ID:        w.ID,
			IsDefault: w.IsDefault,
			Name:      w.Name,
			Addresses: w.Addrs,
			CreatedAt: w.CreatedAt,
		})
	}

	return &dto.WalletsResponse{
		Locked:  true,
		Wallets: respWallets,
	}, nil
}
