package aws

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
)

func InitAwsConfig(k *koanf.Koanf) *aws.Config {
	log.Info().Str("aws connection", "Initiating").Msg("connecting to aws")
	awsConn, err := config.LoadDefaultConfig(context.TODO(), config.WithRegion(k.String("aws.regoin")))
	if err != nil {
		log.Fatal().Str("aws connection", "Failed").Msg("connection to aws failed")
	}
	log.Info().Str("aws connection", "Success").Msg("connected to aws successfully")
	return &awsConn
}
