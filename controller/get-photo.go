package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/aws"
)

func (config *ControllerConfig) GetPhoto(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("recieved request")
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Error().Str("Path", r.Pattern).Str("Method", r.Method).Msg("method not allowed in path")
		writeProblemDetails(w, http.StatusMethodNotAllowed, "Method not allowed",
			"Only GET requests are supported on this endpoint.",
			"https://httpstatuses.com/405",
			r.URL.Path)
		return
	}

	fileName := r.PathValue("photo_id")
	url := aws.GeneratePreSignedUrl(config.awsConfig, config.appConfig, fileName)

	if url == "" {
		writeProblemDetails(w,
			http.StatusInternalServerError,
			"Failed to generate presigned URL",
			"Could not generate URL for the given image, make sure the image is valid",
			"https://httpstatuses.com/500",
			r.URL.Path)
		return
	}
	// Set content type to JSON
	w.Header().Set("Content-Type", "application/json")

	// Prepare response payload
	resp := map[string]string{
		"photoId": fileName,
		"url":     url,
	}

	w.WriteHeader(http.StatusOK)

	// Encode response as JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg("failed to encode json response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("response sent")
}
