package logger

import (
	"os"
	"strings"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/bridges/otelzerolog"
)

func InitLogger(config *koanf.Koanf) {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	switch config.String("logging.level") {
	case "TRACE":
		zerolog.SetGlobalLevel(zerolog.TraceLevel)
	case "DEBUG":
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	case "WARN":
		zerolog.SetGlobalLevel(zerolog.WarnLevel)
	default:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
	}

	hook := otelzerolog.NewHook(config.String("application.name"))
	baseLogger := zerolog.New(os.Stdout).With().Timestamp().Str("app-name", config.String("application.name")).Logger()
	rootLogger := baseLogger.Hook(hook)
	// adding root context to rootLogger
	// Override the global logger to use my logger
	log.Logger = rootLogger

	log.Info().Str("global logging level", strings.ToUpper(zerolog.GlobalLevel().String())).Msg("loging level set")
}
