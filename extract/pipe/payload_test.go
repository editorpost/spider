package pipe_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/collect/events"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestNewPayload(t *testing.T) {

	// get document
	doc := tester.GetDocument(t, "../../tester/fixtures/cases/must_article_title.html")
	doc.Request.URL, _ = url.Parse(gofakeit.URL())

	// get selections
	selections := events.Selections(doc, ".node-article--full", nil)
	require.Greater(t, len(selections), 0)

	// create payload
	pay, err := pipe.NewPayload(doc, selections[0])
	require.NoError(t, err)

	// data is not empty
	assert.NoError(t, article.Article(pay))
	assert.Greater(t, len(pay.Data), 0)
}
