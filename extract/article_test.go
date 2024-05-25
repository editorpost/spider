package extract

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"strings"
	"testing"
)

func TestArticle(t *testing.T) {

	htmlStr, err := LoadHTMLFromFile("article_test.html")
	if err != nil {
		t.Fatalf("failed to load HTML from file: %v", err)
	}

	payload, err := CreatePayload(htmlStr, "https://example.com")
	if err != nil {
		t.Fatalf("failed to create payload: %v", err)
	}

	err = Article(payload)
	if err != nil {
		t.Fatalf("Article function returned an error: %v", err)
	}

	if payload.Data["entity__title"] != "Test Title" {
		t.Errorf("expected title 'Test Title', got '%v'", payload.Data["entity__title"])
	}

	if payload.Data["entity__byline"] != "Test Byline" {
		t.Errorf("expected byline 'Test Byline', got '%v'", payload.Data["entity__byline"])
	}

	if payload.Data["entity__content"] == nil || !strings.Contains(payload.Data["entity__content"].(string), "This is the content of the article.") {
		t.Errorf("expected content to contain 'This is the content of the article.', got '%v'", payload.Data["entity__content"])
	}

	if strings.Contains(payload.Data["entity__content"].(string), "Ad content") {
		t.Errorf("expected content to not contain 'Ad content', but it does")
	}

	if strings.Contains(payload.Data["entity__content"].(string), "alert('test')") {
		t.Errorf("expected content to not contain 'alert('test')', but it does")
	}

	if strings.Contains(payload.Data["entity__content"].(string), "<iframe") {
		t.Errorf("expected content to not contain '<iframe', but it does")
	}
}

func TestArticles(t *testing.T) {
	tests := []struct {
		name        string
		htmlFile    string
		url         string
		expected    map[string]any
		shouldError bool
	}{
		{
			name:     "Valid Article",
			htmlFile: "article_test.html",
			url:      "https://thailand-news.ru",
			expected: map[string]any{
				"entity__title":  "Пхукет в стиле вашего отдыха",
				"entity__byline": "",
			},
			shouldError: false,
		},
		{
			name:        "Invalid URL",
			htmlFile:    "article_test.html",
			url:         "://example.com",
			expected:    nil,
			shouldError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			htmlStr, err := LoadHTMLFromFile(tt.htmlFile)
			require.NoError(t, err, "failed to load HTML from file")

			payload, err := CreatePayload(htmlStr, tt.url)
			if tt.shouldError {
				require.Error(t, err)
				return
			}

			require.NoError(t, err, "failed to create payload")

			err = Article(payload)
			require.NoError(t, err, "Article function returned an error")

			for key, expectedValue := range tt.expected {
				assert.Contains(t, payload.Data, key)
				if key == "entity__content" {
					assert.Contains(t, payload.Data[key].(string), expectedValue.(string))
				} else {
					assert.Equal(t, expectedValue, payload.Data[key])
				}
			}

			assert.NotContains(t, payload.Data["entity__content"].(string), "Ad content")
			assert.NotContains(t, payload.Data["entity__content"].(string), "alert('test')")
			assert.NotContains(t, payload.Data["entity__content"].(string), "<iframe")
		})
	}
}
