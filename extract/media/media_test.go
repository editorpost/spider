package media_test

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/extract/payload"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/require"
	"io"
	"net/url"
	"os"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
)

type Loader struct {
	uploads sync.Map
}

func NewLoader() *Loader {
	return &Loader{}
}

// Upload fetches the media from the specified URL and uploads it to the store.
func (dl *Loader) Upload(src, dst string) (string, error) {
	dl.uploads.Store(dst, src)
	return dst, nil
}

func TestNewMedia(t *testing.T) {

	doc := GetDocument(t)

	loader := NewLoader()
	m := media.NewMedia("https://dst.com/static/media", "static/media", loader)

	p := &payload.Payload{
		Ctx:       context.Background(),
		Doc:       doc,
		Selection: doc.DOM,
		URL:       &url.URL{},
		Data:      map[string]any{},
	}

	// Claims extracts all images urls from `src` attribute in the document.
	require.NoError(t, m.Claims(p))

	// load claims from payload context
	claims := p.Ctx.Value(media.ClaimsCtxKey).(*media.Claims)
	require.NotZero(t, len(claims.All()))

	// Upload requested media from claims
	upload, err := m.Upload()
	require.NoError(t, err)
	require.NoError(t, upload(p))

	count := atomic.Int32{}
	loader.uploads.Range(func(key, value any) bool {
		src, ok := value.(string)
		require.True(t, ok)
		require.NotEmpty(t, src)
		count.Add(1)
		return true
	})

	// no claims requested
	require.Zero(t, count.Load())

	// request claims to upload
	claims.Request(claims.All()[0].Dst)

	// Upload requested media from claims
	require.NoError(t, upload(p))

	count = atomic.Int32{}
	loader.uploads.Range(func(key, value any) bool {
		src, ok := value.(string)
		require.True(t, ok)
		require.NotEmpty(t, src)
		count.Add(1)
		return true
	})

	// 1 claim requested and uploaded
	require.Zero(t, count.Load())
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
