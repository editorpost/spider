package tester

import (
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestPayload(t *testing.T, path string) *pipe.Payload {

	doc := GetDocument(t, path)
	uri := gofakeit.URL()

	doc.Request.URL, _ = url.Parse(uri)

	pay, err := pipe.NewPayload(doc, doc.DOM)
	require.NoError(t, err)
	require.NoError(t, article.Article(pay))

	return pay
}

func TestPayloadWithURI(t *testing.T, path, uri string) *pipe.Payload {

	doc := GetDocument(t, path)

	doc.Request.URL, _ = url.Parse(uri)

	pay, err := pipe.NewPayload(doc, doc.DOM)
	require.NoError(t, err)
	require.NoError(t, article.Article(pay))

	return pay
}
