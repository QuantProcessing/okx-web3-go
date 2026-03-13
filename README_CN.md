# OKX DEX SDK 使用指南

[English](README.md) | 中文

OKX DEX SDK 是一个用于集成 OKX DEX 聚合器的 Go 客户端库，支持跨多链（EVM、Solana 等）的代币交易和查询功能。

## 特性

- ✅ **多链支持**：Ethereum、Base、BNB Chain、Arbitrum、Polygon、Solana、Sui、TON 等 20+ 区块链
- ✅ **流动性聚合**：从 500+ DEX 聚合最优流动性
- ✅ **智能路由**：自动选择最佳交易路径和价格
- ✅ **完整功能**：支持代币查询、报价获取、交易授权、Swap 数据构建、历史查询
- ✅ **类型安全**：完整的类型定义和错误处理
- ✅ **依赖精简**：仅依赖 `go.uber.org/zap` 日志库

## 快速开始

### 安装

```bash
go get github.com/QuantProcessing/okx-web3-go
```

### 初始化客户端

```go
package main

import (
    "context"
    "fmt"
    "log"

    okx "github.com/QuantProcessing/okx-web3-go"
)

func main() {
    // 创建客户端
    client := okx.NewClient()

    // 设置凭证（仅需鉴权接口需要）
    client.WithCredentials(
        "your-api-key",
        "your-secret-key",
        "your-passphrase",
        "your-project-id", // 可选
    )

    ctx := context.Background()

    // 查询支持的链
    chains, err := client.GetSupportedChains(ctx)
    if err != nil {
        log.Fatal(err)
    }

    for _, chain := range chains {
        fmt.Printf("Chain: %s (ID: %s)\n", chain.ChainName, chain.ChainID)
    }
}
```

## 常用链 ID 参考

| 链名称 | Chain ID | 原生代币地址 |
|--------|----------|-------------|
| Ethereum | `1` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Base | `8453` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| BNB Chain | `56` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Polygon | `137` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Arbitrum One | `42161` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Optimism | `10` | `0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE` |
| Solana | `501` | `So11111111111111111111111111111111111111112` |

完整链列表可通过 `GetSupportedChains()` 动态获取。

## 使用示例

### 1. 查询代币列表

```go
// 获取 Base 链上的代币列表
tokens, err := client.GetTokens(ctx, "8453")
if err != nil {
    log.Fatal(err)
}

for _, token := range tokens {
    fmt.Printf("%s (%s): %s\n", token.TokenName, token.TokenSymbol, token.TokenAddress)
}
```

### 2. 获取 Swap 报价

```go
// 查询 Base 链上 0.1 ETH -> USDC 的报价
quote, err := client.GetQuote(ctx, &okx.QuoteRequest{
    ChainID:          "8453",  // Base
    FromTokenAddress: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",  // ETH
    ToTokenAddress:   "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",  // USDC
    Amount:           "100000000000000000",  // 0.1 ETH (18 decimals)
    Slippage:         "0.005",  // 0.5% 滑点
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("报价：\n")
fmt.Printf("  输入: %s %s\n", quote.FromTokenAmount, quote.FromToken.TokenSymbol)
fmt.Printf("  输出: %s %s\n", quote.ToTokenAmount, quote.ToToken.TokenSymbol)
fmt.Printf("  价格影响: %s%%\n", quote.PriceImpactPercentage)
fmt.Printf("  预估 Gas: %s\n", quote.EstimatedGas)
```

### 3. 构建授权交易（仅 ERC-20 代币）

```go
// 注意：原生代币（ETH、BNB 等）无需授权，只有 ERC-20 代币需要
approvalTx, err := client.GetApproveTransaction(ctx, &okx.ApproveTransactionRequest{
    ChainID:              "8453",
    TokenContractAddress: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", // USDC
    ApproveAmount:        "1000000", // 1 USDC (6 decimals)
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("授权交易数据：\n")
fmt.Printf("  To: %s\n", approvalTx.DexContractAddress)
fmt.Printf("  Data: %s\n", approvalTx.Data)
// 此处需要用户签名 approvalTx.Data 并广播上链
```

### 4. 构建 Swap 交易

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

fmt.Printf("Swap 交易数据：\n")
fmt.Printf("  From: %s\n", swapTx.Tx.From)
fmt.Printf("  To: %s\n", swapTx.Tx.To)
fmt.Printf("  Data: %s\n", swapTx.Tx.Data)
fmt.Printf("  Value: %s\n", swapTx.Tx.Value)
fmt.Printf("  输出: %s %s\n", swapTx.RouterResult.ToTokenAmount, swapTx.RouterResult.ToToken.TokenSymbol)
// 此处需要用户签名 swapTx.Tx 并广播上链
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
    ComputeUnitLimit:  "200000", // 可选，类似 EVM 的 gas limit
})

if err != nil {
    log.Fatal(err)
}

fmt.Printf("Solana Swap 指令：\n")
fmt.Printf("  Instructions: %s\n", solanaSwap.Instructions)
fmt.Printf("  Compute Units: %s\n", solanaSwap.ComputeUnitLimit)
// 此处需要用 Solana SDK 解析 Instructions，签名并广播
```

### 6. 查询交易历史

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

## 重要说明

### 职责边界

⚠️ **本 SDK 仅负责 API 调用，不涉及区块链交易签名和广播**：

- **EVM 链**：需要外部使用 `go-ethereum` 或类似库进行交易签名
- **Solana 链**：需要外部使用 Solana SDK 进行交易签名
- SDK 返回的是**交易数据**（calldata/instructions），需由调用方签名并广播上链

### 精度处理

- API 返回的金额/价格均为 `string` 类型（避免浮点数精度损失）
- 金额单位为**最小单位**：
  - ETH: wei (18 decimals)
  - USDC (EVM): 最小单位 (6 decimals)
  - SOL: lamports (9 decimals)
- 在进行金额计算时，建议使用 `shopspring/decimal` 库

### 滑点设置

- 滑点值为字符串格式的小数：`"0.005"` 表示 0.5%
- 推荐范围：0.1% - 1% (即 `"0.001"` - `"0.01"`)
- 滑点过小可能导致交易失败；过大可能遭受损失

### 错误处理

```go
quote, err := client.GetQuote(ctx, req)
if err != nil {
    // 检查是否为 DEX 特定错误
    if okx.IsDexError(err) {
        code := okx.GetErrorCode(err)
        switch code {
        case okx.ErrInsufficientLiquidity:
            fmt.Println("流动性不足")
        case okx.ErrSlippageExceeded:
            fmt.Println("滑点超限")
        case okx.ErrPriceImpactTooHigh:
            fmt.Println("价格影响过大")
        default:
            fmt.Printf("DEX 错误 [%s]: %v\n", code, err)
        }
    } else {
        fmt.Printf("其他错误: %v\n", err)
    }
    return
}
```

### 环境变量

- `DEBUG`: 设置为任意非空值启用调试日志

## 常见问题

### 1. 如何判断代币是否需要授权？

- **原生代币**（ETH、BNB、SOL 等）：**不需要**授权
- **ERC-20 代币**（USDC、USDT、DAI 等）：**需要**授权
- Solana 代币：大部分**不需要**显式授权（通过 Associated Token Account）

### 2. 授权额度应该设置多少？

- **单次授权**：设置为本次交易的确切金额
- **无限授权**：设置为最大值（节省 gas，但存在安全风险）
  - ERC-20 最大值：`115792089237316195423570985008687907853269984665640564039457584007913129639935`

### 3. 如何计算金额的最小单位？

```go
// 示例：将 1.5 USDC 转换为最小单位
// USDC decimals = 6
amount := new(big.Float).SetFloat64(1.5)
multiplier := new(big.Float).SetInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(6), nil))
result := new(big.Float).Mul(amount, multiplier)
amountStr, _ := result.Text('f', 0) // "1500000"
```

### 4. 价格影响百分比如何解读？

- `< 1%`：正常交易
- `1% - 3%`：中等影响，需谨慎
- `> 3%`：高影响，建议拆分订单或等待流动性改善

### 5. 支持哪些 DEX？

常见 DEX 包括：
- **Ethereum**: Uniswap, SushiSwap, Curve, Balancer, 1inch, 0x
- **Base**: Uniswap V3, Aerodrome, BaseSwap
- **BNB Chain**: PancakeSwap, Biswap, ApeSwap
- **Solana**: Jupiter, Orca, Raydium

完整列表通过 `GetLiquiditySources(ctx, chainID)` 查询。

## API 参考

| 方法 | 功能 | 鉴权 |
|------|------|------|
| `GetSupportedChains` | 获取支持的链列表 | 可选 |
| `GetTokens` | 获取链上代币列表 | 可选 |
| `GetLiquiditySources` | 获取流动性来源 | 可选 |
| `GetQuote` | 获取 Swap 报价 | 必需 |
| `GetApproveTransaction` | 构建授权交易 | 必需 |
| `BuildSwapData` | 构建 Swap 交易 | 必需 |
| `GetSwapHistory` | 查询历史记录 | 必需 |
| `GetSolanaSwapInstruction` | 获取 Solana Swap 指令 | 必需 |

## 许可证

本项目基于 [MIT License](LICENSE) 开源。
