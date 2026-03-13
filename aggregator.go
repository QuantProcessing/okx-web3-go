package okx

import (
	"context"
	"net/url"
)

// GetSupportedChains retrieves the list of blockchains supported by OKX DEX.
//
// Endpoint: GET /api/v6/dex/aggregator/supported/chain
// Authentication: Optional (recommended to avoid rate limits)
//
// Returns a list of supported chains with their metadata (chain ID, name, native token, etc.)
func (c *Client) GetSupportedChains(ctx context.Context) ([]SupportedChain, error) {
	path := "/api/v6/dex/aggregator/supported/chain"
	return request[SupportedChain](c, ctx, "GET", path, nil, true)
}

// GetTokens retrieves the list of tokens supported on a specific chain.
//
// Endpoint: GET /api/v6/dex/aggregator/all-tokens
// Parameters:
//   - chainID: Chain ID (e.g., "1" for Ethereum, "501" for Solana)
//
// Authentication: Optional (recommended to avoid rate limits)
//
// Returns a list of supported tokens on the specified chain.
func (c *Client) GetTokens(ctx context.Context, chainID string) ([]Token, error) {
	path := "/api/v6/dex/aggregator/all-tokens"
	params := url.Values{}
	params.Set("chainIndex", chainID)

	return request[Token](c, ctx, "GET", path, params, true)
}

// GetLiquiditySources retrieves the list of liquidity sources (DEXs) supported on a specific chain.
//
// Endpoint: GET /api/v6/dex/aggregator/all-liquidity
// Parameters:
//   - chainID: Chain ID (e.g., "1" for Ethereum, "56" for BNB Chain)
//
// Authentication: Optional (recommended to avoid rate limits)
//
// Returns a list of DEXs and liquidity providers available on the specified chain.
func (c *Client) GetLiquiditySources(ctx context.Context, chainID string) ([]LiquiditySource, error) {
	path := "/api/v6/dex/aggregator/all-liquidity"
	params := url.Values{}
	params.Set("chainIndex", chainID)

	return request[LiquiditySource](c, ctx, "GET", path, params, true)
}
