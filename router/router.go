package router

import (
	"net/http"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/controller"
)

func StartServer(k *koanf.Koanf, cc *controller.ControllerConfig) {
	log.Info().Str("application name", k.String("application.name")).Str("status", "starting").Msg("starting application")

	http.HandleFunc("/"+k.String("application.name")+"/v1/get-content/{video_id}", withHeaders(cc.GetVideo, k))
	http.HandleFunc("/"+k.String("application.name")+"/v1/get-profile-photo/{image_id}", withHeaders(cc.GetImage, k))

	log.Info().Str("port", k.String("server.port")).Msg("starting http server")
	if err := http.ListenAndServe(":"+k.String("server.port"), nil); err != nil {
		log.Fatal().Msg("http server failed to start")
	}
}
