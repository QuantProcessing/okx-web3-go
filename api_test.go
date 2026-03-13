package okx

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestGetSupportedChains(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/supported/chain" {
			t.Errorf("Expected path '/api/v6/dex/aggregator/supported/chain', got '%s'", r.URL.Path)
		}
		if r.Method != "GET" {
			t.Errorf("Expected GET, got %s", r.Method)
		}

		resp := BaseResponse[SupportedChain]{
			Code: "0",
			Msg:  "",
			Data: []SupportedChain{
				{ChainID: "1", ChainName: "Ethereum", NativeTokenSymbol: "ETH"},
				{ChainID: "56", ChainName: "BNB Chain", NativeTokenSymbol: "BNB"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	chains, err := client.GetSupportedChains(context.Background())
	if err != nil {
		t.Fatalf("GetSupportedChains returned error: %v", err)
	}
	if len(chains) != 2 {
		t.Fatalf("Expected 2 chains, got %d", len(chains))
	}
	if chains[0].ChainName != "Ethereum" {
		t.Errorf("Expected first chain 'Ethereum', got '%s'", chains[0].ChainName)
	}
	if chains[1].NativeTokenSymbol != "BNB" {
		t.Errorf("Expected second chain native token 'BNB', got '%s'", chains[1].NativeTokenSymbol)
	}
}

func TestGetTokens(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/all-tokens" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("chainIndex") != "8453" {
			t.Errorf("Expected chainIndex=8453, got %s", r.URL.Query().Get("chainIndex"))
		}

		resp := BaseResponse[Token]{
			Code: "0",
			Data: []Token{
				{TokenName: "USD Coin", TokenSymbol: "USDC", TokenAddress: "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913", Decimals: "6"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	tokens, err := client.GetTokens(context.Background(), "8453")
	if err != nil {
		t.Fatalf("GetTokens returned error: %v", err)
	}
	if len(tokens) != 1 {
		t.Fatalf("Expected 1 token, got %d", len(tokens))
	}
	if tokens[0].TokenSymbol != "USDC" {
		t.Errorf("Expected token symbol 'USDC', got '%s'", tokens[0].TokenSymbol)
	}
}

func TestGetLiquiditySources(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/all-liquidity" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		if r.URL.Query().Get("chainIndex") != "1" {
			t.Errorf("Expected chainIndex=1, got %s", r.URL.Query().Get("chainIndex"))
		}

		resp := BaseResponse[LiquiditySource]{
			Code: "0",
			Data: []LiquiditySource{
				{DexName: "Uniswap V3", DexProtocol: "uniswap_v3"},
				{DexName: "Curve", DexProtocol: "curve"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	sources, err := client.GetLiquiditySources(context.Background(), "1")
	if err != nil {
		t.Fatalf("GetLiquiditySources returned error: %v", err)
	}
	if len(sources) != 2 {
		t.Fatalf("Expected 2 sources, got %d", len(sources))
	}
	if sources[0].DexName != "Uniswap V3" {
		t.Errorf("Expected 'Uniswap V3', got '%s'", sources[0].DexName)
	}
}

func TestGetQuote(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/quote" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("chainIndex") != "8453" {
			t.Errorf("Expected chainIndex=8453")
		}
		if q.Get("amount") != "100000000000000000" {
			t.Errorf("Expected amount=100000000000000000")
		}
		if q.Get("slippage") != "0.005" {
			t.Errorf("Expected slippage=0.005")
		}

		resp := BaseResponse[QuoteResponse]{
			Code: "0",
			Data: []QuoteResponse{
				{
					ChainID:         "8453",
					FromTokenAmount: "100000000000000000",
					ToTokenAmount:   "250000000",
					EstimatedGas:    "150000",
					FromToken:       TokenInfo{TokenSymbol: "ETH"},
					ToToken:         TokenInfo{TokenSymbol: "USDC"},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	quote, err := client.GetQuote(context.Background(), &QuoteRequest{
		ChainID:          "8453",
		FromTokenAddress: "0xEeeeeEeeeEeEeeEeEeEeeEEEeeeeEeeeeeeeEEeE",
		ToTokenAddress:   "0x833589fCD6eDb6E08f4c7C32D4f71b54bdA02913",
		Amount:           "100000000000000000",
		Slippage:         "0.005",
	})
	if err != nil {
		t.Fatalf("GetQuote returned error: %v", err)
	}
	if quote.ToTokenAmount != "250000000" {
		t.Errorf("Expected ToTokenAmount '250000000', got '%s'", quote.ToTokenAmount)
	}
	if quote.FromToken.TokenSymbol != "ETH" {
		t.Errorf("Expected FromToken.TokenSymbol 'ETH', got '%s'", quote.FromToken.TokenSymbol)
	}
}

func TestGetQuote_NoResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := BaseResponse[QuoteResponse]{Code: "0", Data: []QuoteResponse{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	_, err := client.GetQuote(context.Background(), &QuoteRequest{
		ChainID:          "8453",
		FromTokenAddress: "0x1",
		ToTokenAddress:   "0x2",
		Amount:           "1",
		Slippage:         "0.005",
	})
	if err == nil {
		t.Fatal("Expected error for empty results")
	}
	if !IsDexError(err) {
		t.Errorf("Expected DexError, got %T", err)
	}
}

func TestGetApproveTransaction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/approve-transaction" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		resp := BaseResponse[ApproveTransactionResponse]{
			Code: "0",
			Data: []ApproveTransactionResponse{
				{
					DexContractAddress:   "0xRouterAddress",
					TokenContractAddress: "0xTokenAddress",
					Data:                 "0xApproveCalldata",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	result, err := client.GetApproveTransaction(context.Background(), &ApproveTransactionRequest{
		ChainID:              "8453",
		TokenContractAddress: "0xTokenAddress",
		ApproveAmount:        "1000000",
	})
	if err != nil {
		t.Fatalf("GetApproveTransaction returned error: %v", err)
	}
	if result.DexContractAddress != "0xRouterAddress" {
		t.Errorf("Expected router '0xRouterAddress', got '%s'", result.DexContractAddress)
	}
	if result.Data != "0xApproveCalldata" {
		t.Errorf("Expected calldata '0xApproveCalldata', got '%s'", result.Data)
	}
}

func TestBuildSwapData(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/swap" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("userWalletAddress") != "0xMyWallet" {
			t.Errorf("Expected userWalletAddress=0xMyWallet")
		}

		resp := BaseResponse[SwapResponse]{
			Code: "0",
			Data: []SwapResponse{
				{
					RouterResult: RouterResult{
						FromTokenAmount: "100000000000000000",
						ToTokenAmount:   "250000000",
						ToToken:         TokenInfo{TokenSymbol: "USDC"},
					},
					Tx: TxData{
						From:  "0xMyWallet",
						To:    "0xRouter",
						Data:  "0xSwapCalldata",
						Value: "100000000000000000",
					},
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	swap, err := client.BuildSwapData(context.Background(), &SwapRequest{
		ChainID:           "8453",
		FromTokenAddress:  "0xEeee",
		ToTokenAddress:    "0xUSDC",
		Amount:            "100000000000000000",
		Slippage:          "0.005",
		UserWalletAddress: "0xMyWallet",
	})
	if err != nil {
		t.Fatalf("BuildSwapData returned error: %v", err)
	}
	if swap.Tx.From != "0xMyWallet" {
		t.Errorf("Expected Tx.From '0xMyWallet', got '%s'", swap.Tx.From)
	}
	if swap.RouterResult.ToToken.TokenSymbol != "USDC" {
		t.Errorf("Expected ToToken 'USDC', got '%s'", swap.RouterResult.ToToken.TokenSymbol)
	}
}

func TestGetSwapHistory(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/swap/history" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}

		resp := BaseResponse[SwapHistoryResponse]{
			Code: "0",
			Data: []SwapHistoryResponse{
				{TxHash: "0xTx1", Status: "success", FromTokenAmount: "100", ToTokenAmount: "250"},
				{TxHash: "0xTx2", Status: "pending", FromTokenAmount: "200", ToTokenAmount: "500"},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	history, err := client.GetSwapHistory(context.Background(), &SwapHistoryRequest{
		ChainID:           "8453",
		UserWalletAddress: "0xMyWallet",
		Page:              "1",
		Limit:             "20",
	})
	if err != nil {
		t.Fatalf("GetSwapHistory returned error: %v", err)
	}
	if len(history) != 2 {
		t.Fatalf("Expected 2 history items, got %d", len(history))
	}
	if history[0].Status != "success" {
		t.Errorf("Expected first TX status 'success', got '%s'", history[0].Status)
	}
}

func TestGetSolanaSwapInstruction(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v6/dex/aggregator/swap-instruction" {
			t.Errorf("Unexpected path: %s", r.URL.Path)
		}
		q := r.URL.Query()
		if q.Get("chainIndex") != "501" {
			t.Errorf("Expected chainIndex=501")
		}

		resp := BaseResponse[SolanaSwapResponse]{
			Code: "0",
			Data: []SolanaSwapResponse{
				{
					ChainID:          "501",
					FromTokenAmount:  "1000000000",
					ToTokenAmount:    "25000000",
					MinReceiveAmount: "24875000",
					Instructions:     "base64EncodedInstructions",
				},
			},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	result, err := client.GetSolanaSwapInstruction(context.Background(), &SolanaSwapRequest{
		ChainID:           "501",
		FromTokenAddress:  "So11111111111111111111111111111111111111112",
		ToTokenAddress:    "EPjFWdd5AufqSSqeM2qN1xzybapC8G4wEGGkZwyTDt1v",
		Amount:            "1000000000",
		Slippage:          "0.005",
		UserWalletAddress: "SolanaWallet",
	})
	if err != nil {
		t.Fatalf("GetSolanaSwapInstruction returned error: %v", err)
	}
	if result.MinReceiveAmount != "24875000" {
		t.Errorf("Expected MinReceiveAmount '24875000', got '%s'", result.MinReceiveAmount)
	}
	if result.Instructions != "base64EncodedInstructions" {
		t.Errorf("Expected Instructions 'base64EncodedInstructions', got '%s'", result.Instructions)
	}
}

func TestGetSolanaSwapInstruction_NoResults(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		resp := BaseResponse[SolanaSwapResponse]{Code: "0", Data: []SolanaSwapResponse{}}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	_, err := client.GetSolanaSwapInstruction(context.Background(), &SolanaSwapRequest{
		ChainID:           "501",
		FromTokenAddress:  "from",
		ToTokenAddress:    "to",
		Amount:            "1",
		Slippage:          "0.005",
		UserWalletAddress: "wallet",
	})
	if err == nil {
		t.Fatal("Expected error for empty results")
	}
}

func TestBuildSwapData_WithOptionalParams(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		q := r.URL.Query()
		if q.Get("referrerAddress") != "0xReferrer" {
			t.Errorf("Expected referrerAddress=0xReferrer, got %s", q.Get("referrerAddress"))
		}
		if q.Get("feePercent") != "0.01" {
			t.Errorf("Expected feePercent=0.01, got %s", q.Get("feePercent"))
		}

		resp := BaseResponse[SwapResponse]{
			Code: "0",
			Data: []SwapResponse{{Tx: TxData{Data: "0x"}}},
		}
		json.NewEncoder(w).Encode(resp)
	}))
	defer server.Close()

	client := NewClient()
	client.BaseURL = server.URL
	client.WithCredentials("key", "secret", "pass", "proj")

	_, err := client.BuildSwapData(context.Background(), &SwapRequest{
		ChainID:           "8453",
		FromTokenAddress:  "0xA",
		ToTokenAddress:    "0xB",
		Amount:            "1000",
		Slippage:          "0.005",
		UserWalletAddress: "0xWallet",
		ReferrerAddress:   "0xReferrer",
		FeePercent:        "0.01",
	})
	if err != nil {
		t.Fatalf("BuildSwapData returned error: %v", err)
	}
}
