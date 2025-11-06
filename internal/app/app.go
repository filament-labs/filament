package app

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/filament-labs/filament/internal/database"
	"github.com/filament-labs/filament/internal/repository"
	"github.com/filament-labs/filament/internal/service"
	"github.com/filament-labs/filament/pkg/util"
	"github.com/filament-labs/filament/pkg/wallet"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

// App defines the interface for application lifecycle management.
type App interface {
	Bootstrap(srvc *service.Service) error
	Run() error
}

// appDataDir returns the platform-appropriate directory for persistent app data (DB, logs).
// - Linux:   $XDG_DATA_HOME/<appName> or ~/.local/share/<appName>
// - macOS:   ~/Library/Application Support/<appName>
// - Windows: %LOCALAPPDATA%\<appName>
func appDataDir(appName string) (string, error) {
	var dir string

	switch runtime.GOOS {
	case "linux":
		if xdg := os.Getenv("XDG_DATA_HOME"); xdg != "" {
			dir = filepath.Join(xdg, appName)
		} else {
			home, err := os.UserHomeDir()
			if err != nil {
				return "", fmt.Errorf("failed to get home directory: %w", err)
			}
			dir = filepath.Join(home, ".local", "share", appName)
		}
	case "darwin":
		home, err := os.UserHomeDir()
		if err != nil {
			return "", fmt.Errorf("failed to get home directory: %w", err)
		}
		dir = filepath.Join(home, "Library", "Application Support", appName)
	case "windows":
		localAppData := os.Getenv("LOCALAPPDATA")
		if localAppData == "" {
			return "", fmt.Errorf("LOCALAPPDATA environment variable not set")
		}
		dir = filepath.Join(localAppData, appName)
	default:
		base, err := os.UserCacheDir()
		if err != nil {
			return "", fmt.Errorf("failed to get cache directory: %w", err)
		}
		dir = filepath.Join(base, appName)
	}

	if err := os.MkdirAll(dir, 0700); err != nil {
		return "", fmt.Errorf("failed to create data directory %s: %w", dir, err)
	}

	return dir, nil
}

func configureLogger() {
	zerolog.TimeFieldFormat = "15:04:05"
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
}

// Run initializes and starts the application with proper logging.
func Run(appName string, appFunc func() App) error {
	configureLogger()
	log.Info().Msgf("Starting %s...", appName)

	// Get app data directory
	log.Debug().Msg("Resolving app data directory")
	dataDir, err := appDataDir(util.HyphenateAndLower(appName))
	if err != nil {
		log.Error().Err(err).Msg("Failed to resolve app data directory")
		return fmt.Errorf("failed to resolve app data directory: %w", err)
	}
	log.Info().Str("path", dataDir).Msg("App data directory resolved")

	// Initialize database
	log.Debug().Msg("Opening database connection")
	db, err := database.Open(dataDir)
	if err != nil {
		log.Error().Err(err).Str("path", dataDir).Msg("Failed to open database")
		return fmt.Errorf("failed to open database: %w", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			log.Warn().Err(err).Msg("Failed to close database cleanly")
		} else {
			log.Debug().Msg("Database connection closed")
		}
	}()
	log.Info().Msg("Database connection established")

	// Initialize wallet manager
	log.Debug().Msg("Initializing wallet manager")
	walletManager, err := wallet.NewManager()
	if err != nil {
		log.Error().Err(err).Msg("Failed to initialize wallet manager")
		return fmt.Errorf("failed to initialize wallet manager: %w", err)
	}
	log.Info().Msg("Wallet manager initialized")

	// Set up app components
	log.Debug().Msg("Setting up repository and service")
	repo := repository.New(db)
	srvc := service.New(repo, walletManager)
	log.Info().Msg("Repository and service initialized")

	// Bootstrap application
	app := appFunc()
	log.Debug().Msg("Bootstrapping application")
	if err := app.Bootstrap(srvc); err != nil {
		log.Error().Err(err).Msg("Failed to bootstrap application")
		return fmt.Errorf("failed to bootstrap application: %w", err)
	}
	log.Info().Msg("Application bootstrapped successfully")

	// Run application
	log.Info().Msg("Running application")
	if err := app.Run(); err != nil {
		log.Error().Err(err).Msg("Application failed to run")
		return fmt.Errorf("application run failed: %w", err)
	}
	log.Info().Msg("Application shutdown gracefully")

	return nil
}
