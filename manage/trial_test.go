//go:build e2e
// +build e2e

package manage_test

import (
	"fmt"
	"github.com/editorpost/article"
	"github.com/editorpost/spider/collect/config"
	article2 "github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/manage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrial(t *testing.T) {

	args := &config.Args{
		StartURL:        "https://thailand-news.ru/news/turizm/",
		AllowedURL:      "https://thailand-news.ru/news/{turizm,puteshestviya}{any}",
		ExtractURL:      "https://thailand-news.ru/news/{turizm,puteshestviya}/{some}",
		ExtractLimit:    2,
		Depth:           1,
		ExtractSelector: ".node-article--full",
		ProxyEnabled:    false,
	}

	articles, tErr := manage.Trial(args, article2.Article)
	require.NoError(t, tErr)
	assert.NotNil(t, articles)

	for _, payload := range articles {
		a, err := article.NewArticleFromMap(payload.Data)
		assert.NoError(t, err)
		fmt.Printf("%s: %s\n", a.Published.Format("2006-01-02"), a.Title)
	}

	fmt.Printf("Extracted %d items\n", len(articles))
}
