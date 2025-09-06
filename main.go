package main

import (
	"context"
	"errors"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/aws"
	"hangoutsb.in/hangout-content-delivery-api/config"
	"hangoutsb.in/hangout-content-delivery-api/controller"
	"hangoutsb.in/hangout-content-delivery-api/logger"
	"hangoutsb.in/hangout-content-delivery-api/router"
	"hangoutsb.in/hangout-content-delivery-api/telemetry"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var CONFIG = koanf.New(".")

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	config.InitAppConfig(CONFIG)
	logger.InitLogger(CONFIG)

	otelShutdown, err := telemetry.SetUpOTelSDK(ctx, CONFIG)
	if err != nil {
		log.Error().Err(err).Msg("could not set up telemetry")
		return
	}
	// Handle shutdown properly so nothing leaks.
	defer func() {
		err = errors.Join(err, otelShutdown(context.Background()))
	}()
	telemetry.StartProcessMetrics()
	log.Info().Str("endpoint", CONFIG.String("otel.endpoint")).Msg("starting to send traces, logs, metrics")
	awsConn := aws.InitAwsConfig(CONFIG)
	controllerConfig := controller.InitControllerConfig(CONFIG, awsConn)
	router.StartServer(CONFIG, controllerConfig)
}
