package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/model"
)

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
		Message: "Hello, JSON!",
		Status:  "success",
	}
	// Encode response to JSON and write to client
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
}
