package controller

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/rs/zerolog/log"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"hangoutsb.in/hangout-content-delivery-api/aws"
)

func (config *ControllerConfig) GetImage(w http.ResponseWriter, r *http.Request) {
	tr := otel.Tracer("hangout.content-delivery-api.controller")
	ctx := r.Context()
	ctx, span := tr.Start(ctx, "Get Profile Photo")
	defer span.End()
	span.SetAttributes(
		attribute.String("filename", r.PathValue("image_id")),
		attribute.String("pathname", r.Pattern),
		attribute.String("method", r.Method),
	)
	log := log.With().Ctx(ctx).Str("filename", r.PathValue("image_id")).Logger()
	log.Info().Msg("received a request for profile photo")
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

	fileName := r.PathValue("image_id")
	url := aws.GeneratePreSignedUrl(config.awsConfig, config.appConfig, fileName, log, ctx)

	if url == "" {
		err := fmt.Errorf("presigned url could not be generated")

		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

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
		"file":    fileName,
		"url":     url,
		"message": "Signed URL genereated successfully",
	}

	w.WriteHeader(http.StatusOK)

	// Encode response as JSON
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())

		log.Error().Err(err).Msg("failed to encode json response")
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
	}
	log.Info().Str("Path", r.Pattern).Str("Method", r.Method).Msg("response sent")
}
