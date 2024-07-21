package tester

import (
	"context"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/stretchr/testify/require"
	"net/url"
	"testing"
)

func TestPayload(t *testing.T, path string) *pipe.Payload {

	doc := GetDocument(t, path)
	id, err := pipe.Hash(gofakeit.URL())
	require.NoError(t, err)

	pay := &pipe.Payload{
		ID:        id,
		Ctx:       context.Background(),
		Doc:       doc,
		Selection: doc.DOM,
		URL:       &url.URL{},
		Data:      map[string]any{},
	}

	require.NoError(t, article.Article(pay))

	return pay
}
