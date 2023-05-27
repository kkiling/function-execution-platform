package logging

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type Config struct {
	SystemName string
	Level      zerolog.Level
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
}

func InitLogging(cfg *Config) error {
	if cfg.SystemName == "" {
		return fmt.Errorf("system is required property for logging")
	}
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	zerolog.SetGlobalLevel(cfg.Level)
	log.Logger = log.With().Str("system", cfg.SystemName).Logger()
	return nil
}
