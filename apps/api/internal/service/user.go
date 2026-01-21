package service

import (
	"context"

	"github.com/codemaestro64/filament/apps/api/internal/domain"
	"github.com/codemaestro64/filament/apps/api/internal/repository"
	"github.com/codemaestro64/filament/libs/filwallet"
	"github.com/rs/zerolog/log"
)

type UserService interface {
	GetBootstrap(ctx context.Context, req domain.GetBootstrapRequest) (*domain.GetBootstrapResponse, error)
}

type userService struct {
	walletMgr   *filwallet.Manager
	settingRepo repository.SettingRepo
	walletRepo  repository.WalletRepo
}

func newUserService(repo *repository.Repository, walletMgr *filwallet.Manager) UserService {
	return &userService{
		walletMgr:   walletMgr,
		settingRepo: repo.Setting,
		walletRepo:  repo.Wallet,
	}
}

func (s *userService) GetBootstrap(ctx context.Context, req domain.GetBootstrapRequest) (*domain.GetBootstrapResponse, error) {
	walletCount, err := s.walletMgr.WalletsCount(ctx)
	if err != nil {
		log.Error().Err(err).Msg("error fetching wallet count")
		return nil, domain.ErrInternalServer
	}

	resp := &domain.GetBootstrapResponse{
		WalletCount: walletCount,
		Settings:    domain.Settings{
			//Network: s.walletMgr.GetNetwork(),
		},
	}

	return resp, nil
}
