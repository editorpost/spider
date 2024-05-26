package article_test

import (
	"github.com/editorpost/spider/extract/article"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// TestVideoConversions is a table-driven test for the Video struct.
// It verifies the conversion of map data to Video struct, the validation process, and handling of zero-value fields.
//
// Explanation of test cases:
// - Valid Video: Ensures that valid data is correctly converted into a Video struct without errors.
// - Invalid URL: Ensures that an invalid URL triggers a validation error.
// - Missing Required Fields: Ensures that missing mandatory fields trigger a validation error. Specifically, the 'url' field is required.
// - Zero Value Fields: Ensures that empty field values are handled correctly and trigger a validation error.
func TestVideoConversions(t *testing.T) {
	tests := []struct {
		name          string
		inputMap      map[string]any
		expectedVideo *article.Video
		expectError   bool
	}{
		{
			name: "Valid Video",
			inputMap: map[string]any{
				"url":        "https://example.com/video.mp4",
				"embed_code": "<iframe src='https://example.com/video'></iframe>",
				"caption":    "An example caption",
			},
			expectedVideo: &article.Video{
				URL:       "https://example.com/video.mp4",
				EmbedCode: "<iframe src='https://example.com/video'></iframe>",
				Caption:   "An example caption",
			},
			expectError: false,
		},
		{
			name: "Invalid URL",
			inputMap: map[string]any{
				"url":        "invalid-url",
				"embed_code": "<iframe src='https://example.com/video'></iframe>",
				"caption":    "An example caption",
			},
			expectedVideo: nil,
			expectError:   true,
		},
		{
			name: "Missing Required Fields",
			inputMap: map[string]any{
				"caption": "Some caption, but no URL",
			},
			expectedVideo: nil,
			expectError:   true,
		},
		{
			name: "Zero Value Fields",
			inputMap: map[string]any{
				"url":        "",
				"embed_code": "",
				"caption":    "",
			},
			expectedVideo: &article.Video{
				URL:       "",
				EmbedCode: "",
				Caption:   "",
			},
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			vid, err := article.NewVideoFromMap(tt.inputMap)
			if tt.expectError {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expectedVideo, vid)
				assert.Equal(t, tt.inputMap, vid.Map())
			}
		})
	}
}
