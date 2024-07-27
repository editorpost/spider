package article_test

import (
	"errors"
	"github.com/editorpost/spider/extract/media"
	"testing"
)

func TestMarkdownImages(t *testing.T) {

	_ = MockDownloadClaims{}

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
			name: "Image Enabled Failure",
			markdown: `# Sample Markdown

![Alt text](http://example.com/fail.jpg)`,
			expected:    "",
			expectedErr: errors.New("failed to download image"),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// todo
		})
	}
}

// Mock implementation of MediaClaims for testing
type MockDownloadClaims struct{}

func (m MockDownloadClaims) Add(payloadID string, src string) (media.Claim, error) {
	// Mock implementation: Just prepend "downloaded_" to the src URL
	if src == "http://example.com/fail.jpg" {
		return media.Claim{}, errors.New("failed to download image")
	}
	return media.Claim{Dst: "downloaded_" + src}, nil
}
