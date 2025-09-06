package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"hangoutsb.in/hangout-content-delivery-api/aws"
)

func (config *ControllerConfig) GetVideo(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("hangout.content-delivery-api.controller")
	ctx := r.Context()
	ctx, span := tr.Start(ctx, "Get Content")
	defer span.End()
	span.SetAttributes(
		attribute.String("filename", r.PathValue("video_id")),
		attribute.String("pathname", r.Pattern),
		attribute.String("method", r.Method),
	)
	log := log.With().Ctx(ctx).Str("filename", r.PathValue("video_id")).Logger()
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Str("file", r.PathValue("video_id")).Msg("recieved request")
	// Only allow GET requests
	if r.Method != http.MethodGet {
		err := fmt.Errorf("method not allowed")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		log.Error().Err(err).Str("Path", r.Pattern).Str("Method", r.Method).Msg("method not allowed in path")
		writeProblemDetails(w, http.StatusMethodNotAllowed, "Method not allowed",
			"Only GET requests are supported on this endpoint.",
			"https://httpstatuses.com/405",
			r.URL.Path)

		return
	}
	fileName := r.PathValue("video_id")
	dirName := getDirName(fileName)
	preSignedCookies, err := aws.GeneratePreSignedCookies(dirName, config.appConfig, log, ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

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
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

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
