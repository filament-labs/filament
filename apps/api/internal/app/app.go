package app

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/codemaestro64/filament/apps/api/internal/config"
	"github.com/codemaestro64/filament/apps/api/internal/domain"
	"github.com/codemaestro64/filament/apps/api/internal/infra/database"
	"github.com/codemaestro64/filament/apps/api/internal/infra/logger"
	"github.com/codemaestro64/filament/apps/api/internal/repository"
	"github.com/codemaestro64/filament/apps/api/internal/server"
	"github.com/codemaestro64/filament/apps/api/internal/service"
	"github.com/codemaestro64/filament/apps/api/pkg/util"
	"github.com/codemaestro64/filament/libs/filwallet"
	"github.com/rs/zerolog/log"
)

const (
	// DefaultShutdownTimeout ensures we don't hang forever on exit.
	DefaultShutdownTimeout = 20 * time.Second
	// DefaultDirPerm: owner full access, group read, others none.
	DefaultDirPerm = 0750
	// FilePerm: owner read/write, group read.
	FilePerm = 0640
	// SettingsFileName is the Layer 0 bootstrap file.
	SettingsFileName = "settings.json"
	// DefaultNetwork is used for first-time installs.
	//DefaultNetwork = config.CalibrationNet
)

// Runnable defines the contract for components with a controlled lifecycle.
type Runnable interface {
	Name() string
	Start(ctx context.Context) error
	Shutdown(ctx context.Context) error
}

// Run initializes and starts the application infrastructure and services.
func Run(env config.Env) error {
	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer cancel()

	baseDir, err := util.AppDataDir(config.AppName)
	if err != nil {
		return fmt.Errorf("could not determine app data directory: %w", err)
	}

	network, err := bootstrap(util.CalibrationNet, baseDir)
	if err != nil {
		return fmt.Errorf("bootstrap failure: %w", err)
	}

	envName := strings.ToLower(env.String())
	dataDir := filepath.Join(baseDir, envName, network.String())

	if err := os.MkdirAll(dataDir, DefaultDirPerm); err != nil {
		return fmt.Errorf("create data directory structure: %w", err)
	}

	cfg, err := config.Load(env)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	logger.New(cfg.Log, env, dataDir)

	log.Info().
		Str("env", envName).
		Str("network", network.String()).
		Str("data_dir", dataDir).
		Msg("application bootstrap complete")

	db, err := database.New(cfg.Database, dataDir, env)
	if err != nil {
		return fmt.Errorf("init database: %w", err)
	}

	repo := repository.New(db.GetClient())

	// Intialize wallet manager
	walletMgr, err := filwallet.NewManager(ctx, repo.Wallet, &filwallet.Config{
		Network:        network,
		SessionTimeout: cfg.Server.SessionTimeout,
		RPCEndpoint:    "",
		RPCToken:       "",
		DataDir:        dataDir,
	})

	srvc := service.New(repo, walletMgr)

	srvr, err := server.New(srvc, cfg.Server, cancel)
	if err != nil {
		return fmt.Errorf("init server: %w", err)
	}

	components := []Runnable{db, srvr}
	return runWithGracefulShutdown(ctx, DefaultShutdownTimeout, components)
}

// bootstrap handles first-use by creating settings.json or reading the existing network preference.
func bootstrap(defaultNetwork util.Network, baseDir string) (util.Network, error) {
	settingsPath := filepath.Join(baseDir, SettingsFileName)

	// Check if settings file exists
	if _, err := os.Stat(settingsPath); errors.Is(err, os.ErrNotExist) {
		// Ensure baseDir exists so we can write the file
		if err := os.MkdirAll(baseDir, DefaultDirPerm); err != nil {
			return "", fmt.Errorf("create base dir: %w", err)
		}

		// Initialize with configured network
		s := domain.Settings{Network: defaultNetwork}
		bytes, _ := json.MarshalIndent(s, "", "  ")
		if err := os.WriteFile(settingsPath, bytes, FilePerm); err != nil {
			return "", fmt.Errorf("create default settings: %w", err)
		}
		return defaultNetwork, nil
	}

	// Read existing settings
	data, err := os.ReadFile(settingsPath)
	if err != nil {
		return "", fmt.Errorf("read settings: %w", err)
	}

	var s domain.Settings
	if err := json.Unmarshal(data, &s); err != nil || s.Network == "" {
		// Fallback to default network if file is corrupted
		return defaultNetwork, nil
	}

	return s.Network, nil
}

// runWithGracefulShutdown manages startup and ensures reverse-order cleanup on exit.
func runWithGracefulShutdown(ctx context.Context, timeout time.Duration, components []Runnable) error {
	started := make([]Runnable, 0, len(components))

	for _, c := range components {
		log.Info().Str("component", c.Name()).Msg("starting")
		if err := c.Start(ctx); err != nil {
			log.Error().Err(err).Str("component", c.Name()).Msg("startup failed")
			shutdown(timeout, started)
			return fmt.Errorf("component %s failed: %w", c.Name(), err)
		}
		started = append(started, c)
	}

	log.Info().Msg("application is running")

	// Wait for OS signal (Ctrl+C) or internal fatal error (via context cancel)
	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	shutdown(timeout, started)
	return nil
}

// shutdown performs clean closure of all started components in reverse order.
func shutdown(timeout time.Duration, components []Runnable) {
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	for i := len(components) - 1; i >= 0; i-- {
		c := components[i]
		log.Info().Str("component", c.Name()).Msg("shutting down")
		if err := c.Shutdown(shutdownCtx); err != nil {
			log.Error().Err(err).Str("component", c.Name()).Msg("shutdown error")
		}
	}
	log.Info().Msg("graceful shutdown completed")
}
