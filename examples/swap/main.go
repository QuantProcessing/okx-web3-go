// Package main demonstrates how to build a swap transaction using the OKX DEX SDK.
//
// This example shows the complete flow for executing an EVM swap:
//   1. Get a swap quote
//   2. Check if token approval is needed (ERC-20 only)
//   3. Build the swap transaction data
//
// The SDK returns transaction data that you must sign and broadcast
// using a library like go-ethereum.
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

	// Configuration
	chainID := "8453"                                             // Base
	fromToken := "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE"  // ETH (native)
	toToken := "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913"      // USDC
	amount := "10000000000000000"                                 // 0.01 ETH (18 decimals)
	walletAddress := "0xYourWalletAddress"                        // Replace with your address

	// Step 1: Get quote
	fmt.Println("=== Step 1: Getting Quote ===")
	quote, err := client.GetQuote(ctx, &okx.QuoteRequest{
		ChainID:          chainID,
		FromTokenAddress: fromToken,
		ToTokenAddress:   toToken,
		Amount:           amount,
		Slippage:         "0.005",
	})
	if err != nil {
		log.Fatalf("Quote failed: %v", err)
	}
	fmt.Printf("  Quote: %s %s -> %s %s\n",
		quote.FromTokenAmount, quote.FromToken.TokenSymbol,
		quote.ToTokenAmount, quote.ToToken.TokenSymbol)
	fmt.Printf("  Price Impact: %s%%\n", quote.PriceImpactPercentage)
	fmt.Println()

	// Step 2: Build swap transaction
	// Note: For native tokens (ETH), no approval is needed.
	// For ERC-20 tokens, you would call GetApproveTransaction first.
	fmt.Println("=== Step 2: Building Swap Transaction ===")
	swap, err := client.BuildSwapData(ctx, &okx.SwapRequest{
		ChainID:           chainID,
		FromTokenAddress:  fromToken,
		ToTokenAddress:    toToken,
		Amount:            amount,
		Slippage:          "0.005",
		UserWalletAddress: walletAddress,
	})
	if err != nil {
		log.Fatalf("Swap build failed: %v", err)
	}

	fmt.Printf("  From:  %s\n", swap.Tx.From)
	fmt.Printf("  To:    %s\n", swap.Tx.To)
	fmt.Printf("  Value: %s\n", swap.Tx.Value)
	fmt.Printf("  Data:  %s...\n", truncate(swap.Tx.Data, 40))
	fmt.Printf("  Output: %s %s\n",
		swap.RouterResult.ToTokenAmount,
		swap.RouterResult.ToToken.TokenSymbol)
	fmt.Println()

	fmt.Println("=== Next Steps ===")
	fmt.Println("  1. Sign the transaction with your private key using go-ethereum")
	fmt.Println("  2. Broadcast the signed transaction to the network")
	fmt.Println("  3. Wait for confirmation")
}

func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen]
}
