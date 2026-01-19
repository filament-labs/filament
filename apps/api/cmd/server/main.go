package main

import (
	"log"

	"github.com/codemaestro64/filament/apps/api/internal/app"
	"github.com/codemaestro64/filament/apps/api/internal/config"
	"github.com/spf13/cobra"
)

var Environment = config.Development

func main() {
	Execute()
}

func run(cmd *cobra.Command, args []string) {
	if err := app.Run(Environment); err != nil {
		log.Fatal(err)
	}
}
