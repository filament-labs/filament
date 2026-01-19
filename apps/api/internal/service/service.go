package service

import (
	"github.com/codemaestro64/filament/apps/api/internal/repository"
	"github.com/codemaestro64/filament/libs/filwallet"
)

type Service struct {
	User UserService
}

func New(
	repo *repository.Repository,
	walletMgr *filwallet.Manager,
) *Service {
	return &Service{
		User: newUserService(repo, walletMgr),
	}
}
