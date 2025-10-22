package database

import (
	"path"

	"github.com/dgraph-io/badger/v4"
)

func Open(dataDir string) (*badger.DB, error) {
	dbPath := path.Join(dataDir, "db")
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func Close(db *badger.DB) error {
	if db != nil {
		if err := db.Close(); err != nil {
			return err
		}
	}

	return nil
}
