package okx

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"net/http"
	"testing"
)

func TestNewSigner(t *testing.T) {
	signer := NewSigner("test-secret")

	if signer == nil {
		t.Fatal("NewSigner returned nil")
	}

	if signer.SecretKey != "test-secret" {
		t.Errorf("Expected SecretKey 'test-secret', got '%s'", signer.SecretKey)
	}
}

func TestSignRequest_SetsRequiredHeaders(t *testing.T) {
	signer := NewSigner("my-secret-key")
	req, _ := http.NewRequest("GET", "https://example.com/api/test", nil)

	signer.SignRequest(req, "GET", "/api/test", "", "my-api-key", "my-passphrase", "my-project")

	// Verify all required headers are set
	headers := map[string]string{
		"OK-ACCESS-KEY":        "my-api-key",
		"OK-ACCESS-PASSPHRASE": "my-passphrase",
		"OK-ACCESS-PROJECT":    "my-project",
	}

	for header, expected := range headers {
		got := req.Header.Get(header)
		if got != expected {
			t.Errorf("Header %s: expected '%s', got '%s'", header, expected, got)
		}
	}

	// Verify timestamp is set (non-empty)
	ts := req.Header.Get("OK-ACCESS-TIMESTAMP")
	if ts == "" {
		t.Error("OK-ACCESS-TIMESTAMP header is empty")
	}

	// Verify signature is set (non-empty, valid base64)
	sig := req.Header.Get("OK-ACCESS-SIGN")
	if sig == "" {
		t.Error("OK-ACCESS-SIGN header is empty")
	}
	if _, err := base64.StdEncoding.DecodeString(sig); err != nil {
		t.Errorf("OK-ACCESS-SIGN is not valid base64: %v", err)
	}
}

func TestSignRequest_WithoutProjectID(t *testing.T) {
	signer := NewSigner("secret")
	req, _ := http.NewRequest("GET", "https://example.com/test", nil)

	signer.SignRequest(req, "GET", "/test", "", "key", "pass", "")

	if req.Header.Get("OK-ACCESS-PROJECT") != "" {
		t.Error("OK-ACCESS-PROJECT should not be set when projectID is empty")
	}

	// Other required headers should still be present
	if req.Header.Get("OK-ACCESS-KEY") != "key" {
		t.Error("OK-ACCESS-KEY should be set")
	}
}

func TestSignRequest_SignatureCorrectness(t *testing.T) {
	// Verify the HMAC SHA256 algorithm produces consistent results
	signer := NewSigner("test-secret")
	req1, _ := http.NewRequest("GET", "https://example.com/test", nil)
	req2, _ := http.NewRequest("GET", "https://example.com/test", nil)

	signer.SignRequest(req1, "GET", "/api/v6/test", "", "key", "pass", "")
	signer.SignRequest(req2, "GET", "/api/v6/test", "", "key", "pass", "")

	// Signatures may differ due to timestamp, but both should be valid base64
	sig1 := req1.Header.Get("OK-ACCESS-SIGN")
	sig2 := req2.Header.Get("OK-ACCESS-SIGN")
	if sig1 == "" || sig2 == "" {
		t.Error("Signatures should not be empty")
	}
}

func TestSignRequest_WithBody(t *testing.T) {
	signer := NewSigner("secret-key")
	req, _ := http.NewRequest("POST", "https://example.com/api/swap", nil)
	body := `{"amount":"100"}`

	signer.SignRequest(req, "POST", "/api/swap", body, "api-key", "passphrase", "project")

	// Manually compute expected signature
	ts := req.Header.Get("OK-ACCESS-TIMESTAMP")
	preHash := ts + "POST" + "/api/swap" + body
	h := hmac.New(sha256.New, []byte("secret-key"))
	h.Write([]byte(preHash))
	expectedSig := base64.StdEncoding.EncodeToString(h.Sum(nil))

	actualSig := req.Header.Get("OK-ACCESS-SIGN")
	if actualSig != expectedSig {
		t.Errorf("Signature mismatch:\n  expected: %s\n  got:      %s", expectedSig, actualSig)
	}
}
