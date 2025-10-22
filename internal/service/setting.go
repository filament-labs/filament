package service

import "github.com/filament-labs/filament/internal/repository"

type SettingService interface {
}

type settingService struct {
	settingRepo repository.SettingRepo
}

func NewSettingService(repo *repository.Repository) SettingService {
	return &settingService{
		settingRepo: repo.Setting,
	}
}
