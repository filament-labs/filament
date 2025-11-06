package main

import (
	"log"

	"github.com/filament-labs/filament/internal/app"
	"github.com/filament-labs/filament/internal/wsapi"
)

func main() {
	err := app.Run("Filament", wsapi.New)
	if err != nil {
		log.Fatal(err)
	}
}
