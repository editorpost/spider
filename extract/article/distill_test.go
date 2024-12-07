package article_test

import (
	"github.com/go-shiori/dom"
	distiller "github.com/markusmobius/go-domdistiller"
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestDistiller(t *testing.T) {

	markup := GetArticleHTML(t)

	// Run distiller
	result, err := distiller.ApplyForReader(strings.NewReader(markup), &distiller.Options{
		OriginalURL: GetArticleURL(t),
	})
	if err != nil {
		panic(err)
	}

	rawHTML := dom.OuterHTML(result.Node)
	assert.NotEmpty(t, rawHTML)
}
