package storage

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

var sess *session.Session

func BlobStorageConnInit(config *koanf.Koanf) {
	s, err := session.NewSession(&aws.Config{
		Region:      aws.String(config.String("storage.region")),
		Credentials: credentials.NewStaticCredentials(config.String("storage.accessKey"), config.String("storage.secretKey"), ""),
	},
	)

	if err != nil {
		log.Fatal().Msg("could not login to aws")
	}
	sess = s
	log.Info().Msg("logged in to aws")
}

func GetPreSignedUrl(filenameWithExtension string, config *koanf.Koanf) string {

	// Create S3 service client
	svc := s3.New(sess)

	req, _ := svc.GetObjectRequest(&s3.GetObjectInput{
		Bucket: aws.String(config.String("storage.bucketName")),
		Key:    aws.String(filenameWithExtension),
	})
	urlStr, err := req.Presign(5 * time.Minute)

	if err != nil {
		log.Fatal().Str("error", err.Error()).Msg("Failed to Sign request")
	}

	log.Debug().Str("PreSigned URL", urlStr).Msg("Presigned url generated")
	return urlStr
}
