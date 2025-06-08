package storage

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/knadh/koanf/v2"
	"github.com/rs/zerolog/log"
	"hangoutsb.in/hangout-content-delivery-api/model"
)

func GeneratePreSignedCookies(filenameWithExtension string, config *koanf.Koanf) (
	expiresCookie, signatureCookie, keyPairIDCookie, policyCookie *http.Cookie, err error) {
	// 1. Read and parse the private key
	privateKeyBytes, err := os.ReadFile(config.String("cloudfront.privateKeyPath"))
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to read private key: %w", err)
	}

	privateKey, err := parsePrivateKey(privateKeyBytes)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// 2. Define the expiration time (Unix epoch time)
	expiresAt := time.Now().Add(time.Duration(config.Int("cloudfront.expirationDurationInSeconds")) * time.Second)
	expiresEpoch := expiresAt.Unix()

	// 3. Construct the Custom Policy JSON
	policy := model.Policy{
		Statement: []model.PolicyStatement{
			{
				Resource: fmt.Sprintf("http*://%s%s", config.String("cloudfront.domain"), getBaseFileNamePhoto(filenameWithExtension)+"/"+filenameWithExtension),
				Condition: model.PolicyCondition{
					DateLessThan: map[string]int64{
						"AWS:EpochTime": expiresEpoch,
					},
				},
			},
		},
	}

	policyBytes, err := json.Marshal(policy)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to marshal policy: %w", err)
	}
	policyString := string(policyBytes)
	log.Printf("Generated Policy: %s", policyString)

	// 4. Sign the policy
	h := sha1.New()
	h.Write(policyBytes)
	hashedPolicy := h.Sum(nil)

	signature, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, hashedPolicy)
	if err != nil {
		return nil, nil, nil, nil, fmt.Errorf("failed to sign policy: %w", err)
	}

	// 5. Base64 URL-safe encode the signature and the policy
	encodedSignature := base64.URLEncoding.EncodeToString(signature)
	encodedPolicy := base64.URLEncoding.EncodeToString(policyBytes)

	// 6. Create the http.Cookie objects
	expiresCookie = &http.Cookie{
		Name:     "CloudFront-Expires",
		Value:    fmt.Sprintf("%d", expiresEpoch),
		Path:     "/",
		Domain:   "." + config.String("cloudfront.domain"), // Leading dot for subdomains
		Expires:  expiresAt,
		Secure:   true,                 // Always use secure cookies for HTTPS
		HttpOnly: true,                 // Recommended to prevent client-side JavaScript access
		SameSite: http.SameSiteLaxMode, // Adjust as per your CORS policy
	}

	signatureCookie = &http.Cookie{
		Name:     "CloudFront-Signature",
		Value:    encodedSignature,
		Path:     "/",
		Domain:   "." + config.String("cloudfront.domain"),
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	keyPairIDCookie = &http.Cookie{
		Name:     "CloudFront-Key-Pair-Id",
		Value:    config.String("cloudfront.publicKeyID"),
		Path:     "/",
		Domain:   "." + config.String("cloudfront.domain"),
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
	}

	policyCookie = &http.Cookie{
		Name:     "CloudFront-Policy", // IMPORTANT: This name is fixed
		Value:    encodedPolicy,       // Use the encodedPolicy variable here
		Path:     "/",
		Domain:   "." + config.String("cloudfront.domain"),
		Expires:  expiresAt,
		Secure:   true,
		HttpOnly: true, // You can make this false if your frontend needs to read it, but generally keep HttpOnly.
		SameSite: http.SameSiteLaxMode,
	}

	return expiresCookie, signatureCookie, keyPairIDCookie, policyCookie, nil
}

// get base file name from filename with extension
func getBaseFileNamePhoto(filenameWithExtension string) string {
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

// parsePrivateKey parses a PEM-encoded RSA private key.
func parsePrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
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
