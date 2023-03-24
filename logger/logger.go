package logger

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"os"
)

type Configurer interface {
	LogLevel() string
	LogFormat() string
}

func Init(config Configurer) {

	switch config.LogLevel() {
	case "info":
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	case "debug":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	default:
		panic(fmt.Sprintf("Unknown LogLevel: %s\n", config.LogLevel()))
	}

	switch config.LogFormat() {
	case "console":
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	case "json":
		log.Logger = log.With().Caller().Logger()
	default:
		panic(fmt.Sprintf("Unknown LogFormat: %s\n", config.LogFormat()))

	}

	log.Debug().Msg("Logger initialized")
}
