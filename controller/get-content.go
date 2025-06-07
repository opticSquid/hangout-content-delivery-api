package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/contrib/bridges/otelslog"
	"go.opentelemetry.io/otel"
	"hangoutsb.in/hangout-content-delivery-api/model"
)

const name = "hangoutsb.in/hangout-content-delivery-api/controller"

var (
	tracer = otel.Tracer(name)
	meter  = otel.Meter(name)
	logger = otelslog.NewLogger(name)
)

func GetContent(w http.ResponseWriter, r *http.Request) {
	ctx, span := tracer.Start(r.Context(), "controller")
	defer span.End()
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
