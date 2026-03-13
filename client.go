package okx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"

	"go.uber.org/zap"
)

const (
	BaseURL = "https://web3.okx.com"
)

// Client is the HTTP client for OKX DEX API.
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
	Signer     *Signer
	Logger     *zap.SugaredLogger

	// Credentials
	ApiKey     string
	SecretKey  string
	Passphrase string
	ProjectID  string
}

// NewClient creates a new OKX DEX API client.
// By default uses a no-op logger. Set client.Logger to inject your own.
func NewClient() *Client {
	return &Client{
		BaseURL:    BaseURL,
		HTTPClient: &http.Client{},
		Logger:     zap.NewNop().Sugar(),
	}
}

// WithCredentials sets the API credentials for authenticated requests.
// Returns the client for method chaining.
func (c *Client) WithCredentials(apiKey, secretKey, passphrase, projectID string) *Client {
	c.ApiKey = apiKey
	c.SecretKey = secretKey
	c.Passphrase = passphrase
	c.ProjectID = projectID
	c.Signer = NewSigner(secretKey)
	return c
}

// Do executes an HTTP request and returns the raw response body.
// It handles authentication signatures if auth is required.
//
// Parameters:
//   - ctx: Context for request timeout/cancellation
//   - method: HTTP method (GET, POST, etc.)
//   - path: API endpoint path (e.g., "/api/v6/dex/aggregator/supported/chain")
//   - params: Query parameters (nil for POST with body)
//   - auth: Whether authentication is required
//
// Returns raw response body or error.
func (c *Client) Do(ctx context.Context, method, path string, params url.Values, auth bool) ([]byte, error) {
	// Build full URL
	fullURL := c.BaseURL + path
	if len(params) > 0 {
		fullURL += "?" + params.Encode()
	}

	// Create request
	req, err := http.NewRequestWithContext(ctx, method, fullURL, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication if required
	if auth {
		if c.Signer == nil {
			return nil, fmt.Errorf("credentials required for authenticated request")
		}

		// For GET requests with query params, include them in signature
		pathForSign := path
		if len(params) > 0 {
			pathForSign += "?" + params.Encode()
		}

		c.Signer.SignRequest(req, method, pathForSign, "", c.ApiKey, c.Passphrase, c.ProjectID)
	}

	// Log request details in DEBUG mode
	if os.Getenv("DEBUG") != "" {
		c.Logger.Debugf("DEX API Request: %s %s", method, fullURL)
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// DoPost executes a POST request with JSON body.
// Used for endpoints that require POST with request body.
func (c *Client) DoPost(ctx context.Context, path string, payload interface{}, auth bool) ([]byte, error) {
	// Marshal payload to JSON
	var bodyReader io.Reader
	var bodyString string

	if payload != nil {
		jsonBytes, err := json.Marshal(payload)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal payload: %w", err)
		}
		bodyReader = bytes.NewBuffer(jsonBytes)
		bodyString = string(jsonBytes)
	}

	// Build full URL
	fullURL := c.BaseURL + path

	// Create request
	req, err := http.NewRequestWithContext(ctx, "POST", fullURL, bodyReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	// Set headers
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	// Add authentication if required
	if auth {
		if c.Signer == nil {
			return nil, fmt.Errorf("credentials required for authenticated request")
		}
		c.Signer.SignRequest(req, "POST", path, bodyString, c.ApiKey, c.Passphrase, c.ProjectID)
	}

	// Log request details in DEBUG mode
	if os.Getenv("DEBUG") != "" {
		c.Logger.Debugf("DEX API Request: POST %s, Body: %s", fullURL, bodyString)
	}

	// Execute request
	resp, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to execute request: %w", err)
	}
	defer resp.Body.Close()

	// Read response body
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Check HTTP status code
	if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("http error %d: %s", resp.StatusCode, string(data))
	}

	return data, nil
}

// request is a generic helper function that executes a request and parses the response.
// It automatically handles the BaseResponse wrapper and error checking.
func request[T any](c *Client, ctx context.Context, method, path string, params url.Values, auth bool) ([]T, error) {
	var data []byte
	var err error

	if method == "GET" {
		data, err = c.Do(ctx, method, path, params, auth)
	} else if method == "POST" {
		data, err = c.DoPost(ctx, path, nil, auth)
	} else {
		return nil, fmt.Errorf("unsupported HTTP method: %s", method)
	}

	if err != nil {
		return nil, err
	}

	// Parse response
	var baseResp BaseResponse[T]
	if err := json.Unmarshal(data, &baseResp); err != nil {
		return nil, fmt.Errorf("failed to unmarshal response: %w", err)
	}

	// Check API error code
	if baseResp.Code != "0" {
		return nil, &DexError{
			Code:    baseResp.Code,
			Message: baseResp.Msg,
		}
	}

	return baseResp.Data, nil
}
