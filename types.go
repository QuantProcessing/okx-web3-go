package okx

// BaseResponse is the standard response wrapper for OKX DEX API.
type BaseResponse[T any] struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []T    `json:"data"`
}

// -------------------------
// Chain and Token Types
// -------------------------

// SupportedChain represents a blockchain supported by OKX DEX.
type SupportedChain struct {
	ChainID            ChainIDType `json:"chainId"`            // Chain ID - API returns as number or string
	ChainName          string `json:"chainName"`          // Chain name (e.g., "Ethereum", "Solana")
	DexTokenApproveAddress string `json:"dexTokenApproveAddress"` // Router/approval contract address
	NativeTokenAddress string `json:"nativeTokenAddress"` // Native token address
	NativeTokenSymbol  string `json:"nativeTokenSymbol"`  // Native token symbol (e.g., "ETH", "SOL")
	NativeTokenDecimals string `json:"nativeTokenDecimals"` // Native token decimals
}

// Token represents a token supported by OKX DEX on a specific chain.
type Token struct {
	ChainID        string `json:"chainId"`
	TokenName      string `json:"tokenName"`
	TokenSymbol    string `json:"tokenSymbol"`
	TokenAddress   string `json:"tokenContractAddress"`
	Decimals       string `json:"decimal"` // Token decimals
	IsHoneyPot     bool   `json:"isHoneyPot,omitempty"`
	MakerDaoVerify bool   `json:"makerDaoVerify,omitempty"`
}

// LiquiditySource represents a DEX or liquidity provider.
type LiquiditySource struct {
	ChainID       string `json:"chainId"`
	DexName       string `json:"dexName"`       // DEX name (e.g., "Uniswap", "PancakeSwap")
	DexProtocol   string `json:"dexProtocol"`   // Protocol identifier
	IsDynamicFee  bool   `json:"isDynamicFee"`  // Whether the DEX uses dynamic fees
	ProtocolVersion string `json:"protocolVersion,omitempty"`
}

// -------------------------
// Quote Types
// -------------------------

// QuoteRequest represents a request to get a swap quote.
type QuoteRequest struct {
	ChainID          string `json:"chainId"`
	FromTokenAddress string `json:"fromTokenAddress"`
	ToTokenAddress   string `json:"toTokenAddress"`
	Amount           string `json:"amount"`           // Amount in smallest unit (e.g., wei for ETH)
	Slippage         string `json:"slippage"`         // Slippage tolerance (e.g., "0.005" for 0.5%)
	UserWalletAddress string `json:"userWalletAddress,omitempty"` // Optional, for personalized routing
}

// QuoteResponse represents the response from a quote request.
type QuoteResponse struct {
	ChainID              string       `json:"chainId"`
	FromToken            TokenInfo    `json:"fromToken"`
	ToToken              TokenInfo    `json:"toToken"`
	FromTokenAmount      string       `json:"fromTokenAmount"`      // Input amount
	ToTokenAmount        string       `json:"toTokenAmount"`        // Output amount
	EstimatedGas         string       `json:"estimatedGas"`         // Estimated gas units
	GasPrice             string       `json:"gasPrice,omitempty"`   // Gas price (EVM only)
	GasUnitPrice         string       `json:"gasUnitPrice,omitempty"` // Alternative gas price field
	PriceImpactPercentage string      `json:"priceImpactPercentage"` // Price impact (e.g., "0.5")
	TradeFee             string       `json:"tradeFee,omitempty"`
	RouteList            []RouteInfo  `json:"routerList,omitempty"`  // Swap route details
	QuoteCompareList     []QuoteCompare `json:"quoteCompareList,omitempty"` // Alternative quotes
}

// TokenInfo contains detailed information about a token in a quote.
type TokenInfo struct {
	Decimals           string `json:"decimal"`
	IsHoneyPot         bool   `json:"isHoneyPot"` // Changed to bool
	TokenSymbol        string `json:"tokenSymbol"`
	TokenAddress       string `json:"tokenContractAddress"`
	TokenUnitPrice     string `json:"tokenUnitPrice,omitempty"` // USD price per token
}

// RouteInfo represents a swap route through one or more DEXs.
type RouteInfo struct {
	Router           string       `json:"router"`           // Router contract address
	RouterPercent    string       `json:"routerPercent"`    // Percentage of swap through this route
	SubRouterList    []SubRouter  `json:"subRouterList"`
	TradeFee         string       `json:"tradeFee,omitempty"`
}

// SubRouter represents a single hop in a swap route.
type SubRouter struct {
	DexProtocol     string   `json:"dexProtocol"`
	FromToken       string   `json:"fromToken"`
	ToToken         string   `json:"toToken"`
	FromTokenAmount string   `json:"fromTokenAmount"`
	ToTokenAmount   string   `json:"toTokenAmount"`
	TradeFee        string   `json:"tradeFee,omitempty"`
}

// QuoteCompare provides comparison quotes from different routing strategies.
type QuoteCompare struct {
	DexName         string `json:"dexName"`
	ToTokenAmount   string `json:"toTokenAmount"`
	TradeFee        string `json:"tradeFee,omitempty"`
}

// -------------------------
// Approval Types
// -------------------------

// ApproveTransactionRequest represents a request to get approval transaction data.
type ApproveTransactionRequest struct {
	ChainID              string `json:"chainId"`
	TokenContractAddress string `json:"tokenContractAddress"`
	ApproveAmount        string `json:"approveAmount"` // Amount to approve, or max uint256 for unlimited
}

// ApproveTransactionResponse represents the response with approval transaction data.
type ApproveTransactionResponse struct {
	ChainID              string `json:"chainId"`
	DexContractAddress   string `json:"dexContractAddress"`   // Spender address (router)
	TokenContractAddress string `json:"tokenContractAddress"`
	ApproveAmount        string `json:"approveAmount"`
	Data                 string `json:"data"`                 // Transaction calldata
	GasPrice             string `json:"gasPrice,omitempty"`   // Suggested gas price
	GasLimit             string `json:"gasLimit,omitempty"`   // Suggested gas limit
}

// -------------------------
// Swap Types
// -------------------------

// SwapRequest represents a request to build swap transaction data.
type SwapRequest struct {
	ChainID           string  `json:"chainId"`
	FromTokenAddress  string  `json:"fromTokenAddress"`
	ToTokenAddress    string  `json:"toTokenAddress"`
	Amount            string  `json:"amount"`
	Slippage          string  `json:"slippage"`
	UserWalletAddress string  `json:"userWalletAddress"`
	ReferrerAddress   string  `json:"referrerAddress,omitempty"`   // For fee sharing
	FeePercent        string  `json:"feePercent,omitempty"`        // Fee percentage for referrer
	PriceImpactProtectionPercentage string `json:"priceImpactProtectionPercentage,omitempty"` // Max acceptable price impact
}

// SwapResponse represents the response with swap transaction data.
type SwapResponse struct {
	RouterResult         RouterResult `json:"routerResult"`        // Router result object (V6 API)
	Tx                   TxData       `json:"tx"`                  // Transaction data to sign
}

// RouterResult contains detailed routing and token information from V6 API
type RouterResult struct {
	ChainIndex           string      `json:"chainIndex"`
	FromToken            TokenInfo   `json:"fromToken"`
	ToToken              TokenInfo   `json:"toToken"`
	FromTokenAmount      string      `json:"fromTokenAmount"`
	ToTokenAmount        string      `json:"toTokenAmount"`
	EstimateGasFee       string      `json:"estimateGasFee"`
	PriceImpactPercent   string      `json:"priceImpactPercent"`
	Router               string      `json:"router"`               // Router path
	SwapMode             string      `json:"swapMode"`             // "exactIn" or "exactOut"
	TradeFee             string      `json:"tradeFee,omitempty"`
	DexRouterList        []DexRouter `json:"dexRouterList,omitempty"`
}

// DexRouter represents routing through a specific DEX
type DexRouter struct {
	DexProtocol     DexProtocol `json:"dexProtocol"`
	FromToken       TokenInfo   `json:"fromToken"`
	ToToken         TokenInfo   `json:"toToken"`
	FromTokenIndex  string      `json:"fromTokenIndex"`
	ToTokenIndex    string      `json:"toTokenIndex"`
}

// DexProtocol contains DEX name and percentage
type DexProtocol struct {
	DexName string `json:"dexName"`
	Percent string `json:"percent"`
}

// TxData contains the transaction data ready to be signed and broadcast.
type TxData struct {
	From     string `json:"from"`               // User wallet address
	To       string `json:"to"`                 // Router contract address
	Data     string `json:"data"`               // Transaction calldata
	Value    string `json:"value"`              // Native token amount (for ETH swaps)
	GasPrice string `json:"gasPrice,omitempty"` // Gas price in wei
	GasLimit string `json:"gas,omitempty"`      // Gas limit
	MaxPriorityFeePerGas string `json:"maxPriorityFeePerGas,omitempty"` // EIP-1559
	MaxFeePerGas         string `json:"maxFeePerGas,omitempty"`         // EIP-1559
}

// -------------------------
// Swap History Types
// -------------------------

// SwapHistoryRequest represents a request to query swap history.
type SwapHistoryRequest struct {
	ChainID           string `json:"chainId"`
	UserWalletAddress string `json:"userWalletAddress"`
	Page              string `json:"page,omitempty"`  // Page number (default: 1)
	Limit             string `json:"limit,omitempty"` // Page size (default: 20, max: 100)
}

// SwapHistoryResponse represents a swap transaction in history.
type SwapHistoryResponse struct {
	ChainID         string `json:"chainId"`
	TxHash          string `json:"txId"`
	BlockNumber     string `json:"blockNumber,omitempty"`
	FromToken       string `json:"fromToken"`
	ToToken         string `json:"toToken"`
	FromTokenAmount string `json:"fromTokenAmount"`
	ToTokenAmount   string `json:"toTokenAmount"`
	Status          string `json:"status"` // "success", "failed", "pending"
	Timestamp       string `json:"txTime"`
	GasUsed         string `json:"gasUsed,omitempty"`
}

// -------------------------
// Solana-Specific Types
// -------------------------

// SolanaSwapRequest represents a request to get Solana swap instruction.
type SolanaSwapRequest struct {
	ChainID            string `json:"chainId"` // "501" for Solana mainnet
	FromTokenAddress   string `json:"fromTokenAddress"`
	ToTokenAddress     string `json:"toTokenAddress"`
	Amount             string `json:"amount"`
	Slippage           string `json:"slippage"`
	UserWalletAddress  string `json:"userWalletAddress"`
	ComputeUnitLimit   string `json:"computeUnitLimit,omitempty"`   // Gas limit equivalent for Solana
	ComputeUnitPrice   string `json:"computeUnitPrice,omitempty"`   // Priority fee
	FeeAccount         string `json:"feeAccount,omitempty"`         // For fee sharing
}

// SolanaSwapResponse represents the response with Solana swap instruction.
type SolanaSwapResponse struct {
	ChainID              string    `json:"chainId"`
	FromToken            TokenInfo `json:"fromToken"`
	ToToken              TokenInfo `json:"toToken"`
	FromTokenAmount      string    `json:"fromTokenAmount"`
	ToTokenAmount        string    `json:"toTokenAmount"`
	MinReceiveAmount     string    `json:"minReceiveAmount"`
	PriceImpactPercentage string   `json:"priceImpactPercentage"`
	EstimatedGas         string    `json:"estimatedGas"` // Compute units
	Instructions         string    `json:"tx"`           // Base64 encoded instruction data
	ComputeUnitLimit     string    `json:"computeUnitLimit,omitempty"`
	ComputeUnitPrice     string    `json:"computeUnitPrice,omitempty"`
}
