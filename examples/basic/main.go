// Package main demonstrates basic usage of the OKX DEX SDK.
//
// This example shows how to:
//   - Create a client and set credentials
//   - Query supported chains
//   - Get token list on a specific chain
//   - Get a swap quote
//
// Usage:
//
//	export OKX_API_KEY="your-api-key"
//	export OKX_SECRET_KEY="your-secret-key"
//	export OKX_PASSPHRASE="your-passphrase"
//	export OKX_PROJECT_ID="your-project-id"  # optional
//	go run .
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	okx "github.com/QuantProcessing/okx-web3-go"
)

func main() {
	// Read credentials from environment
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_SECRET_KEY")
	passphrase := os.Getenv("OKX_PASSPHRASE")
	projectID := os.Getenv("OKX_PROJECT_ID")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("Please set OKX_API_KEY, OKX_SECRET_KEY, and OKX_PASSPHRASE environment variables")
	}

	// Create client with credentials
	client := okx.NewClient().WithCredentials(apiKey, secretKey, passphrase, projectID)

	ctx := context.Background()

	// 1. Query supported chains
	fmt.Println("=== Supported Chains ===")
	chains, err := client.GetSupportedChains(ctx)
	if err != nil {
		log.Fatalf("Failed to get supported chains: %v", err)
	}
	for _, chain := range chains {
		fmt.Printf("  %s (ID: %s, Native: %s)\n", chain.ChainName, chain.ChainID, chain.NativeTokenSymbol)
	}
	fmt.Println()

	// 2. Query tokens on Base chain
	fmt.Println("=== Tokens on Base (first 5) ===")
	tokens, err := client.GetTokens(ctx, "8453")
	if err != nil {
		log.Fatalf("Failed to get tokens: %v", err)
	}
	limit := 5
	if len(tokens) < limit {
		limit = len(tokens)
	}
	for _, token := range tokens[:limit] {
		fmt.Printf("  %s (%s): %s\n", token.TokenName, token.TokenSymbol, token.TokenAddress)
	}
	fmt.Printf("  ... and %d more tokens\n\n", len(tokens)-limit)

	// 3. Get a swap quote: 0.01 ETH -> USDC on Base
	fmt.Println("=== Swap Quote: 0.01 ETH -> USDC (Base) ===")
	quote, err := client.GetQuote(ctx, &okx.QuoteRequest{
		ChainID:          "8453",
		FromTokenAddress: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // ETH
		ToTokenAddress:   "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", // USDC
		Amount:           "10000000000000000",                           // 0.01 ETH
		Slippage:         "0.005",
	})
	if err != nil {
		if okx.IsDexError(err) {
			fmt.Printf("  DEX Error [%s]: %v\n", okx.GetErrorCode(err), err)
		} else {
			log.Fatalf("Failed to get quote: %v", err)
		}
	} else {
		fmt.Printf("  Input:        %s %s\n", quote.FromTokenAmount, quote.FromToken.TokenSymbol)
		fmt.Printf("  Output:       %s %s\n", quote.ToTokenAmount, quote.ToToken.TokenSymbol)
		fmt.Printf("  Price Impact: %s%%\n", quote.PriceImpactPercentage)
		fmt.Printf("  Estimated Gas: %s\n", quote.EstimatedGas)
	}
}
