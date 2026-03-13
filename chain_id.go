package okx

import "encoding/json"

// ChainIDType is a flexible type that can handle both string and number chainId from API
type ChainIDType string

// UnmarshalJSON implements custom unmarshaling to handle both string and number
func (c *ChainIDType) UnmarshalJSON(data []byte) error {
	// Try to unmarshal as string first
	var s string
	if err := json.Unmarshal(data, &s); err == nil {
		*c = ChainIDType(s)
		return nil
	}

	// If that fails, try as number
	var n json.Number
	if err := json.Unmarshal(data, &n); err == nil {
		*c = ChainIDType(n.String())
		return nil
	}

	return nil
}

// String returns the chainId as a string
func (c ChainIDType) String() string {
	return string(c)
}
