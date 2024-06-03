package config_test

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestNormalizeURLs(t *testing.T) {

	tests := []struct {
		name  string
		start string
		allow string
		err   bool
	}{
		{
			name:  "empty start url",
			start: "",
			err:   true,
		},
		{
			name:  "invalid start url",
			start: "http://",
			err:   true,
		},
		{
			name:  "valid start url",
			start: "http://example.com",
		},
		{
			name:  "empty allowed url",
			start: "http://example.com",
			allow: "",
		},
		{
			name:  "valid allowed url",
			start: "http://example.com",
			allow: ".*",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			args := &config.Args{
				StartURL:   tt.start,
				AllowedURL: tt.allow,
			}

			err := args.NormalizeURLs()
			if tt.err {
				assert.Error(t, err)
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

func TestNormalizeUserAgent(t *testing.T) {

	tests := []struct {
		name string
		ua   string
	}{
		{
			name: "empty user agent",
			ua:   "",
		},
		{
			name: "valid user agent",
			ua:   "Mozilla/5.0",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			args := &config.Args{
				UserAgent: tt.ua,
			}

			err := args.NormalizeUserAgent()
			require.NoError(t, err)
		})
	}
}

func TestRootUrl(t *testing.T) {

	tests := []struct {
		name string
		url  string
		root string
	}{
		{
			name: "root url",
			url:  "https://example.com",
			root: "https://example.com",
		},
		{
			name: "root url with path",
			url:  "https://example.com/articles/1234/5678",
			root: "https://example.com",
		},
	}

	for _, tt := range tests {

		t.Run(tt.name, func(t *testing.T) {

			u, err := url.ParseRequestURI(tt.url)
			require.NoError(t, err)

			root := config.RootUrl(u)
			require.Equal(t, tt.root, root)
		})
	}
}
