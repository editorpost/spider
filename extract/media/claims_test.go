package media_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/spider/extract/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log/slog"
	"net/url"
	"os"
	"strings"
	"testing"
)

// Claim claim to get src media and save it to dst.
type UploadClaim struct {
	// Src is the URL of the media to download.
	Src string
	// Dst is the path to save the downloaded media.
	Dst string
	// Requested is true if the media is requested to download.
	Requested bool
	// Done is true if the media is downloaded to destination.
	Uploaded bool
}

type UploadClaims struct {
	// dstURL is a prefix of public path of the replaced media url.
	dstURL string
	// claims keyed with source url
	claims map[string]UploadClaim
}

// NewClaims creates a new Claim for each image and replace src path in document and selection.
// Replacement path from media.Filename. Replaces src url in selection.
func NewClaims(uri string) *UploadClaims {

	claims := &UploadClaims{
		dstURL: uri,
		claims: make(map[string]UploadClaim),
	}

	return claims
}

// Extract Claim for each img tag and replace src path in selection.
func (list *UploadClaims) Extract(selection *goquery.Selection) {

	selection.Find("img").Each(func(i int, el *goquery.Selection) {

		// has src
		src, exists := el.Attr("src")
		if !exists {
			return
		}

		// already claimed
		if _, exists = list.claims[src]; exists {
			return
		}

		// already replaced
		if strings.HasPrefix(src, list.dstURL) {
			return
		}

		// filename as src url hash
		filename, err := media.Filename(src)
		if err != nil {
			slog.Error("failed to hash filename", slog.String("src", src), slog.String("err", err.Error()))
			return
		}

		// full url
		dst, err := url.JoinPath(list.dstURL, filename)
		if err != nil {
			slog.Error("failed to join url", slog.String("dst", list.dstURL), slog.String("filename", filename), slog.String("err", err.Error()))
			return
		}

		// replace url in selection
		el.SetAttr("src", dst)

		// add claim
		list.Add(UploadClaim{
			Src: src,
			Dst: dst,
		})
	})
}

func (list *UploadClaims) Add(c UploadClaim) *UploadClaims {
	list.claims[c.Src] = c
	return list
}

func (list *UploadClaims) Len() int {
	return len(list.claims)
}

func TestNewClaims(t *testing.T) {
	claims := NewClaims("http://example.com")
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
	claims := NewClaims(prefix)
	claims.Extract(sel)
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
