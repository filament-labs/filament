package app

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/codemaestro64/filament/apps/api/internal/config"
	"github.com/codemaestro64/filament/apps/api/internal/infrastructure/database"
	"github.com/codemaestro64/filament/apps/api/internal/infrastructure/logger"
	"github.com/codemaestro64/filament/apps/api/internal/repository"
	"github.com/codemaestro64/filament/apps/api/internal/server"
	"github.com/codemaestro64/filament/apps/api/internal/service"
	"github.com/codemaestro64/filament/apps/api/pkg/util"
	"github.com/rs/zerolog/log"
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

	log.Info().Str("env", env.String()).Msg("bootstrap app")

	// Load Configuration
	log.Info().Msg("loading config...")
	cfg, err := config.Load(env)
	if err != nil {
		return fmt.Errorf("load config: %w", err)
	}

	// Setup Data Directory
	log.Info().Msg("setting up app data directory...")
	dataDir, err := util.AppDataDir(config.AppName)
	if err != nil {
		return fmt.Errorf("get app data directory: %w", err)
	}

	// drwxr-x--- owner full access, group read, others none.
	if err := os.MkdirAll(dataDir, 0750); err != nil {
		return fmt.Errorf("create data directory: %w", err)
	}

	// Initialize logger
	logger.New(cfg.Log, env, dataDir)

	log.Info().Msg("initializing database...")
	db, err := database.New(cfg.Database, dataDir, env)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}

	// Initialize Domain Layers
	log.Info().Msg("initializing repository and service layers...")
	repo := repository.New(db.GetClient())
	srvc := service.New(repo)

	log.Info().Msg("initializing server...")
	srvr, err := server.New(cfg.Server, srvc, cancel)
	if err != nil {
		return fmt.Errorf("init server: %w", err)
	}

	// Define Lifecycle Components
	// DB is first (dependency), Server is second (interface).
	components := []Runnable{
		db,
		srvr,
	}

	return runWithGracefulShutdown(ctx, 20*time.Second, components)
}

// runWithGracefulShutdown manages the sequential startup and reverse-order shutdown.
func runWithGracefulShutdown(ctx context.Context, timeout time.Duration, components []Runnable) error {
	started := make([]Runnable, 0, len(components))

	for _, component := range components {
		log.Info().Str("component", component.Name()).Msg("starting component")
		if err := component.Start(ctx); err != nil {
			log.Error().Err(err).Str("component", component.Name()).Msg("failed to start component")
			// If a component fails to start, we shut down whatever was already started.
			shutdown(timeout, started)
			return fmt.Errorf("start component %s: %w", component.Name(), err)
		}
		started = append(started, component)
	}

	log.Info().Msg("application is running")

	// Block until the context is canceled (via OS signal or parent context)
	<-ctx.Done()
	log.Info().Msg("shutdown signal received")

	// Gracefully shut down in reverse order
	shutdown(timeout, started)
	return nil
}

// shutdown performs the cleanup logic for all started components.
func shutdown(timeout time.Duration, components []Runnable) {
	// Context for shutdown ensures one component doesn't hang the entire process exit.
	shutdownCtx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	// Shutdown in reverse order: Server closes first to stop traffic, DB closes last.
	for i := len(components) - 1; i >= 0; i-- {
		component := components[i]
		log.Info().Str("component", component.Name()).Msg("shutting down component")
		if err := component.Shutdown(shutdownCtx); err != nil {
			log.Error().
				Err(err).
				Str("component", component.Name()).
				Msg("shutdown error")
		}
	}
	log.Info().Msg("graceful shutdown completed")
}
