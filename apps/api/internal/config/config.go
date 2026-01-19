package config

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/viper"
)

const AppName = "filament"

const (
	// App
	KeyEnvironment = "server.environment"
	KeyServerPort  = "server.port"
	KeyServerHost  = "server.host"

	// Database
	KeyDBDriver   = "database.driver"
	KeyDBHost     = "database.host"
	KeyDBPort     = "database.port"
	KeyDBName     = "database.name"
	KeyDBUser     = "database.user"
	KeyDBPassword = "database.password"

	// Logs
	KeyLogLevel      = "log.level"
	KeyLogMaxSize    = "log.max_size"
	KeyLogMaxBackups = "log.max_backups"
	KeyLogMaxAge     = "log.max_age"
)

// App Config structs
type ServerConfig struct {
	Port        int
	Host        string
	Environment Env
}

type DatabaseConfig struct {
	Driver   string
	Host     string
	Port     int
	Name     string
	User     string
	Password string
}

type LogConfig struct {
	Level      string
	MaxSize    int
	MaxBackups int
	MaxAge     int
}

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Log      LogConfig
}

// Load config from Viper
func Load(env Env) (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Environment: env,
			Port:        viper.GetInt(KeyServerPort),
			Host:        viper.GetString(KeyServerHost),
		},
		Database: DatabaseConfig{
			Driver:   viper.GetString(KeyDBDriver),
			Host:     viper.GetString(KeyDBHost),
			Port:     viper.GetInt(KeyDBPort),
			Name:     viper.GetString(KeyDBName),
			User:     viper.GetString(KeyDBUser),
			Password: viper.GetString(KeyDBPassword),
		},
		Log: LogConfig{
			Level:      viper.GetString(KeyLogLevel),
			MaxSize:    viper.GetInt(KeyLogMaxSize),
			MaxBackups: viper.GetInt(KeyLogMaxBackups),
			MaxAge:     viper.GetInt(KeyLogMaxAge),
		},
	}

	if err := validate(cfg); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return cfg, nil
}

func validate(cfg *Config) error {
	var errs []error

	driver := strings.ToLower(strings.TrimSpace(cfg.Database.Driver))

	if driver == "" {
		errs = append(errs, fmt.Errorf("database.driver is required"))
	}

	switch driver {
	case "sqlite":
		if strings.TrimSpace(cfg.Database.Name) == "" {
			errs = append(errs, fmt.Errorf("database.name is required when using sqlite"))
		}

	default:
		if strings.TrimSpace(cfg.Database.Host) == "" {
			errs = append(errs, fmt.Errorf("database.host is required when using %s", driver))
		}
		if cfg.Database.Port == 0 {
			errs = append(errs, fmt.Errorf("database.port is required when using %s", driver))
		}
		if strings.TrimSpace(cfg.Database.Name) == "" {
			errs = append(errs, fmt.Errorf("database.name is required when using %s", driver))
		}
		if strings.TrimSpace(cfg.Database.User) == "" {
			errs = append(errs, fmt.Errorf("database.user is required when using %s", driver))
		}
	}

	if len(errs) > 0 {
		return errors.Join(errs...)
	}

	return nil
}
