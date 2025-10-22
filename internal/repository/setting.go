package repository

import "github.com/dgraph-io/badger/v4"

type SettingRepo interface {
}

type settingRepo struct {
	db *badger.DB
}

func NewSettingRepo(db *badger.DB) SettingRepo {
	return &settingRepo{
		db: db,
	}
}
