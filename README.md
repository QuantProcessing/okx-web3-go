# OKX DEX SDK for Go

English | [中文](README_CN.md)

[![Go Reference](https://pkg.go.dev/badge/github.com/QuantProcessing/okx-web3-go.svg)](https://pkg.go.dev/github.com/QuantProcessing/okx-web3-go)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

A Go client library for the [OKX DEX Aggregator API](https://web3.okx.com/), supporting cross-chain token swaps on EVM chains, Solana, and more.

## Features

- ✅ **Multi-chain support**: Ethereum, Base, BNB Chain, Arbitrum, Polygon, Solana, Sui, TON, and 20+ blockchains
- ✅ **Liquidity aggregation**: Best price routing across 500+ DEXs
- ✅ **Smart routing**: Automatic split-routing for optimal execution
- ✅ **Complete API coverage**: Token queries, quotes, approvals, swap data building, and history
- ✅ **Type-safe**: Full type definitions and structured error handling
- ✅ **Minimal dependencies**: Only `go.uber.org/zap` for logging

## Installation

```bash
go get github.com/QuantProcessing/okx-web3-go
```

## Quick Start

```go
package main

import (
    "context"
    "fmt"
    "log"

    okx "github.com/QuantProcessing/okx-web3-go"
)

func main() {
    // Create client
    client := okx.NewClient()

    // Set credentials (required for authenticated endpoints)
    client.WithCredentials(
        "your-api-key",
        "your-secret-key",
        "your-passphrase",
        "your-project-id", // optional
    )

    ctx := context.Background()

    // Query supported chains
    chains, err := client.GetSupportedChains(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, chain := range chains {
        fmt.Printf("Chain: %s (ID: %s)\n", chain.ChainName, chain.ChainID)
    }
}
```

## Common Chain IDs

| Chain | Chain ID | Native Token Address |
|-------|----------|---------------------|
| Ethereum | `1` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Base | `8453` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| BNB Chain | `56` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Polygon | `137` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Arbitrum One | `42161` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Optimism | `10` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Solana | `501` | `So11111111111111111111111111111111111111112` |

Full chain list available via `GetSupportedChains()`.

## Usage Examples

### 1. Query Token List

```go
// Get tokens on Base chain
tokens, err := client.GetTokens(ctx, "8453")
if err != nil {
    log.Fatal(err)
}

for _, token := range tokens {
    fmt.Printf("%s (%s): %s\n", token.TokenName, token.TokenSymbol, token.TokenAddress)
}
```

### 2. Get Swap Quote

```go
// Query a quote for 0.1 ETH -> USDC on Base
quote, err := client.GetQuote(ctx, &okx.QuoteRequest{
    ChainID:          "8453",  // Base
    FromTokenAddress: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",  // ETH
    ToTokenAddress:   "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",  // USDC
    Amount:           "100000000000000000",  // 0.1 ETH (18 decimals)
    Slippage:         "0.005",  // 0.5% slippage
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Quote:\n")
fmt.Printf("  Input:  %s %s\n", quote.FromTokenAmount, quote.FromToken.TokenSymbol)
fmt.Printf("  Output: %s %s\n", quote.ToTokenAmount, quote.ToToken.TokenSymbol)
fmt.Printf("  Price Impact: %s%%\n", quote.PriceImpactPercentage)
fmt.Printf("  Estimated Gas: %s\n", quote.EstimatedGas)
```

### 3. Build Approval Transaction (ERC-20 tokens only)

```go
// Note: Native tokens (ETH, BNB, etc.) do NOT need approval
approvalTx, err := client.GetApproveTransaction(ctx, &okx.ApproveTransactionRequest{
    ChainID:              "8453",
    TokenContractAddress: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", // USDC
    ApproveAmount:        "1000000", // 1 USDC (6 decimals)
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Approval TX:\n")
fmt.Printf("  To: %s\n", approvalTx.DexContractAddress)
fmt.Printf("  Data: %s\n", approvalTx.Data)
// Sign approvalTx.Data and broadcast to chain
```

### 4. Build Swap Transaction

```go
swapTx, err := client.BuildSwapData(ctx, &okx.SwapRequest{
    ChainID:           "8453",
    FromTokenAddress:  "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE", // ETH
    ToTokenAddress:    "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", // USDC
    Amount:            "100000000000000000", // 0.1 ETH
    Slippage:          "0.005",
    UserWalletAddress: "0xYourWalletAddress",
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Swap TX:\n")
fmt.Printf("  From: %s\n", swapTx.Tx.From)
fmt.Printf("  To: %s\n", swapTx.Tx.To)
fmt.Printf("  Data: %s\n", swapTx.Tx.Data)
fmt.Printf("  Value: %s\n", swapTx.Tx.Value)
fmt.Printf("  Output: %s %s\n", swapTx.RouterResult.ToTokenAmount, swapTx.RouterResult.ToToken.TokenSymbol)
// Sign swapTx.Tx and broadcast to chain
```

### 5. Solana Swap

```go
solanaSwap, err := client.GetSolanaSwapInstruction(ctx, &okx.SolanaSwapRequest{
    ChainID:           "501", // Solana mainnet
    FromTokenAddress:  "So11111111111111111111111111111111111111112", // SOL
    ToTokenAddress:    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v", // USDC
    Amount:            "1000000000", // 1 SOL (9 decimals)
    Slippage:          "0.005",
    UserWalletAddress: "YourSolanaWalletAddress",
    ComputeUnitLimit:  "200000", // optional, similar to EVM gas limit
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Solana Swap:\n")
fmt.Printf("  Instructions: %s\n", solanaSwap.Instructions)
fmt.Printf("  Compute Units: %s\n", solanaSwap.ComputeUnitLimit)
// Decode instructions with Solana SDK, sign, and broadcast
```

### 6. Query Swap History

```go
history, err := client.GetSwapHistory(ctx, &okx.SwapHistoryRequest{
    ChainID:           "8453",
    UserWalletAddress: "0xYourWalletAddress",
    Page:              "1",
    Limit:             "20",
})

if err != nil {
    log.Fatal(err)
}

for _, tx := range history {
    fmt.Printf("TX: %s, Status: %s, %s %s -> %s %s\n",
        tx.TxHash, tx.Status,
        tx.FromTokenAmount, tx.FromToken,
        tx.ToTokenAmount, tx.ToToken,
    )
}
```

## Important Notes

### Scope

⚠️ **This SDK handles API calls only — it does NOT sign or broadcast transactions:**

- **EVM chains**: Use `go-ethereum` or similar libraries for transaction signing
- **Solana**: Use a Solana SDK for transaction signing
- The SDK returns **transaction data** (calldata/instructions) that must be signed and broadcast by the caller

### Amount Precision

- All amounts/prices are returned as `string` to avoid floating-point precision loss
- Amounts are in **smallest units**:
  - ETH: wei (18 decimals)
  - USDC (EVM): smallest unit (6 decimals)
  - SOL: lamports (9 decimals)
- Use `shopspring/decimal` for amount calculations

### Slippage

- Slippage is a decimal string: `"0.005"` = 0.5%
- Recommended range: 0.1% – 1% (`"0.001"` – `"0.01"`)
- Too low → transaction may fail; too high → may suffer losses

### Error Handling

```go
quote, err := client.GetQuote(ctx, req)
if err != nil {
    if okx.IsDexError(err) {
        code := okx.GetErrorCode(err)
        switch code {
        case okx.ErrInsufficientLiquidity:
            fmt.Println("Insufficient liquidity")
        case okx.ErrSlippageExceeded:
            fmt.Println("Slippage exceeded")
        case okx.ErrPriceImpactTooHigh:
            fmt.Println("Price impact too high")
        default:
            fmt.Printf("DEX error [%s]: %v\n", code, err)
        }
    } else {
        fmt.Printf("Other error: %v\n", err)
    }
    return
}
```

### Environment Variables

- `DEBUG`: Set to any non-empty value to enable debug logging

## FAQ

### 1. Which tokens need approval?

- **Native tokens** (ETH, BNB, SOL, etc.): **No** approval needed
- **ERC-20 tokens** (USDC, USDT, DAI, etc.): **Approval required**
- **Solana tokens**: Most do **not** need explicit approval (uses Associated Token Accounts)

### 2. How much should I approve?

- **Exact amount**: Set to the trade amount (more secure)
- **Unlimited**: Set to max uint256 (saves gas on repeated swaps, but carries security risk)
  - Max ERC-20: `115792089237316195423570985008687907853269984665640564039457584007913129639935`

### 3. How to convert amounts to smallest units?

```go
// Example: convert 1.5 USDC to smallest unit (6 decimals)
amount := new(big.Float).SetFloat64(1.5)
multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))
result := new(big.Float).Mul(amount, multiplier)
amountStr, _ := result.Text('f', 0) // "1500000"
```

### 4. How to interpret price impact?

- `< 1%`: Normal trade
- `1% – 3%`: Moderate impact, proceed with caution
- `> 3%`: High impact — consider splitting the order or waiting for better liquidity

### 5. Which DEXs are supported?

Common DEXs include:
- **Ethereum**: Uniswap, SushiSwap, Curve, Balancer, 1inch, 0x
- **Base**: Uniswap V3, Aerodrome, BaseSwap
- **BNB Chain**: PancakeSwap, Biswap, ApeSwap
- **Solana**: Jupiter, Orca, Raydium

Full list via `GetLiquiditySources(ctx, chainID)`.

## API Reference

| Method | Description | Auth |
|--------|-------------|------|
| `GetSupportedChains` | List supported blockchains | Optional |
| `GetTokens` | List tokens on a chain | Optional |
| `GetLiquiditySources` | List DEX liquidity sources | Optional |
| `GetQuote` | Get swap quote | Required |
| `GetApproveTransaction` | Build token approval TX | Required |
| `BuildSwapData` | Build swap transaction data | Required |
| `GetSwapHistory` | Query swap history | Required |
| `GetSolanaSwapInstruction` | Get Solana swap instruction | Required |

## License

This project is licensed under the [MIT License](LICENSE).
