//go:build e2e

package manage_test

import (
	"fmt"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/manage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSingleTurn_ExtractsDataFromURL(t *testing.T) {

	uri := "https://thailand-news.ru/news/puteshestviya/pkhuket-v-stile-vashego-otdykha/"
	selector := "article"

	payload, err := manage.Single(uri, selector, article.Article)
	require.NoError(t, err)
	assert.NotNil(t, payload)
	fmt.Print(payload.Data)
}
