package repository

import "github.com/codemaestro64/filament/apps/api/internal/infra/database/orm"

type Repository struct {
	Setting SettingRepo
	Wallet  WalletRepo
}

func New(dbClient *orm.Client) *Repository {
	return &Repository{
		Setting: newSettingRepo(dbClient),
		Wallet:  newWalletRepo(dbClient),
	}
}
