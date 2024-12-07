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

	// remove h1 from the selection
	sel.Find("h1").Remove()

	// remove url links from the selection
	// replace text links with text content
	replaceLinks(sel)

	// get selection html content
	content, err := sel.Html()
	if err != nil {
		return nil, fmt.Errorf("failed to get HTML from selection: %w", err)
	}

	// fill data from readability
	// note article.Markup is still HTML here, not markdown
	if err = readabilityArticle(payload, content, a); err != nil {
		return nil, fmt.Errorf("failed to get readability article: %w", err)
	}

	// fallback: published
	a.Published = lo.Ternary(a.Published.IsZero(), legacyPublished(content), a.Published)
	a.Author = lo.Ternary(a.Author == "", legacyAuthor(content), a.Author)

	// html to markdown
	// article.Markup is now converted to markdown
	if a.Markup, err = HTMLToMarkdown(a.Markup, HostUrl(payload.URL)); err != nil {
		return nil, err
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

func readabilityArticle(payload *pipe.Payload, content string, a *dto.Article) error {

	// get head tags as a string
	head, err := payload.Doc.DOM.Find("head").Html()
	if err != nil {
		return fmt.Errorf("failed to get head tags: %w", err)
	}

	// and attach to the html content
	html := strings.Join([]string{head, content}, "")

	// readability: title, summary, text, markup, html, language, summary
	read, err := readability.FromReader(strings.NewReader(html), payload.URL)
	if err != nil {
		return nil
	}

	a.Title = read.Title
	a.Text = read.TextContent
	a.Markup = read.Content
	a.Language = read.Language
	a.Summary = lo.Ternary(len(read.Byline) > 0, read.Byline, read.Excerpt)

	// fallback fields applied only if the fields are empty
	a.Title = lo.Ternary(a.Title == "", read.Title, a.Title)
	a.Summary = lo.Ternary(a.Summary == "", read.Excerpt, a.Summary)
	a.Text = lo.Ternary(a.Text == "", read.TextContent, a.Text)
	a.Language = lo.Ternary(a.Language == "", read.Language, a.Language)

	a.Markup = read.Content

	if read.PublishedTime != nil {
		a.Published = *read.PublishedTime
	}

	if len(a.Title) == 0 {
		a.Title = payload.URL.String()
	}

	return nil
}
