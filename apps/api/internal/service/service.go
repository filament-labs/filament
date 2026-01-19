package service

import "github.com/codemaestro64/filament/apps/api/internal/repository"

type Service struct {
}

func New(repo *repository.Repository) *Service {
	return &Service{}
}
