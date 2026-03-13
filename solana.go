package okx

import (
	"context"
	"net/url"
)

// GetSolanaSwapInstruction retrieves the swap instruction data for Solana.
// Solana uses a different transaction model than EVM chains with instructions instead of calldata.
//
// Endpoint: GET /api/v6/dex/aggregator/swap-instruction
// Parameters defined in SolanaSwapRequest
// Authentication: Required
//
// Returns Solana-specific swap instruction data that must be decoded, signed, and broadcast.
//
// Note: Unlike EVM swaps, Solana swaps:
//   - Use "instructions" instead of calldata
//   - Use "computeUnitLimit" instead of gas limit
//   - Use "computeUnitPrice" for priority fees
//   - Native SOL address: So11111111111111111111111111111111111111112
func (c *Client) GetSolanaSwapInstruction(ctx context.Context, req *SolanaSwapRequest) (*SolanaSwapResponse, error) {
	path := "/api/v6/dex/aggregator/swap-instruction"
	params := url.Values{}
	params.Set("chainIndex", req.ChainID)
	params.Set("fromTokenAddress", req.FromTokenAddress)
	params.Set("toTokenAddress", req.ToTokenAddress)
	params.Set("amount", req.Amount)
	params.Set("slippage", req.Slippage)
	params.Set("userWalletAddress", req.UserWalletAddress)

	if req.ComputeUnitLimit != "" {
		params.Set("computeUnitLimit", req.ComputeUnitLimit)
	}
	if req.ComputeUnitPrice != "" {
		params.Set("computeUnitPrice", req.ComputeUnitPrice)
	}
	if req.FeeAccount != "" {
		params.Set("feeAccount", req.FeeAccount)
	}

	results, err := request[SolanaSwapResponse](c, ctx, "GET", path, params, true)
	if err != nil {
		return nil, err
	}

	if len(results) == 0 {
		return nil, &DexError{
			Code:    ErrSwapFailed,
			Message: "no Solana swap instruction data returned",
		}
	}

	return &results[0], nil
}
