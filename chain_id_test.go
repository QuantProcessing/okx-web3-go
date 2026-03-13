package okx

import (
	"encoding/json"
	"testing"
)

func TestChainIDType_UnmarshalString(t *testing.T) {
	input := `"8453"`
	var c ChainIDType
	if err := json.Unmarshal([]byte(input), &c); err != nil {
		t.Fatalf("Failed to unmarshal string chainId: %v", err)
	}
	if c.String() != "8453" {
		t.Errorf("Expected '8453', got '%s'", c.String())
	}
}

func TestChainIDType_UnmarshalNumber(t *testing.T) {
	input := `8453`
	var c ChainIDType
	if err := json.Unmarshal([]byte(input), &c); err != nil {
		t.Fatalf("Failed to unmarshal number chainId: %v", err)
	}
	if c.String() != "8453" {
		t.Errorf("Expected '8453', got '%s'", c.String())
	}
}

func TestChainIDType_UnmarshalInStruct(t *testing.T) {
	tests := []struct {
		name     string
		json     string
		expected string
	}{
		{
			name:     "string chainId",
			json:     `{"chainId":"1"}`,
			expected: "1",
		},
		{
			name:     "number chainId",
			json:     `{"chainId":56}`,
			expected: "56",
		},
		{
			name:     "large number chainId",
			json:     `{"chainId":42161}`,
			expected: "42161",
		},
		{
			name:     "solana chainId string",
			json:     `{"chainId":"501"}`,
			expected: "501",
		},
	}

	type testStruct struct {
		ChainID ChainIDType `json:"chainId"`
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var s testStruct
			if err := json.Unmarshal([]byte(tt.json), &s); err != nil {
				t.Fatalf("Failed to unmarshal: %v", err)
			}
			if s.ChainID.String() != tt.expected {
				t.Errorf("Expected '%s', got '%s'", tt.expected, s.ChainID.String())
			}
		})
	}
}

func TestChainIDType_String(t *testing.T) {
	c := ChainIDType("137")
	if c.String() != "137" {
		t.Errorf("Expected '137', got '%s'", c.String())
	}
}

func TestChainIDType_UnmarshalNull(t *testing.T) {
	input := `null`
	var c ChainIDType
	// null should not cause an error
	_ = json.Unmarshal([]byte(input), &c)
	// c should be empty string
	if c.String() != "" {
		t.Errorf("Expected empty string for null, got '%s'", c.String())
	}
}
