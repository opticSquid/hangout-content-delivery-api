package storage

import (
	"context"
	"net/url"
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/rs/zerolog/log"
)

var blobClient *minio.Client

func BlobStorageConnInit(config *koanf.Koanf) {
	log.Debug().Str("storage url", config.String("storage.endpoint")).Msg("connecting to blob storage")
	s3Client, err := minio.New(config.String("storage.endpoint"), &minio.Options{
		Creds:  credentials.NewStaticV4(config.String("storage.accessKey"), config.String("storage.secretKey"), ""),
		Secure: config.Bool("storage.isSecure"),
	})
	if err != nil {
		log.Fatal().Str("reasson", err.Error()).Msg("could not connect to blob storage")
	}
	blobClient = s3Client
}

func GetPreSignedUrl(filenameWithExtension string, config *koanf.Koanf) *url.URL {
	log.Debug().Str("file", filenameWithExtension).Msg("fetching presigned url")
	// Extract base file name
	baseFileName := getBaseFileName(filenameWithExtension)

	// Construct the correct object path: <baseFileName>/<filenameWithExtension>
	objectKey := baseFileName + "/" + filenameWithExtension

	reqParams := make(url.Values)
	reqParams.Set("response-content-disposition", "attachment; filename="+filenameWithExtension)

	presignedURL, err := blobClient.PresignedGetObject(context.Background(), config.String("storage.bucketName"), objectKey, time.Duration(300)*time.Second, reqParams)
	if err != nil {
		log.Fatal().Str("file", filenameWithExtension).Str("reasson", err.Error()).Msg("could not fetch presigned url of file")
	}

	return presignedURL
}

func getBaseFileName(filenameWithExtension string) string {
	// Expected format: <base_name>_<resolution>.<ext>
	lastDotIndex := strings.LastIndex(filenameWithExtension, "_")
	if lastDotIndex != -1 && lastDotIndex < len(filenameWithExtension)-1 {
		// Extract the substring before the last dot
		log.Debug().Str("base-file-name", filenameWithExtension[:lastDotIndex])
		return filenameWithExtension[:lastDotIndex]
	}
	log.Debug().Str("base-file-name", filenameWithExtension[:lastDotIndex])
	return filenameWithExtension
}
