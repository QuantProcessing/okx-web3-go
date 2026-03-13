package okx

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestNewClient(t *testing.T) {
	client := NewClient()

	if client == nil {
		t.Fatal("NewClient() returned nil")
	}

	if client.BaseURL != BaseURL {
		t.Errorf("Expected BaseURL %s, got %s", BaseURL, client.BaseURL)
	}

	if client.HTTPClient == nil {
		t.Error("HTTPClient is nil")
	}

	if client.Logger == nil {
		t.Error("Logger is nil")
	}
}

func TestWithCredentials(t *testing.T) {
	client := NewClient()

	result := client.WithCredentials("key", "secret", "pass", "proj")

	if result != client {
		t.Error("WithCredentials should return the same client for chaining")
	}

	if client.ApiKey != "key" {
		t.Errorf("Expected ApiKey 'key', got '%s'", client.ApiKey)
	}
	if client.SecretKey != "secret" {
		t.Errorf("Expected SecretKey 'secret', got '%s'", client.SecretKey)
	}
	if client.Passphrase != "pass" {
		t.Errorf("Expected Passphrase 'pass', got '%s'", client.Passphrase)
	}
	if client.ProjectID != "proj" {
		t.Errorf("Expected ProjectID 'proj', got '%s'", client.ProjectID)
	}
	if client.Signer == nil {
		t.Error("Signer should be set after WithCredentials")
	}
}

func TestDo_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}
		if r.Header.Get("Content-Type") != "application/json" {
			t.Error("Content-Type header not set")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"0","data":[]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL

	data, err := client.Do(context.Background(), "GET", "/test", nil, false)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}

	expected := `{"code":"0","data":[]}`
	if string(data) != expected {
		t.Errorf("Expected '%s', got '%s'", expected, string(data))
	}
}

func TestDo_WithQueryParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Query().Get("chainIndex") != "1" {
			t.Errorf("Expected chainIndex=1, got %s", r.URL.Query().Get("chainIndex"))
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"0","data":[]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL

	params := url.Values{}
	params.Set("chainIndex", "1")

	_, err := client.Do(context.Background(), "GET", "/test", params, false)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
}

func TestDo_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("internal server error"))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL

	_, err := client.Do(context.Background(), "GET", "/test", nil, false)
	if err == nil {
		t.Fatal("Expected error for HTTP 500")
	}

	expectedMsg := "http error 500: internal server error"
	if err.Error() != expectedMsg {
		t.Errorf("Expected error '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestDo_AuthWithoutCredentials(t *testing.T) {
	client := NewClient()

	_, err := client.Do(context.Background(), "GET", "/test", nil, true)
	if err == nil {
		t.Fatal("Expected error when auth=true but no credentials")
	}

	if err.Error() != "credentials required for authenticated request" {
		t.Errorf("Unexpected error: %s", err.Error())
	}
}

func TestDo_AuthSetsHeaders(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify auth headers are set
		if r.Header.Get("OK-ACCESS-KEY") != "test-key" {
			t.Errorf("Expected OK-ACCESS-KEY 'test-key', got '%s'", r.Header.Get("OK-ACCESS-KEY"))
		}
		if r.Header.Get("OK-ACCESS-PASSPHRASE") != "test-pass" {
			t.Errorf("Expected OK-ACCESS-PASSPHRASE 'test-pass'")
		}
		if r.Header.Get("OK-ACCESS-SIGN") == "" {
			t.Error("OK-ACCESS-SIGN should not be empty")
		}
		if r.Header.Get("OK-ACCESS-TIMESTAMP") == "" {
			t.Error("OK-ACCESS-TIMESTAMP should not be empty")
		}
		if r.Header.Get("OK-ACCESS-PROJECT") != "test-proj" {
			t.Errorf("Expected OK-ACCESS-PROJECT 'test-proj'")
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"0","data":[]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("test-key", "test-secret", "test-pass", "test-proj")

	_, err := client.Do(context.Background(), "GET", "/test", nil, true)
	if err != nil {
		t.Fatalf("Do returned error: %v", err)
	}
}

func TestDoPost_Success(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != "POST" {
			t.Errorf("Expected POST, got %s", r.Method)
		}
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`{"code":"0","data":[{"result":"ok"}]}`))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL

	data, err := client.DoPost(context.Background(), "/test", map[string]string{"key": "value"}, false)
	if err != nil {
		t.Fatalf("DoPost returned error: %v", err)
	}

	if string(data) == "" {
		t.Error("DoPost returned empty data")
	}
}

func TestDoPost_AuthWithoutCredentials(t *testing.T) {
	client := NewClient()

	_, err := client.DoPost(context.Background(), "/test", nil, true)
	if err == nil {
		t.Fatal("Expected error when auth=true but no credentials")
	}
}

func TestDoPost_HTTPError(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL

	_, err := client.DoPost(context.Background(), "/test", nil, false)
	if err == nil {
		t.Fatal("Expected error for HTTP 400")
	}
}

func TestDo_ContextCancelled(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte(`ok`))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL

	ctx, cancel := context.WithCancel(context.Background())
	cancel() // Cancel immediately

	_, err := client.Do(ctx, "GET", "/test", nil, false)
	if err == nil {
		t.Fatal("Expected error for cancelled context")
	}
}

// --- Tests for generic request[T] helper ---

// newTestServerWithResponse creates a test server that returns the given BaseResponse JSON.
func newTestServerWithResponse(t *testing.T, code string, msg string, data interface{}) *httptest.Server {
	t.Helper()
	return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := map[string]interface{}{
			"code": code,
			"msg":  msg,
			"data": data,
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(resp)
	}))
}

func TestRequest_Success(t *testing.T) {
	server := newTestServerWithResponse(t, "0", "", []map[string]string{
		{"chainId": "1", "chainName": "Ethereum"},
	})
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	results, err := request[SupportedChain](client, context.Background(), "GET", "/test", nil, true)
	if err != nil {
		t.Fatalf("request returned error: %v", err)
	}
	if len(results) != 1 {
		t.Fatalf("Expected 1 result, got %d", len(results))
	}
	if results[0].ChainName != "Ethereum" {
		t.Errorf("Expected ChainName 'Ethereum', got '%s'", results[0].ChainName)
	}
}

func TestRequest_APIError(t *testing.T) {
	server := newTestServerWithResponse(t, "51003", "Insufficient liquidity", []interface{}{})
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	_, err := request[QuoteResponse](client, context.Background(), "GET", "/test", nil, true)
	if err == nil {
		t.Fatal("Expected error for non-zero API code")
	}

	if !IsDexError(err) {
		t.Errorf("Expected DexError, got %T", err)
	}

	if GetErrorCode(err) != "51003" {
		t.Errorf("Expected error code '51003', got '%s'", GetErrorCode(err))
	}
}

func TestRequest_UnsupportedMethod(t *testing.T) {
	client := NewClient()
	client.WithCredentials("key", "secret", "pass", "proj")

	_, err := request[SupportedChain](client, context.Background(), "DELETE", "/test", nil, false)
	if err == nil {
		t.Fatal("Expected error for unsupported method")
	}

	expectedMsg := "unsupported HTTP method: DELETE"
	if err.Error() != expectedMsg {
		t.Errorf("Expected '%s', got '%s'", expectedMsg, err.Error())
	}
}

func TestRequest_InvalidJSON(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("not json"))
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL

	_, err := request[SupportedChain](client, context.Background(), "GET", "/test", nil, false)
	if err == nil {
		t.Fatal("Expected error for invalid JSON")
	}
}

// --- Tests for DexError ---

func TestDexError_Error(t *testing.T) {
	err := &DexError{
		Code:    "51003",
		Message: "Insufficient liquidity",
	}

	expected := "OKX DEX API error [51003]: Insufficient liquidity"
	if err.Error() != expected {
		t.Errorf("Expected '%s', got '%s'", expected, err.Error())
	}
}

func TestIsDexError(t *testing.T) {
	dexErr := &DexError{Code: "51003", Message: "test"}
	stdErr := fmt.Errorf("standard error")

	if !IsDexError(dexErr) {
		t.Error("IsDexError should return true for DexError")
	}

	if IsDexError(stdErr) {
		t.Error("IsDexError should return false for standard error")
	}

	if IsDexError(nil) {
		t.Error("IsDexError should return false for nil")
	}
}

func TestGetErrorCode(t *testing.T) {
	dexErr := &DexError{Code: "51003", Message: "test"}
	stdErr := fmt.Errorf("standard error")

	if code := GetErrorCode(dexErr); code != "51003" {
		t.Errorf("Expected '51003', got '%s'", code)
	}

	if code := GetErrorCode(stdErr); code != "" {
		t.Errorf("Expected empty string, got '%s'", code)
	}

	if code := GetErrorCode(nil); code != "" {
		t.Errorf("Expected empty string for nil, got '%s'", code)
	}
}
