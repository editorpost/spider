package tester

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/require"
	"io"
	"os"
	"strings"
	"testing"
)

func GetHTML(t *testing.T, path string) string {

	t.Helper()

	f, err := os.Open(path)
	require.NoError(t, err)
	defer f.Close()

	// read file as a string
	buf := new(strings.Builder)
	_, err = io.Copy(buf, f)
	require.NoError(t, err)

	return buf.String()
}

func GetDocument(t *testing.T, path string) *colly.HTMLElement {
	t.Helper()
	return GetDocumentSelections(t, path, "html")
}

func GetDocumentSelections(t *testing.T, path string, selector string) *colly.HTMLElement {

	t.Helper()

	// parse html
	query, err := goquery.NewDocumentFromReader(strings.NewReader(GetHTML(t, path)))
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

	selections := query.Find(selector)

	// todo it must be refactored, there is no need to loop through
	selections.Each(func(_ int, s *goquery.Selection) {
		for _, n := range s.Nodes {
			if doc != nil {
				doc = colly.NewHTMLElementFromSelectionNode(resp, s, n, 0)
			}
		}
	})

	return doc
}
