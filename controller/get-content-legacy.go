package controller

import (
	"encoding/json"
	"net/http"

	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/model"
	"hangoutsb.in/hangout-content-delivery-api/storage"
)

func (config *ControllerConfig) GetContentLegacy(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("recieved request")
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Error().Str("Path", r.Pattern).Str("Method", r.Method).Msg("method not allowed in path")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	fileName := r.PathValue("content_id")
	// Set cookies in response Object
	expiresCookie, signatureCookie, keyPairIDCookie, policyCookie, err := storage.GeneratePreSignedCookies(fileName, config.appConfig)

	var resp model.Response

	if err != nil {
		log.Error().Str("error", err.Error()).Msg("failed to generate presigned cookies")
		resp = model.Response{
			ContentId: fileName,
		}
		w.WriteHeader(http.StatusInternalServerError)
	} else {
		// --- 3. Set the Cookies in the HTTP Response ---
		// These cookies will be automatically sent by the browser
		// with subsequent requests to your CloudFront distribution.
		http.SetCookie(w, expiresCookie)
		http.SetCookie(w, signatureCookie)
		http.SetCookie(w, keyPairIDCookie)
		http.SetCookie(w, policyCookie) // <-- The crucial fourth cookie!

		resp = model.Response{
			ContentId: fileName,
		}
		w.WriteHeader(http.StatusOK)

	}
	// Set the content type to JSON
	w.Header().Set("Content-Type", "application/json")
	// Encode response to JSON and write to client
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Str("Path", r.Pattern).Str("Method", r.Method).Str("reason", err.Error()).Msg("error encoding response to JSON")
		http.Error(w, "Failed to encode JSON", http.StatusInternalServerError)
	}
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("response sent")
}
