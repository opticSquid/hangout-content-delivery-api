package controller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/aws"
)

func (config *ControllerConfig) GetVideo(w http.ResponseWriter, r *http.Request) {
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Str("file", r.PathValue("video_id")).Msg("recieved request")
	// Only allow GET requests
	if r.Method != http.MethodGet {
		log.Error().Str("Path", r.Pattern).Str("Method", r.Method).Msg("method not allowed in path")
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	fileName := r.PathValue("video_id")
	dirName := getDirName(fileName)
	preSignedCookies, err := aws.GeneratePreSignedCookies(dirName, config.appConfig)
	if err != nil {
		log.Error().Err(err).Str("video file", fileName).Msg("could not generate cookies for the given file")
		writeProblemDetails(w, http.StatusInternalServerError, "Failed to generate Presigned cookies", "Could not generate presigned cookies for the given file. Make sure the filename is valid", "https://httpstatuses.com/500", r.URL.Path)
		return
	}
	for _, c := range preSignedCookies {
		http.SetCookie(w, c)
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	// Prepare response payload
	resp := map[string]string{
		"file":    fileName,
		"message": "Signed cookies set successfully",
	}
	// Encode response as JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Error().Err(err).Msg("failed to encode json response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("response sent")

}

func getDirName(filename string) string {
	fName, _, found := strings.Cut(filename, ".")
	if found {
		return fName
	} else {
		return filename
	}
}
