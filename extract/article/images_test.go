package article_test

import (
	"errors"
	dto "github.com/editorpost/article"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/extract/media"
	"github.com/stretchr/testify/assert"
	"strings"
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

![Alt text](http://storage.s3/image1.png)`,
			expectedErr: nil,
		},
		{
			name: "Multiple ImageDownloader",
			markdown: `# Sample Markdown

![Alt text](http://example.com/image1.png)

![Another alt text](http://example.com/image2.jpg)`,
			expected: `# Sample Markdown

![Alt text](http://storage.s3/image1.png)

![Another alt text](http://storage.s3/image2.jpg)`,
			expectedErr: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			a := dto.NewArticle()

			a.Markup = tt.markdown

			article.Images("test", a, &MockDownloadClaims{})
			assert.Equal(t, tt.expected, a.Markup)
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

	return media.Claim{
		Src: src,
		Dst: strings.Replace(src, "example.com", "storage.s3", 1),
	}, nil
}
