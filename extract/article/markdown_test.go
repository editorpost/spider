package article_test

import (
	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/editorpost/spider/collect/events"
	"github.com/editorpost/spider/extract/article"
	"github.com/editorpost/spider/tester"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func HTMLToMarkdown(html string) string {
	markdown, err := md.ConvertString(html)
	return lo.Ternary(err == nil, markdown, "")
}

// just a simple test
func TestHTMLToMarkdown(t *testing.T) {
	in := "<h1>Hello, World!</h1>"
	out := "# Hello, World!"
	assert.Equal(t, out, HTMLToMarkdown(in))
}

// image conversion test
func TestHTMLToMarkdownImage(t *testing.T) {
	in := `<img src="https://example.com/image.jpg" alt="example image" />`
	out := "![example image](https://example.com/image.jpg)"
	mark, err := article.HTMLToStripMarkdown(in, "")
	require.NoError(t, err)
	assert.Equal(t, out, mark)
}

// extract specific image
func TestHTMLToMarkdownImageFromMany(t *testing.T) {

	doc := tester.GetDocument(t, "../../tester/fixtures/cases/must_article_title.html")
	selections := events.Selections(doc, ".node-article--full", nil)
	require.Greater(t, len(selections), 0)

	in, err := selections[0].Html()
	require.NoError(t, err)

	image := "![](/sites/default/files/storage/images/2016-20/rambutan-thaiskii-frukt.jpg)"
	mark, err := article.HTMLToStripMarkdown(in, "")
	require.NoError(t, err)
	assert.Contains(t, mark, image)
}

// ensure links are removed from markdown
func TestHTMLToMarkdownStripLinks(t *testing.T) {

	// markdown has a link
	in := `<a href="https://example.com">Example</a>`
	out := "Example"

	mark, err := article.HTMLToStripMarkdown(in, "")
	require.NoError(t, err)

	// markdown has no link
	assert.Equal(t, out, mark)
}

// ensure links are removed from markdown
func TestHTMLToMarkdownRemoveUrlLinks(t *testing.T) {

	// markdown has a link
	in := `<a href="https://example.com">https://example.com</a>`
	out := ""

	mark, err := article.HTMLToStripMarkdown(in, "")
	require.NoError(t, err)

	// markdown has no link
	assert.Equal(t, out, mark)
}
