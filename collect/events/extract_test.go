package events_test

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/collect/events"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestSelections(t *testing.T) {

	selector := "html"
	doc := tester.GetDocumentSelections(t, "../../tester/fixtures/cases/must_article_title.html", selector)

	// parse DOM and get elements by selector
	selections := events.Selections(doc, "html", nil)

	// html is only one element
	require.Len(t, selections, 1)

	// h1 is only one element too
	selections = events.Selections(doc, "h1", nil)
	require.Len(t, selections, 1)

	// articles are 27 elements
	selections = events.Selections(doc, "article", nil)
	require.Len(t, selections, 27)

	// but only 1 element with class ".node-article--full"
	selections = events.Selections(doc, ".node-article--full", nil)
	require.Len(t, selections, 1)

	// let's get article from the first selection
	selection := selections[0]
	doc.Request.URL, _ = url.Parse(gofakeit.URL())

	pay, err := pipe.NewPayload(doc, selection)
	require.NoError(t, err)
	require.NoError(t, article.Article(pay))
}

func TestSelectionsArticle(t *testing.T) {

	selector := "body"
	doc := tester.GetDocumentSelections(t, "../../tester/fixtures/cases/must_article_title.html", selector)
	uri := gofakeit.URL()

	doc.Request.URL, _ = url.Parse(uri)

	pay, err := pipe.NewPayload(doc, doc.DOM)
	require.NoError(t, err)
	require.NoError(t, article.Article(pay))
}
