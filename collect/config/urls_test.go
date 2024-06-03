package config_test

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCreateRegexpString(t *testing.T) {
	tc := []struct {
		name     string
		pattern  string
		expected string
		matches  []string
		invalid  []string
	}{
		{
			name:     "dir and some",
			pattern:  "https://example.com/articles/{dir}/{some}",
			expected: "^https:\\/\\/example\\.com\\/articles\\/([^/]+)\\/(.+)$",
			matches: []string{
				"https://example.com/articles/one/two",
				"https://example.com/articles/one/two/three",
			},
			invalid: []string{
				"https://example.com/articles/one/",
			},
		},
		{
			name:     "or",
			pattern:  "https://example.com/{one,two}{any}",
			expected: "^https:\\/\\/example\\.com\\/(one|two)(.*)$",
			matches: []string{
				"https://example.com/one",
				"https://example.com/one/",
				"https://example.com/one/two/three",
				"https://example.com/two",
			},
			invalid: []string{
				"https://example.com/three/",
			},
		},
		{
			name:     "or_spaced",
			pattern:  "https://example.com/{one, two}{any}",
			expected: "^https:\\/\\/example\\.com\\/(one|two)(.*)$",
			matches: []string{
				"https://example.com/one",
				"https://example.com/one/",
				"https://example.com/one/two/three",
				"https://example.com/two",
			},
			invalid: []string{
				"https://example.com/three/",
			},
		},
		{
			name:     "or_strict",
			pattern:  "https://example.com/{one,two}",
			expected: "^https:\\/\\/example\\.com\\/(one|two)$",
			matches: []string{
				"https://example.com/one",
				"https://example.com/two",
			},
			invalid: []string{
				"https://example.com/one/",
				"https://example.com/two/",
				"https://example.com/three",
			},
		},
		{
			name:     "dir and any",
			pattern:  "https://example.com/articles/{dir}/{any}",
			expected: "^https:\\/\\/example\\.com\\/articles\\/([^/]+)\\/(.*)$",
			matches: []string{
				"https://example.com/articles/one/",
				"https://example.com/articles/one/two",
				"https://example.com/articles/one/two/three",
			},
			invalid: []string{
				"https://example.com/articles/one", // This now matches due to .* being able to match an empty string
			},
		},
		{
			name:     "dir and num",
			pattern:  "https://example.com/articles/{dir}/{num}",
			expected: "^https:\\/\\/example\\.com\\/articles\\/([^/]+)\\/(\\d+)$",
			matches: []string{
				"https://example.com/articles/one/123",
			},
			invalid: []string{
				"https://example.com/articles/one/two",
			},
		},
	}

	for _, tt := range tc {
		t.Run(tt.name, func(t *testing.T) {
			actual := config.RegexPattern(tt.pattern)
			assert.Equal(t, tt.expected, actual)

			for _, match := range tt.matches {
				assert.Regexp(t, actual, match)
			}

			for _, invalid := range tt.invalid {
				assert.NotRegexp(t, actual, invalid)
			}
		})
	}
}
