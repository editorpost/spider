package media_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

func TestNewClaims(t *testing.T) {
	claims := media.NewClaims("http://example.com")
	assert.NotNil(t, claims)
}

// Test extract

func TestNewExtract(t *testing.T) {

	code := GetArticleHTML(t)
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(code))

	// check if url replaced in selection
	count := 0
	doc.Find("img").Each(func(i int, el *goquery.Selection) {
		if _, ok := el.Attr("src"); ok {
			count++
		}
	})

	require.NoError(t, err)
	sel := doc.Find("body")

	prefix := "http://example.com"
	claims := media.NewClaims(prefix)
	claims.ExtractAndReplace(sel)
	assert.Equal(t, count, claims.Len())

	// check if url replaced in selection
	sel.Find("img").Each(func(i int, el *goquery.Selection) {
		src, _ := el.Attr("src")
		assert.True(t, strings.HasPrefix(src, prefix))
	})

	// check if url replaced in document (pointers are the same)
	doc.Find("img").Each(func(i int, el *goquery.Selection) {
		src, _ := el.Attr("src")
		assert.True(t, strings.HasPrefix(src, prefix))
	})
}

func GetArticleHTML(t *testing.T) string {

	t.Helper()

	// open file `article_test.html` return as string
	f, err := os.Open("claims_test.html")
	require.NoError(t, err)
	defer f.Close()

	// read file as a string
	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)
	require.NoError(t, err)

	return buf.String()
}