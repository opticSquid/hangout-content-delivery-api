package controller

import (
	"encoding/json"
	"net/http"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/model"
	"hangoutsb.in/hangout-content-delivery-api/storage"
)

var config *koanf.Koanf

func SetConfig(c *koanf.Koanf) {
	config = c
}
func GetContent(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("recieved request")
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Error().Str("Path", r.Pattern).Str("Method", r.Method).Msg("method not allowed in path")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Create a response object
	resp := model.Response{
		Url:       storage.GetPreSignedUrl(r.PathValue("content_id"), config),
		ContentId: r.PathValue("content_id"),
	}
	// Encode response to JSON and write to client
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Str("Path", r.Pattern).Str("Method", r.Method).Str("reason", err.Error()).Msg("error encoding response to JSON")
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("response sent")
}
