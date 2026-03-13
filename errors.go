package okx

import (
	"fmt"
)

// DexError represents an error returned by the OKX DEX API.
type DexError struct {
	Code    string `json:"code"`
	Message string `json:"msg"`
}

func (e *DexError) Error() string {
	return fmt.Sprintf("OKX DEX API error [%s]: %s", e.Code, e.Message)
}

// Common DEX API error codes
const (
	ErrChainNotSupported        = "51001" // Chain not supported
	ErrTokenNotSupported        = "51002" // Token not supported
	ErrInsufficientLiquidity    = "51003" // Insufficient liquidity
	ErrSlippageExceeded         = "51004" // Slippage exceeded
	ErrInvalidAmount            = "51005" // Invalid amount
	ErrInvalidAddress           = "51006" // Invalid address
	ErrInvalidChainID           = "51007" // Invalid chain ID
	ErrPriceImpactTooHigh       = "51013" // Price impact too high
	ErrQuoteExpired             = "51014" // Quote expired
	ErrInsufficientAllowance    = "51015" // Insufficient token allowance
	ErrSwapFailed               = "51016" // Swap transaction failed
	ErrInvalidSlippage          = "51017" // Invalid slippage value
	ErrRateLimitExceeded        = "50011" // Rate limit exceeded
	ErrInvalidAPIKey            = "50100" // Invalid API key
	ErrInvalidSignature         = "50113" // Invalid signature
	ErrSystemBusy               = "50001" // System busy
)

// IsDexError checks if an error is a DexError.
func IsDexError(err error) bool {
	_, ok := err.(*DexError)
	return ok
}

// GetErrorCode extracts the error code from a DexError.
// Returns empty string if err is not a DexError.
func GetErrorCode(err error) string {
	if dexErr, ok := err.(*DexError); ok {
		return dexErr.Code
	}
	return ""
}
