package aws

import (
	"context"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

// for photo
func GeneratePreSignedUrl(cfg *aws.Config, k *koanf.Koanf, image string) string {
	client := s3.NewFromConfig(*cfg)
	presignClient := s3.NewPresignClient(client)
	bucket := k.String("aws.image.s3.bucket")
	expiration := k.Int("aws.image.s3.expirationDurationInSeconds")
	req, err := presignClient.PresignGetObject(context.TODO(), &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(image),
	}, s3.WithPresignExpires(time.Duration(expiration)*time.Second))
	if err != nil {
		log.Error().Err(err).Str("image", image).Msg("could not presign the image")
		return ""
	}
	return req.URL
}
