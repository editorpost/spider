package article_test

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/editorpost/spider/tester"
	"github.com/samber/lo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func HTMLToMarkdown(html string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
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
	assert.Equal(t, out, HTMLToMarkdown(in))
}

// extract specific image
func TestHTMLToMarkdownImageFromMany(t *testing.T) {
	in := tester.GetHTML(t, "../../tester/fixtures/cases/must_article_title.html")
	image := "![](/sites/default/files/storage/images/2016-20/rambutan-thaiskii-frukt.jpg)"
	assert.Contains(t, HTMLToMarkdown(in), image)
}
