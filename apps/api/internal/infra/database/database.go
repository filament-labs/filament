package database

import (
	"context"
	"fmt"
	"net/url"
	"path/filepath"

	"entgo.io/ent/dialect"
	"github.com/codemaestro64/filament/apps/api/internal/config"
	"github.com/codemaestro64/filament/apps/api/internal/database/orm"

	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	client *orm.Client
}

func New(cfg config.DatabaseConfig, dataDir string, env config.Env) (*Database, error) {
	var dsn string
	driver := cfg.Driver

	switch cfg.Driver {
	case "sqlite", "sqlite3":
		driver = dialect.SQLite
		dbPath := filepath.Join(dataDir, cfg.Name)
		if env == config.Development {
			dbPath += "_dev"
		}
		// PRAGMAs are part of the DSN; Ent handles this during Open.
		dsn = fmt.Sprintf("file:%s.db?_pragma=foreign_keys(1)&_pragma=journal_mode(WAL)", dbPath)

	case "postgres", "postgresql":
		driver = dialect.Postgres
		password := url.QueryEscape(cfg.Password)
		dsn = fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
			cfg.User, password, cfg.Host, cfg.Port, cfg.Name)

	case "mysql":
		driver = dialect.MySQL
		dsn = fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True",
			cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name)

	default:
		return nil, fmt.Errorf("unsupported driver: %s", cfg.Driver)
	}

	client, err := orm.Open(driver, dsn)
	if err != nil {
		return nil, fmt.Errorf("orm open: %w", err)
	}

	return &Database{client: client}, nil
}

func (d *Database) Name() string { return "database" }

func (d *Database) GetClient() *orm.Client {
	return d.client
}

func (d *Database) Start(ctx context.Context) error {
	// Verify the connection.
	if d.client == nil {
		return fmt.Errorf("database client not initialized")
	}

	return nil
}

func (d *Database) Shutdown(ctx context.Context) error {
	if d.client == nil {
		return nil
	}
	return d.client.Close()
}
