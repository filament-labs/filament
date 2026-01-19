package repository

import "github.com/codemaestro64/filament/apps/api/internal/database/orm"

type Repository struct {
}

func New(dbClient *orm.Client) *Repository {
	return &Repository{}
}
