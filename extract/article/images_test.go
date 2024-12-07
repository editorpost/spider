package article_test

import (
	"github.com/editorpost/spider/extract/article"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestMarkdownSourceUrls(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "Simple URL",
			input:    `![alt](http://example.com/image.png)`,
			expected: []string{"http://example.com/image.png"},
		},
		{
			name:     "URL with Title",
			input:    `![alt](http://example.com/image.png "title")`,
			expected: []string{"http://example.com/image.png"},
		},
		{
			name:     "URL with Single-Quoted Title",
			input:    `![alt](http://example.com/image.png 'title')`,
			expected: []string{"http://example.com/image.png"},
		},
		{
			name:     "Multiple Images",
			input:    `![alt1](http://example.com/img1.png) ![alt2](http://example.com/img2.png "title2")`,
			expected: []string{"http://example.com/img1.png", "http://example.com/img2.png"},
		},
		{
			name:     "Complex Title",
			input:    `![alt](http://example.com/image.png "complex title with spaces")`,
			expected: []string{"http://example.com/image.png"},
		},
		{
			name:     "No Images",
			input:    `No markdown image tags here!`,
			expected: nil,
		},
		{
			name:     "Malformed Image Tag",
			input:    `![alt](http://example.com/image.png "title)`,
			expected: nil,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			actual := article.MarkdownSourceUrls(tt.input)
			assert.Equal(t, tt.expected, actual)
		})
	}
}
