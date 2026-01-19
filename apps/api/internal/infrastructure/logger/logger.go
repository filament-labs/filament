package logger

import (
	"os"
	"path"

	"github.com/codemaestro64/filament/apps/api/internal/config"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/rs/zerolog/pkgerrors"
	"gopkg.in/natefinch/lumberjack.v2"
)

func New(cfg config.LogConfig, env config.Env, dataDir string) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.ErrorStackMarshaler = pkgerrors.MarshalStack

	if env == config.Development {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr}).With().Timestamp().Logger()
		return
	}
	log.Logger = zerolog.New(newRollingFile(cfg, dataDir)).With().Timestamp().Logger()
}

func newRollingFile(cfg config.LogConfig, dataDir string) *lumberjack.Logger {
	logFilePath := path.Join(dataDir, "logs", "log.log")

	return &lumberjack.Logger{
		Filename:   logFilePath,
		MaxBackups: cfg.MaxBackups,
		MaxSize:    cfg.MaxSize,
		MaxAge:     cfg.MaxAge,
	}
}
