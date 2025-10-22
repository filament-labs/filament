package service

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/filament-labs/filament/internal/dto"
	"github.com/filament-labs/filament/internal/model"
	"github.com/filament-labs/filament/internal/repository"
	"golang.org/x/crypto/pbkdf2"
)

type WalletService interface {
	LoadWallets(ctx context.Context, req dto.GetWalletsRequest) (*dto.GetWalletsResponse, error)
}

type walletService struct {
	walletRepo repository.WalletRepo
	wallets    []model.Wallet
	masterKey  []byte
}

func NewWalletService(repo *repository.Repository) WalletService {
	return &walletService{
		walletRepo: repo.Wallet,
	}
}

func (s *walletService) LoadWallets(ctx context.Context, req dto.GetWalletsRequest) (*dto.GetWalletsResponse, error) {
	wallets, err := s.walletRepo.GetWallets()
	if err != nil {
		return nil, err
	}

	res := &dto.GetWalletsResponse{
		Wallets: make([]dto.GetWalletResponse, len(wallets)),
	}
	for idx, wallet := range wallets {
		res.Wallets[idx] = dto.GetWalletResponse{
			ID: wallet.ID,
		}
	}
	s.wallets = wallets

	return res, nil
}

func (s *walletService) UnlockWallets(ctx context.Context, req dto.UnlockWalletsRequest) (*dto.UnlockWalletsResponse, error) {
	s.masterKey = pbkdf2.Key([]byte(req.Password), []byte("filament-global-salt"), 4096, 36, sha256.New)
	if s.wallets == nil {
		_, err := s.LoadWallets(ctx, dto.GetWalletsRequest{})
		if err != nil {
			return nil, err
		}
	}

	if len(s.wallets) == 0 {
		return &dto.UnlockWalletsResponse{}, nil
	}

	// TODO unlock
	if _, err := s.wallets[0].GetPrivateKey(s.masterKey); err != nil {
		return nil, fmt.Errorf("")
	}

	return &dto.UnlockWalletsResponse{}, nil
}
