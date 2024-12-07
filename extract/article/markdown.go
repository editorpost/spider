package article

import (
	md "github.com/JohannesKaufmann/html-to-markdown/v2"
	"github.com/JohannesKaufmann/html-to-markdown/v2/converter"
	"github.com/PuerkitoBio/goquery"
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
	linksToText(doc.Selection)

	// Return cleaned HTML
	out, _ := doc.Html()
	return out
}

// linksToText cleans links based on the requirements
func linksToText(selection *goquery.Selection) {
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
