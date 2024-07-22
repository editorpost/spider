package media_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/tester"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"sync"
	"testing"
)

type Loader struct {
	uploads sync.Map
}

func NewLoader() *Loader {
	return &Loader{}
}

// Upload fetches the media from the specified Endpoint and uploads it to the store.
func (dl *Loader) Download(src, dst string) error {
	dl.uploads.Store(dst, src)
	return nil
}

// Upload fetches the media from the specified Endpoint and uploads it to the store.
func (dl *Loader) Has(dst string) bool {
	_, ok := dl.uploads.Load(dst)
	return ok
}

func (dl *Loader) Len() int {
	count := 0
	dl.uploads.Range(func(_, _ any) bool {
		count++
		return true
	})
	return count
}

func TestNewMedia(t *testing.T) {

	loader := NewLoader()
	m := pipe.NewMedia("https://dst.com/static/media", loader)

	// create payload
	payload := tester.TestPayload(t, "../../tester/fixtures/news/article-1.html")
	require.NoError(t, m.Claims(payload))

	// empty
	require.NoError(t, m.Upload(payload))
	assert.Zero(t, loader.Len())

	// add claim
	src := gofakeit.URL()
	dst, err := payload.Download(src)
	require.NoError(t, err)
	assert.NotEmpty(t, dst)

	// claim added
	require.Equal(t, 1, payload.Claims().Len())
	// claim not uploaded yet
	require.Equal(t, 0, loader.Len())
	// upload
	require.NoError(t, m.Upload(payload))
	require.Equal(t, 1, loader.Len())
}

func GetHTML(t *testing.T) string {

	t.Helper()

	// open file `article_test.html` return as string
	f, err := os.Open("media_test.html")
	require.NoError(t, err)
	defer f.Close()

	// read file as a string
	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)
	require.NoError(t, err)

	return buf.String()
}

func GetDocument(t *testing.T) *colly.HTMLElement {

	t.Helper()

	// parse html
	query, err := goquery.NewDocumentFromReader(strings.NewReader(GetHTML(t)))
	require.NoError(t, err)

	ctx := &colly.Context{}
	resp := &colly.Response{
		Request: &colly.Request{
			Ctx: ctx,
		},
		Ctx: ctx,
	}

	var doc *colly.HTMLElement
	doc = colly.NewHTMLElementFromSelectionNode(resp, query.Selection, query.Nodes[0], 0)

	query.Find("html").Each(func(_ int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			if doc != nil {
				doc = colly.NewHTMLElementFromSelectionNode(resp, s, n, 0)
			}
		}
	})

	return doc
}
