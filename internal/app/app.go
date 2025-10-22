package app

import (
	"fmt"

	"github.com/filament-labs/filament/internal/database"
	"github.com/filament-labs/filament/internal/repository"
	"github.com/filament-labs/filament/internal/service"
	"github.com/filament-labs/filament/pkg/util"
)

type App interface {
	InitHandlers(srvc *service.Service)
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

	// init app
	app := appFunc()

	repo := repository.New(db)
	srvc := service.New(repo)

	app.InitHandlers(srvc)

	return app.Run()
}
