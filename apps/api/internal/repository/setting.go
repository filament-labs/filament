package repository

import (
	"context"

	"github.com/codemaestro64/filament/apps/api/internal/database/orm"
)

type SettingRepo interface {
}

type settingRepo struct {
	db *orm.Client
}

func newSettingRepo(db *orm.Client) SettingRepo {
	return &settingRepo{
		db: db,
	}
}

func (s *settingRepo) GetSettings(ctx context.Context) (*orm.Setting, error) {
	return nil, nil
}
