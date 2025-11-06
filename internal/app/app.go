package app

import (
	"fmt"

	"github.com/filament-labs/filament/internal/database"
	"github.com/filament-labs/filament/internal/repository"
	"github.com/filament-labs/filament/internal/service"
	"github.com/filament-labs/filament/pkg/util"
	"github.com/filament-labs/filament/pkg/wallet"
)

type App interface {
	Bootstrap(srvc *service.Service) error
	Run() error
}

func Run(appName string, appFunc func() App) error {
	// get data directory from name
	dataDir := util.HyphenateAndLower(appName)

	// initialize database connection
	db, err := database.Open(dataDir)
	if err != nil {
		return fmt.Errorf("error opening database connection: %w", err)
	}

	walletManager, err := wallet.NewManager(db)
	if err != nil {
		return fmt.Errorf("error initializing wallet manager: %w", err)
	}

	// init app
	app := appFunc()

	repo := repository.New(db)
	srvc := service.New(repo, walletManager)

	app.Bootstrap(srvc)

	return app.Run()
}
