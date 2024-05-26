package article_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/editorpost/spider/extract/article"
)

// TestQuoteConversions is a table-driven test for the Quote struct.
// It verifies the conversion of map data to Quote struct, the validation process, and handling of zero-value fields.
//
// Explanation of test cases:
// - Valid Quote: Ensures that valid data is correctly converted into a Quote struct without errors.
// - Invalid Source URL: Ensures that an invalid source URL triggers a validation error.
// - Missing Required Fields: Ensures that missing mandatory fields trigger a validation error. Specifically, 'text', 'author', 'source', and 'platform' fields are required.
// - Zero Value Fields: Ensures that empty field values are handled correctly and trigger a validation error.
func TestQuoteConversions(t *testing.T) {
	tests := []struct {
		name          string
		inputMap      map[string]any
		expectedQuote *article.Quote
		expectError   bool
	}{
		{
			name: "Valid Quote",
			inputMap: map[string]any{
				"text":     "This is a quote",
				"author":   "John Doe",
				"source":   "https://example.com",
				"platform": "Twitter",
			},
			expectedQuote: &article.Quote{
				Text:     "This is a quote",
				Author:   "John Doe",
				Source:   "https://example.com",
				Platform: "Twitter",
			},
			expectError: false,
		},
		{
			name: "Invalid Source URL",
			inputMap: map[string]any{
				"text":     "This is a quote",
				"author":   "John Doe",
				"source":   "invalid-url",
				"platform": "Twitter",
			},
			expectedQuote: nil,
			expectError:   true,
		},
		{
			name: "Missing Required Fields",
			inputMap: map[string]any{
				"text":   "This is a quote",
				"author": "John Doe",
			},
			expectedQuote: nil,
			expectError:   true,
		},
		{
			name: "Zero Value Fields",
			inputMap: map[string]any{
				"text":     "",
				"author":   "",
				"source":   "",
				"platform": "",
			},
			expectedQuote: &article.Quote{
				Text:     "",
				Author:   "",
				Source:   "",
				Platform: "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			quote, err := article.NewQuoteFromMap(tt.inputMap)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedQuote, quote)
				assert.Equal(t, tt.inputMap, quote.Map())
			}
		})
	}
}
