package service

import (
	"github.com/filament-labs/filament/internal/repository"
	walletclient "github.com/filament-labs/filament/pkg/wallet_client"
)

type Service struct {
}

func New(repo *repository.Repository, walletClient *walletclient.WalletClient) *Service {
	return &Service{}
}
