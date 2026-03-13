package okx

import (
	"context"
	"net/url"
)

// GetQuote retrieves a swap quote from OKX DEX aggregator.
//
// Endpoint: GET /api/v6/dex/aggregator/quote
// Parameters defined in QuoteRequest
// Authentication: Recommended (for better routing and to avoid rate limits)
//
// Returns the best quote with routing information, estimated gas, price impact, etc.
func (c *Client) GetQuote(ctx context.Context, req *QuoteRequest) (*QuoteResponse, error) {
	path := "/api/v6/dex/aggregator/quote"
	params := url.Values{}
	params.Set("chainIndex", req.ChainID)
	params.Set("fromTokenAddress", req.FromTokenAddress)
	params.Set("toTokenAddress", req.ToTokenAddress)
	params.Set("amount", req.Amount)
	params.Set("slippage", req.Slippage)

	if req.UserWalletAddress != "" {
		params.Set("userWalletAddress", req.UserWalletAddress)
	}

	results, err := request[QuoteResponse](c, ctx, "GET", path, params, true)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, &DexError{
			Code:    ErrInsufficientLiquidity,
			Message: "no quote available",
		}
	}

	return &results[0], nil
}

// GetApproveTransaction retrieves the transaction data needed to approve token spending.
// This is only required for ERC-20 tokens on EVM chains (not needed for native tokens like ETH).
//
// Endpoint: GET /api/v6/dex/aggregator/approve-transaction
// Parameters defined in ApproveTransactionRequest
// Authentication: Required
//
// Returns transaction data that can be signed and broadcast to approve the DEX router
// to spend the specified amount of tokens.
func (c *Client) GetApproveTransaction(ctx context.Context, req *ApproveTransactionRequest) (*ApproveTransactionResponse, error) {
	path := "/api/v6/dex/aggregator/approve-transaction"
	params := url.Values{}
	params.Set("chainIndex", req.ChainID)
	params.Set("tokenContractAddress", req.TokenContractAddress)
	params.Set("approveAmount", req.ApproveAmount)

	results, err := request[ApproveTransactionResponse](c, ctx, "GET", path, params, true)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, &DexError{
			Code:    ErrInvalidAmount,
			Message: "no approval transaction data returned",
		}
	}

	return &results[0], nil
}

// BuildSwapData builds the transaction data for executing a swap.
// The returned transaction data must be signed by the user and broadcast to the blockchain.
//
// Endpoint: GET /api/v6/dex/aggregator/swap
// Parameters defined in SwapRequest
// Authentication: Required
//
// Returns complete transaction data (calldata, gas estimates, routing info) ready to be signed.
// Note: This method does NOT execute the swap - it only builds the transaction data.
func (c *Client) BuildSwapData(ctx context.Context, req *SwapRequest) (*SwapResponse, error) {
	path := "/api/v6/dex/aggregator/swap"
	params := url.Values{}
	params.Set("chainIndex", req.ChainID)
	params.Set("fromTokenAddress", req.FromTokenAddress)
	params.Set("toTokenAddress", req.ToTokenAddress)
	params.Set("amount", req.Amount)
	params.Set("slippagePercent", req.Slippage)
	params.Set("userWalletAddress", req.UserWalletAddress)

	if req.ReferrerAddress != "" {
		params.Set("referrerAddress", req.ReferrerAddress)
	}
	if req.FeePercent != "" {
		params.Set("feePercent", req.FeePercent)
	}
	if req.PriceImpactProtectionPercentage != "" {
		params.Set("priceImpactProtectionPercentage", req.PriceImpactProtectionPercentage)
	}

	results, err := request[SwapResponse](c, ctx, "GET", path, params, true)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, &DexError{
			Code:    ErrSwapFailed,
			Message: "no swap transaction data returned",
		}
	}

	return &results[0], nil
}

// GetSwapHistory retrieves the swap transaction history for a specific wallet address.
//
// Endpoint: GET /api/v6/dex/aggregator/swap/history
// Parameters defined in SwapHistoryRequest
// Authentication: Required
//
// Returns a list of historical swap transactions.
func (c *Client) GetSwapHistory(ctx context.Context, req *SwapHistoryRequest) ([]SwapHistoryResponse, error) {
	path := "/api/v6/dex/aggregator/swap/history"
	params := url.Values{}
	params.Set("chainIndex", req.ChainID)
	params.Set("userWalletAddress", req.UserWalletAddress)

	if req.Page != "" {
		params.Set("page", req.Page)
	}
	if req.Limit != "" {
		params.Set("limit", req.Limit)
	}

	return request[SwapHistoryResponse](c, ctx, "GET", path, params, true)
}
