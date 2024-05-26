package article_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/editorpost/spider/extract/article" // Замените на фактический путь к вашему пакету
)

func TestSocialProfileNormalize(t *testing.T) {
	sp := &article.SocialProfile{
		Platform: "  " + gofakeit.Word() + "  ",
		URL:      "  " + gofakeit.URL() + "  ",
	}

	sp.Normalize()

	assert.NotEmpty(t, sp.URL)
	assert.Equal(t, strings.TrimSpace(sp.Platform), sp.Platform)
	assert.Equal(t, strings.TrimSpace(sp.URL), sp.URL)
}

// TestSocialProfileConversions is a table-driven test for the SocialProfile struct.
// It verifies the conversion of map data to SocialProfile struct, the validation process, and handling of zero-value fields.
//
// Explanation of test cases:
// - Valid SocialProfile: Ensures that valid data is correctly converted into a SocialProfile struct without errors.
// - Invalid URL: Ensures that an invalid URL triggers a validation error.
// - Missing Required Fields: Ensures that missing mandatory fields trigger a validation error. Specifically, the 'url' field is required.
// - Zero Value Fields: Ensures that empty field values are handled correctly and trigger a validation error.
func TestSocialProfileConversions(t *testing.T) {
	tests := []struct {
		name                  string
		inputMap              map[string]any
		expectedSocialProfile *article.SocialProfile
		expectError           bool
	}{
		{
			name: "Valid SocialProfile",
			inputMap: map[string]any{
				"platform": "Twitter",
				"url":      "https://twitter.com/example",
			},
			expectedSocialProfile: &article.SocialProfile{
				Platform: "Twitter",
				URL:      "https://twitter.com/example",
			},
			expectError: false,
		},
		{
			name: "Invalid URL",
			inputMap: map[string]any{
				"platform": "Twitter",
				"url":      "invalid-url",
			},
			expectedSocialProfile: nil,
			expectError:           true,
		},
		{
			name: "Missing Required Fields",
			inputMap: map[string]any{
				"platform": "Twitter",
			},
			expectedSocialProfile: nil,
			expectError:           true,
		},
		{
			name: "Zero Value Fields",
			inputMap: map[string]any{
				"platform": "",
				"url":      "",
			},
			expectedSocialProfile: &article.SocialProfile{
				Platform: "",
				URL:      "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			profile, err := article.NewSocialProfileFromMap(tt.inputMap)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedSocialProfile, profile)
				assert.Equal(t, tt.inputMap, profile.Map())
			}
		})
	}
}
