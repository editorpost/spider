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
	in := `<img src="https://example.com/image.jpg" width="1200" height="800" alt="Image Alt" title="Image Title" class="some-css-class es" style="margin-bottom: 10px;" />`
	out := "![Image Alt](https://example.com/image.jpg \"Image Title\")"
	mark, err := article.HTMLToMarkdown(in, "")
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

func TestArticle_ReadabilityMainImage(t *testing.T) {

	// @note readability might remove the main image from the content
	// @see article.readabilityArticle()

	// get document
	pay := tester.TestPayload(t, "../../tester/fixtures/cases/must_article_image.html")

	// data is not empty
	assert.NoError(t, article.Article(pay))
	assert.Greater(t, len(pay.Data), 0)

	// the image must be in the markdown
	image := "/sites/default/files/storage/images/2016-20/rambutan-thaiskii-frukt.jpg)"

	assert.Contains(t, pay.Data["markup"], image)
}
