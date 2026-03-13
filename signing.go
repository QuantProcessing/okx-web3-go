package okx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"time"
)

// Signer handles OKX DEX API signature generation.
// Uses the same HMAC SHA256 algorithm as the V5 trading API.
type Signer struct {
	SecretKey string
}

// NewSigner creates a new Signer with the given secret key.
func NewSigner(secretKey string) *Signer {
	return &Signer{SecretKey: secretKey}
}

// SignRequest adds necessary authentication headers to the request.
// OKX DEX API requires:
// - OK-ACCESS-KEY: API key
// - OK-ACCESS-SIGN: HMAC SHA256 signature
// - OK-ACCESS-TIMESTAMP: ISO8601 timestamp
// - OK-ACCESS-PASSPHRASE: API key passphrase
// - OK-ACCESS-PROJECT: Project ID (optional, DEX-specific)
func (s *Signer) SignRequest(req *http.Request, method, path, body, apiKey, passphrase, projectID string) {
	// 1. Timestamp in ISO8601 format
	ts := time.Now().UTC().Format("2006-01-02T15:04:05.999Z07:00")

	// 2. PreHash String: timestamp + method + requestPath + body
	preHash := ts + method + path + body

	// 3. HMAC SHA256 signature
	h := hmac.New(sha256.New, []byte(s.SecretKey))
	h.Write([]byte(preHash))
	signature := base64.StdEncoding.EncodeToString(h.Sum(nil))

	// 4. Set required headers
	req.Header.Set("OK-ACCESS-KEY", apiKey)
	req.Header.Set("OK-ACCESS-SIGN", signature)
	req.Header.Set("OK-ACCESS-TIMESTAMP", ts)
	req.Header.Set("OK-ACCESS-PASSPHRASE", passphrase)

	// 5. Set optional project ID header (DEX-specific)
	if projectID != "" {
		req.Header.Set("OK-ACCESS-PROJECT", projectID)
	}
}
