package aws

import (
	"context"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/feature/cloudfront/sign"
	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
)

func GeneratePreSignedCookies(dirName string, k *koanf.Koanf, log zerolog.Logger, ctx context.Context) ([]*http.Cookie, error) {
	tr := otel.Tracer("hangout.content-delivery-api.aws")
	ctx, span := tr.Start(ctx, "generate presigned cookies")
	defer span.End()
	span.SetAttributes(
		attribute.String("dirname", dirName),
		attribute.String("sign-type", "presign cookie"),
	)
	log = log.With().Ctx(ctx).Str("sign-type", "presign cookie").Logger()
	privateKeyPath := k.String("aws.video.cloudfront.private-key-path")
	privateKeyBytes, err := os.ReadFile(privateKeyPath)
	if err != nil {
		log.Error().Err(err).Msg("failed to read private key")
		return nil, fmt.Errorf("failed to read private key: %w", err)
	}
	log.Debug().Msg("successfully read private key")

	privateKey, err := parsePrivateKey(privateKeyBytes, log, ctx)
	if err != nil {
		log.Error().Err(err).Msg("failed to parse private key")
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}
	log.Debug().Msg("successfully parsed private key")

	publicKeyId := k.String("aws.video.cloudfront.public-key-id")
	expiresAt := time.Now().Add(time.Duration(k.Int("aws.video.cloudfront.expirationDurationInSeconds")) * time.Second)
	resource := fmt.Sprintf("http*://%s/%s/*", k.String("aws.video.cloudfront.domain"), dirName)

	log.Debug().Msg("starting to generate cookies")
	cookies, err := sign.NewCookieSigner(publicKeyId, privateKey).Sign(resource, expiresAt)
	if err != nil {
		return nil, fmt.Errorf("failed to generate cookies: %w", err)
	}
	log.Debug().Msg("successfully generated cookies")
	for _, c := range cookies {
		customizeCookie(c, dirName, expiresAt, k)
	}
	return cookies, nil
}

// parsePrivateKey parses a PEM-encoded RSA private key.
func parsePrivateKey(pemBytes []byte, log zerolog.Logger, ctx context.Context) (*rsa.PrivateKey, error) {
	log = log.With().Ctx(ctx).Logger()
	block, _ := pem.Decode(pemBytes)
	if block == nil || block.Type != "PRIVATE KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing RSA private key")
	}

	log.Debug().Msg("Attempting to parse private key as PKCS#8")
	privateKeyInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse PKCS#8 private key: %w", err)
	}

	rsaKey, ok := privateKeyInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("parsed PKCS#8 private key is not an RSA private key (got %T)", privateKeyInterface)
	}
	return rsaKey, nil
}

func customizeCookie(c *http.Cookie, dirName string, expiresAt time.Time, k *koanf.Koanf) *http.Cookie {
	c.Domain = "." + k.String("cookie.domain")
	c.Path = "/" + dirName
	c.Expires = expiresAt
	return c
}
