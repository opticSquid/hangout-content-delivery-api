package logger

import (
	"strings"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
)

func InitLogger() {
	// UNIX Time is faster and smaller than most timestamps
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	log.Info().Msg("configured zerolog for logging")
	zerolog.SetGlobalLevel(zerolog.TraceLevel)
}

func SetGlobalLogLevel(config *koanf.Koanf) {
	log.Info().Str("requested global logging level", config.String("logging.level")).Msg("")
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
	log.Info().Str("configured global logging level", strings.ToUpper(zerolog.GlobalLevel().String())).Msg("")
}
