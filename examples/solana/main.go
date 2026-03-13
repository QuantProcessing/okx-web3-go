// Package main demonstrates Solana swap instruction retrieval using the OKX DEX SDK.
//
// Solana uses a different transaction model than EVM chains:
//   - Instructions instead of calldata
//   - Compute units instead of gas
//   - Priority fees via computeUnitPrice
//
// Usage:
//
//	export OKX_API_KEY="your-api-key"
//	export OKX_SECRET_KEY="your-secret-key"
//	export OKX_PASSPHRASE="your-passphrase"
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
	apiKey := os.Getenv("OKX_API_KEY")
	secretKey := os.Getenv("OKX_SECRET_KEY")
	passphrase := os.Getenv("OKX_PASSPHRASE")
	projectID := os.Getenv("OKX_PROJECT_ID")

	if apiKey == "" || secretKey == "" || passphrase == "" {
		log.Fatal("Please set OKX_API_KEY, OKX_SECRET_KEY, and OKX_PASSPHRASE environment variables")
	}

	client := okx.NewClient().WithCredentials(apiKey, secretKey, passphrase, projectID)
	ctx := context.Background()

	// Solana swap: 0.1 SOL -> USDC
	fmt.Println("=== Solana Swap: 0.1 SOL -> USDC ===")
	result, err := client.GetSolanaSwapInstruction(ctx, &okx.SolanaSwapRequest{
		ChainID:           "501",                                             // Solana mainnet
		FromTokenAddress:  "So11111111111111111111111111111111111111112",       // SOL
		ToTokenAddress:    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",   // USDC
		Amount:            "100000000",                                       // 0.1 SOL (9 decimals)
		Slippage:          "0.005",
		UserWalletAddress: "YourSolanaWalletAddress", // Replace with your address
		ComputeUnitLimit:  "200000",
	})
	if err != nil {
		if okx.IsDexError(err) {
			fmt.Printf("DEX Error [%s]: %v\n", okx.GetErrorCode(err), err)
			return
		}
		log.Fatalf("Failed: %v", err)
	}

	fmt.Printf("  From:             %s %s\n", result.FromTokenAmount, result.FromToken.TokenSymbol)
	fmt.Printf("  To:               %s %s\n", result.ToTokenAmount, result.ToToken.TokenSymbol)
	fmt.Printf("  Min Receive:      %s\n", result.MinReceiveAmount)
	fmt.Printf("  Price Impact:     %s%%\n", result.PriceImpactPercentage)
	fmt.Printf("  Compute Units:    %s\n", result.ComputeUnitLimit)
	fmt.Printf("  Instructions:     %s...\n", truncate(result.Instructions, 40))
	fmt.Println()
	fmt.Println("Next: Decode instructions with Solana SDK, sign, and broadcast")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
