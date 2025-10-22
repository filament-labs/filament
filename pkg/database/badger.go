package database

import (
	"fmt"
	"path"
	"strings"

	"github.com/dgraph-io/badger/v4"
	"github.com/filament-labs/filament/pkg/util"
)

type Database struct {
	db *badger.DB
}

// New opens or creates a new BadgerDB instance
func NewBadgerDB(appDataDir string) (Store, error) {
	dbPath := path.Join(appDataDir, "db")
	opts := badger.DefaultOptions(dbPath)

	db, err := badger.Open(opts)
	if err != nil {
		return nil, fmt.Errorf("error opening or creating database `%s`: %w", dbPath, err)
	}

	return &Database{db: db}, nil
}

// Save stores any Go value (struct, slice, map, etc.) under the given key
func (d *Database) Save(key string, value any) error {
	data, err := util.Encode(value)
	if err != nil {
		return fmt.Errorf("error encoding value: %w", err)
	}

	if err := d.db.Update(func(txn *badger.Txn) error {
		return txn.Set([]byte(key), data)
	}); err != nil {
		return fmt.Errorf("error saving to database: %w", err)
	}

	return nil
}

// Get retrieves a value by key and decodes it into dest
func (d *Database) Get(key string, dest any) error {
	return d.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("error getting item from database: %w", err)
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("error copying value: %w", err)
		}

		if err := util.Decode(val, dest); err != nil {
			return fmt.Errorf("error decoding value: %w", err)
		}

		return nil
	})
}

func (d *Database) GetMany(prefix string, decodeFunc func([]byte) (any, error)) ([]any, error) {
	var results []any
	err := d.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := string(item.Key())
			if !strings.HasPrefix(k, prefix) {
				continue
			}

			val, err := item.ValueCopy(nil)
			if err != nil {
				return err
			}

			obj, err := decodeFunc(val)
			if err != nil {
				return err
			}
			results = append(results, obj)
		}
		return nil
	})
	return results, err
}

// Delete removes a key from the database
func (d *Database) Delete(key string) error {
	if err := d.db.Update(func(txn *badger.Txn) error {
		return txn.Delete([]byte(key))
	}); err != nil {
		return fmt.Errorf("error deleting key: %w", err)
	}

	return nil
}

// Update modifies an existing value atomically using the provided update function
func (d *Database) Update(key string, updateFn func(current any) (any, error)) error {
	return d.db.Update(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(key))
		if err != nil {
			return fmt.Errorf("error getting key for update: %w", err)
		}

		val, err := item.ValueCopy(nil)
		if err != nil {
			return fmt.Errorf("error copying value for update: %w", err)
		}

		var current any
		if err := util.Decode(val, &current); err != nil {
			return fmt.Errorf("error decoding current value: %w", err)
		}

		newVal, err := updateFn(current)
		if err != nil {
			return fmt.Errorf("update function error: %w", err)
		}

		data, err := util.Encode(newVal)
		if err != nil {
			return fmt.Errorf("error encoding updated value: %w", err)
		}

		return txn.Set([]byte(key), data)
	})
}

// Close closes the database
func (d *Database) Close() error {
	return d.db.Close()
}
