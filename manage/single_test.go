package manage_test

import (
	"fmt"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestSingleTurn_ExtractsDataFromURL(t *testing.T) {

	uri := "https://www.bbc.com/news/uk-england-gloucestershire-69055101"
	selector := "article"

	payload, err := manage.SingleTurn(uri, selector, extract.Article)
	require.NoError(t, err)
	assert.NotNil(t, payload)
	fmt.Print(payload.Data)
}
