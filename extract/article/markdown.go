package article

import (
	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/PuerkitoBio/goquery"
	"net/url"
	"strings"
)

// HTMLToMarkdown converts HTML.
func HTMLToMarkdown(html, domain string) (string, error) {

	if domain == "" {
		return md.ConvertString(html)
	}

	return md.ConvertString(html,
		// replacing relative URLs with absolute
		converter.WithDomain(domain),
	)
}

// HTMLToStripMarkdown converts HTML to Markdown and cleans unwanted links
func HTMLToStripMarkdown(html, domain string) (string, error) {
	html = removeLinksFromHTML(html)
	return HTMLToMarkdown(html, domain)
}

// removeLinksFromHTML normalizes the HTML (e.g., removes links containing "http")
func removeLinksFromHTML(html string) string {
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(html))
	if err != nil {
		return html // Fallback to raw HTML if parsing fails
	}

	// Normalize links
	replaceLinks(doc.Selection)

	// Return cleaned HTML
	out, _ := doc.Html()
	return out
}

// replaceLinks cleans links based on the requirements
func replaceLinks(selection *goquery.Selection) {
	selection.Find("a").Each(func(i int, s *goquery.Selection) {

		text := s.Text()

		// Remove links where text contains "http"
		if isUriLikeString(text) {
			s.Remove()
		} else {
			// Replace links with their text content
			s.ReplaceWithHtml(text)
		}
	})
}

// isUriLikeString checks if a string is a URI
func isUriLikeString(s string) bool {

	// has schema in it
	if strings.Contains(s, "://") {
		return true
	}

	return false
}

// QueryToMarkup extracts HTML from a goquery selection and converts it to cleaned Markdown
func QueryToMarkup(selection *goquery.Selection, uri *url.URL) (string, error) {
	// Remove h1 from the selection
	selection.Find("h1").Remove()

	// Get the HTML from the selection
	html, err := selection.Html()
	if err != nil {
		return "", err
	}

	// Convert the HTML to Markdown
	return HTMLToStripMarkdown(html, uri.String())
}
