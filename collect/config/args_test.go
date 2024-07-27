package config_test

import (
	"errors"
	"github.com/editorpost/spider/collect/config"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
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

			args := &config.Config{
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

			args := &config.Config{
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

type ConfigTestSuite struct {
	suite.Suite
}

func (suite *ConfigTestSuite) TestNormalize() {
	tests := []struct {
		args     config.Config
		expected error
	}{
		{
			args: config.Config{
				StartURL: "https://example.com",
			},
			expected: nil,
		},
		{
			args: config.Config{
				StartURL: "",
			},
			expected: errors.New("start url is required"),
		},
		{
			args: config.Config{
				StartURL: "invalid-url",
			},
			expected: errors.New("start url is invalid: parse \"invalid-url\": invalid URI for request"),
		},
		{
			args: config.Config{
				StartURL: "https://",
			},
			expected: errors.New("start url host is invalid, add domain name"),
		},
	}

	for _, tt := range tests {
		err := tt.args.Normalize()
		if tt.expected == nil {
			assert.NoError(suite.T(), err)
		} else {
			assert.EqualError(suite.T(), err, tt.expected.Error())
		}
	}
}

func (suite *ConfigTestSuite) TestNormalizeURLs() {
	tests := []struct {
		args     config.Config
		expected error
	}{
		{
			args: config.Config{
				StartURL: "https://example.com",
			},
			expected: nil,
		},
		{
			args: config.Config{
				StartURL: "",
			},
			expected: errors.New("start url is required"),
		},
		{
			args: config.Config{
				StartURL: "invalid-url",
			},
			expected: errors.New("start url is invalid: parse \"invalid-url\": invalid URI for request"),
		},
		{
			args: config.Config{
				StartURL: "https://",
			},
			expected: errors.New("start url host is invalid, add domain name"),
		},
	}

	for _, tt := range tests {
		err := tt.args.NormalizeURLs()
		if tt.expected == nil {
			assert.NoError(suite.T(), err)
		} else {
			assert.EqualError(suite.T(), err, tt.expected.Error())
		}
	}
}

func (suite *ConfigTestSuite) TestNormalizeUserAgent() {
	tests := []struct {
		args     config.Config
		expected string
	}{
		{
			args: config.Config{
				UserAgent: "Mozilla/5.0",
			},
			expected: "Mozilla/5.0",
		},
		{
			args: config.Config{
				UserAgent: "",
			},
			expected: config.DefaultUserAgent,
		},
	}

	for _, tt := range tests {
		err := tt.args.NormalizeUserAgent()
		assert.NoError(suite.T(), err)
		assert.Equal(suite.T(), tt.expected, tt.args.UserAgent)
	}
}

func (suite *ConfigTestSuite) TestNormalizeExtractSelector() {
	tests := []struct {
		args     config.Config
		expected string
	}{
		{
			args: config.Config{
				ExtractSelector: " article ",
			},
			expected: "article",
		},
		{
			args: config.Config{
				ExtractSelector: "",
			},
			expected: "html",
		},
	}

	for _, tt := range tests {
		tt.args.NormalizeExtractSelector()
		assert.Equal(suite.T(), tt.expected, tt.args.ExtractSelector)
	}
}

func (suite *ConfigTestSuite) TestRootUrl() {
	tests := []struct {
		url      string
		expected string
	}{
		{
			url:      "https://example.com/articles/1234/5678",
			expected: "https://example.com",
		},
		{
			url:      "https://example.com",
			expected: "https://example.com",
		},
	}

	for _, tt := range tests {
		u, err := url.Parse(tt.url)
		assert.NoError(suite.T(), err)
		got := config.RootUrl(u)
		assert.Equal(suite.T(), tt.expected, got)
	}
}

func TestConfigTestSuite(t *testing.T) {
	suite.Run(t, new(ConfigTestSuite))
}
