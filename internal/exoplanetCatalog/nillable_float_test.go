package exoplanetCatalog

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNillableFloat_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		name     string
		jsonData string
		expected *float64
	}{
		{
			name:     "Valid float",
			jsonData: "123.45",
			expected: ptrFloat64(123.45),
		},
		{
			name:     "Valid integer treated as float",
			jsonData: "67",
			expected: ptrFloat64(67),
		},
		{
			name:     "Null value",
			jsonData: "null",
			expected: nil,
		},
		{
			name:     "Non-numeric string",
			jsonData: `"abc"`,
			expected: nil,
		},
		{
			name:     "Empty string",
			jsonData: `""`,
			expected: nil,
		},
		{
			name:     "Boolean true",
			jsonData: "true",
			expected: nil,
		},
		{
			name:     "Boolean false",
			jsonData: "false",
			expected: nil,
		},
		{
			name:     "Array",
			jsonData: "[1, 2, 3]",
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var nf nillableFloat
			err := json.Unmarshal([]byte(tt.jsonData), &nf)
			assert.Nil(t, err)
			assert.Equal(t, tt.expected, nf.Value)
		})
	}
}
