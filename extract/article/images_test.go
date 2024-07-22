package article_test

import (
	"errors"
	"github.com/editorpost/spider/extract/article"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarkdownImages(t *testing.T) {
	d := MockDownloadClaims{}

	tests := []struct {
		name        string
		markdown    string
		expected    string
		expectedErr error
	}{
		{
			name: "No ImageDownloader",
			markdown: `# Sample Markdown

This is a sample markdown without images.`,
			expected: `# Sample Markdown

This is a sample markdown without images.`,
			expectedErr: nil,
		},
		{
			name: "Single Image",
			markdown: `# Sample Markdown

![Alt text](http://example.com/image1.png)`,
			expected: `# Sample Markdown

![Alt text](downloaded_http://example.com/image1.png)`,
			expectedErr: nil,
		},
		{
			name: "Multiple ImageDownloader",
			markdown: `# Sample Markdown

![Alt text](http://example.com/image1.png)

![Another alt text](http://example.com/image2.jpg)`,
			expected: `# Sample Markdown

![Alt text](downloaded_http://example.com/image1.png)

![Another alt text](downloaded_http://example.com/image2.jpg)`,
			expectedErr: nil,
		},
		{
			name: "Image Download Failure",
			markdown: `# Sample Markdown

![Alt text](http://example.com/fail.jpg)`,
			expected:    "",
			expectedErr: errors.New("failed to download image"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := article.MarkdownImages(tt.markdown, d)
			if tt.expectedErr != nil {
				require.Error(t, err)
				assert.Equal(t, tt.expectedErr.Error(), err.Error())
			} else {
				require.NoError(t, err)
				assert.Equal(t, tt.expected, result)
			}
		})
	}
}

// Mock implementation of MediaClaims for testing
type MockDownloadClaims struct{}

func (m MockDownloadClaims) Add(src string) (string, error) {
	// Mock implementation: Just prepend "downloaded_" to the src URL
	if src == "http://example.com/fail.jpg" {
		return "", errors.New("failed to download image")
	}
	return "downloaded_" + src, nil
}
