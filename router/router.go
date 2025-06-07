package router

import (
	"net/http"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/controller"
)

func StartServer(config *koanf.Koanf) {
	http.HandleFunc("/get-content", controller.GetContent)
	log.Info().Str("port", config.String("server.port")).Msg("starting http server")
	if err := http.ListenAndServe(":"+config.String("server.port"), nil); err != nil {
		log.Fatal().Msg("http server failed to start")
	}
}
