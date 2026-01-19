package main

import (
	"fmt"
	"os"

	"github.com/codemaestro64/filament/apps/api/internal/config"
	"github.com/codemaestro64/filament/apps/api/pkg/util"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	cfgFile string
	rootCmd = &cobra.Command{
		Use:   "api-server",
		Short: "Filament wallet API server",
		Long:  "Filament wallet connectRPC API server for the filament wallet web app",
		Run:   run,
	}
)

func init() {
	cobra.OnInitialize(initConfig)

	/// Server flags
	rootCmd.Flags().Int("port", 0, "Server port")
	rootCmd.Flags().String("host", "0.0.0.0", "Server host")
	rootCmd.Flags().String("network", "calibration", "Network")
	rootCmd.Flags().Int64("session_timeout", 30, "Session Timeout (mins)")

	// Database flags
	rootCmd.Flags().String("db-driver", "sqlite", "Database driver")
	rootCmd.Flags().String("db-host", "", "Database host")
	rootCmd.Flags().Int("db-port", 3306, "Database port")
	rootCmd.Flags().String("db-name", "filament-db", "Database name")
	rootCmd.Flags().String("db-user", "root", "Database user")
	rootCmd.Flags().String("db-password", "", "Database password")

	// Log flags
	rootCmd.Flags().String("log-level", "info", "Log level")
	rootCmd.Flags().Int("log-max-size", 500, "Log max size")
	rootCmd.Flags().Int("log-max-backups", 3, "Log max backups")
	rootCmd.Flags().Int("log-max-age", 28, "Log max age")

	// Bind flags → config keys
	_ = viper.BindPFlag(config.KeyServerPort, rootCmd.Flags().Lookup("port"))
	_ = viper.BindPFlag(config.KeyServerHost, rootCmd.Flags().Lookup("host"))
	_ = viper.BindPFlag(config.KeyNetwork, rootCmd.Flags().Lookup("network"))
	_ = viper.BindPFlag(config.KeySessionTimeout, rootCmd.Flags().Lookup("session_timeout"))

	_ = viper.BindPFlag(config.KeyDBDriver, rootCmd.Flags().Lookup("db-driver"))
	_ = viper.BindPFlag(config.KeyDBHost, rootCmd.Flags().Lookup("db-host"))
	_ = viper.BindPFlag(config.KeyDBPort, rootCmd.Flags().Lookup("db-port"))
	_ = viper.BindPFlag(config.KeyDBName, rootCmd.Flags().Lookup("db-name"))
	_ = viper.BindPFlag(config.KeyDBUser, rootCmd.Flags().Lookup("db-user"))
	_ = viper.BindPFlag(config.KeyDBPassword, rootCmd.Flags().Lookup("db-password"))

	_ = viper.BindPFlag(config.KeyLogLevel, rootCmd.Flags().Lookup("log-level"))
	_ = viper.BindPFlag(config.KeyLogMaxSize, rootCmd.Flags().Lookup("log-max-size"))
	_ = viper.BindPFlag(config.KeyLogMaxBackups, rootCmd.Flags().Lookup("log-max-backups"))
	_ = viper.BindPFlag(config.KeyLogMaxAge, rootCmd.Flags().Lookup("log-max-age"))
}

func initConfig() {
	// Server
	viper.SetDefault(config.KeyServerPort, 0)
	viper.SetDefault(config.KeyServerHost, "0.0.0.0")
	viper.SetDefault(config.KeyNetwork, util.CalibrationNet.String())
	viper.SetDefault(config.KeySessionTimeout, 30)

	// Database
	viper.SetDefault(config.KeyDBDriver, "sqlite")
	viper.SetDefault(config.KeyDBHost, "")
	viper.SetDefault(config.KeyDBPort, 3306)
	viper.SetDefault(config.KeyDBName, "filament-db")
	viper.SetDefault(config.KeyDBUser, "root")
	viper.SetDefault(config.KeyDBPassword, "")

	// Logs
	viper.SetDefault(config.KeyLogLevel, "info")
	viper.SetDefault(config.KeyLogMaxSize, 500)
	viper.SetDefault(config.KeyLogMaxBackups, 3)
	viper.SetDefault(config.KeyLogMaxAge, 28)

	// Config file
	if cfgFile != "" {
		viper.SetConfigFile(cfgFile)
	} else {
		envFile := getEnvFilePath()
		if envFile != "" {
			if _, err := os.Stat(envFile); err == nil {
				viper.AddConfigPath(".")
				viper.SetConfigName(envFile)
				viper.SetConfigType("env")
			}
		}
	}

	// Read config file (optional)
	if err := viper.ReadInConfig(); err == nil {
		fmt.Println("Using config file:", viper.ConfigFileUsed())
	}

	// Environment variables override config
	viper.AutomaticEnv()

	// Flags override everything
	_ = viper.BindPFlags(rootCmd.Flags())
	_ = viper.BindPFlags(rootCmd.PersistentFlags())
}

func getEnvFilePath() string {
	switch Environment {
	case config.Production:
		return ".env.production"
	case config.Development:
		return ".env"
	default:
		return ""
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
