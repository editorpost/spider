package article

import (
	"fmt"
	dto "github.com/editorpost/article"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/go-shiori/go-readability"
	"github.com/samber/lo"
	"log/slog"
	"net/url"
	"strings"
)

// Article extracts the dto from the HTML
// and sets the dto fields to the payload
func Article(payload *pipe.Payload) error {

	art, err := ArticleFromPayload(payload)
	if err != nil {
		slog.Warn("failed to extract dto", slog.String("err", err.Error()))
		return err
	}

	// download claims
	if claims := media.ClaimsFromContext(payload.Ctx); claims != nil {
		// set article images, replace links in markdown
		Images(payload.ID, art, claims)
	}

	// set the dto to the payload
	for k, v := range art.Map() {
		payload.Data[k] = v
	}

	slog.Debug("extract success", slog.String("title", art.Title))

	return nil
}

// ArticleFromPayload extracts Article
func ArticleFromPayload(payload *pipe.Payload) (a *dto.Article, err error) {

	a = dto.NewArticle()
	a.SourceURL = payload.URL.String()

	// get the selection
	sel := payload.Selection.Clone()
	// remove title from content
	sel.Find("h1").Remove()
	// remove links, keep anchor text
	linksToText(sel)

	// get full DOM, including meta, scripts, etc.
	pageHTML, err := payload.Doc.DOM.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTML from document: %w", err)
	}

	// fill data from readability
	// note article.Markup is still HTML here, not markdown
	// get full DOM, including meta, scripts, etc.
	if err = articleMetadata(pageHTML, payload.URL, a); err != nil {
		return nil, fmt.Errorf("failed to get article metadata: %w", err)
	}

	// get selection html content
	contentHTML, err := sel.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTML from selection: %w", err)
	}

	// article html and text content from selection
	// @todo: avoid using readability twice
	if err = articleContent(contentHTML, payload.URL, a); err != nil {
		return nil, fmt.Errorf("failed to get article content: %w", err)
	}

	// html to markdown
	// article.Markup is now converted to markdown
	if a.Markup, err = HTMLToMarkdown(a.Markup, payload.URL); err != nil {
		return nil, fmt.Errorf("failed to convert HTML to markdown: %w", err)
	}

	// nil dto if it's invalid
	if err = a.Normalize(); err != nil {
		return nil, err
	}

	return a, nil
}

// HostUrl returns the host URL without path
func HostUrl(base *url.URL) string {
	return base.Scheme + "://" + base.Host
}

func articleMetadata(html string, addr *url.URL, a *dto.Article) error {

	// readability: title, summary, text, markup, html, language, summary
	read, err := readability.FromReader(strings.NewReader(html), addr)
	if err != nil {
		return nil
	}

	a.Title = read.Title

	a.Language = read.Language
	a.Summary = lo.Ternary(len(read.Byline) > 0, read.Byline, read.Excerpt)

	// fallback fields applied only if the fields are empty
	a.Title = lo.Ternary(a.Title == "", read.Title, a.Title)
	a.Summary = lo.Ternary(a.Summary == "", read.Excerpt, a.Summary)
	a.Language = lo.Ternary(a.Language == "", read.Language, a.Language)

	if read.PublishedTime != nil {
		a.Published = *read.PublishedTime
	}

	if len(a.Title) == 0 {
		a.Title = addr.String()
	}

	return nil
}

func articleContent(htmlSelection string, addr *url.URL, a *dto.Article) error {

	read, err := readability.FromReader(strings.NewReader(htmlSelection), addr)
	if err != nil {
		return nil
	}

	a.Markup = read.Content
	a.Text = lo.Ternary(a.Text == "", read.TextContent, a.Text)

	if read.Image != "" && !strings.Contains(a.Markup, read.Image) {
		// @note readability might remove the main image from the content
		// @see payload_test.articleMetadata()
		a.Markup = "<img src=\"" + read.Image + "\" />" + a.Markup
	}

	return nil
}
