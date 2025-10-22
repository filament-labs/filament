package database

type Store interface {
	Save(key string, value any) error
	Get(key string, dest any) error
	GetMany(prefix string, decodeFunc func([]byte) (any, error)) ([]any, error)
	Delete(key string) error
	Update(key string, updateFn func(current any) (any, error)) error
	Close() error
}
