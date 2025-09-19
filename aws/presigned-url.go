package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

// for photo
func GeneratePreSignedUrl(cfg *aws.Config, k *koanf.Koanf, image string, log zerolog.Logger, ctx context.Context) string {
	tr := otel.Tracer("hangout.content-delivery-api.aws")
	ctx, span := tr.Start(ctx, "generate presigned url")
	defer span.End()
	span.SetAttributes(
		attribute.String("image", image),
		attribute.String("sign-type", "presign url"),
	)
	log = log.With().Ctx(ctx).Str("sign-type", "presign url").Logger()
	log.Debug().Msg("creating presign url client")

	client := s3.NewFromConfig(*cfg)
	presignClient := s3.NewPresignClient(client)
	bucket := k.String("aws.image.s3.bucket")
	expiration := k.Int("aws.image.s3.expiration-duration-seconds")
	log.Info().Msg("presigning file")
	req, err := presignClient.PresignGetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(image),
	}, s3.WithPresignExpires(time.Duration(expiration)*time.Second))
	if err != nil {
		log.Error().Err(err).Str("image", image).Msg("could not presign the image")
		return ""
	}
	log.Info().Msg("successfully presigned file")
	return req.URL
}
