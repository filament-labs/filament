package service

import (
	"github.com/filament-labs/filament/internal/repository"
)

type Service struct {
}

func New(repo *repository.Repository) *Service {
	return &Service{}
}
