package extract

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	"github.com/editorpost/article"
	"github.com/go-shiori/dom"
	"github.com/go-shiori/go-readability"
	"github.com/goodsign/monday"
	distiller "github.com/markusmobius/go-domdistiller"
	"github.com/samber/lo"
	"log/slog"
	"net/url"
	"strings"
	"time"
)

// Article extracts the article from the HTML
// and sets the article fields to the payload
func Article(p *Payload) error {

	if p.Selection == nil {
		p.Selection = p.Doc.DOM
	}

	htmlStr, err := p.Selection.Html()
	if err != nil {
		slog.Warn("failed to get HTML from selection", slog.String("err", err.Error()))
		return err
	}

	art, err := ArticleFromHTML(htmlStr, p.URL)
	if err != nil {
		slog.Warn("failed to extract article", slog.String("err", err.Error()))
		return err
	}

	// set the article to the payload
	for k, v := range art.Map() {
		p.Data[k] = v
	}

	slog.Debug("extract success", slog.String("title", art.Title))

	return nil
}

// ArticleFromHTML extracts Article
func ArticleFromHTML(html string, resource *url.URL) (*article.Article, error) {

	a := article.NewArticle()
	a.SourceURL = resource.String()

	// readability: title, summary, text, html, language
	readabilityArticle(html, resource, a)

	// distiller: category, images, source name, author
	distillArticle(html, resource, a)

	// fallback: published
	a.Published = lo.Ternary(a.Published.IsZero(), legacyPublished(html), a.Published)
	a.Author = lo.Ternary(a.Author == "", legacyAuthor(html), a.Author)

	// html to markdown
	a.Markup = markupText(a.Markup)

	// nil article if it's invalid
	if err := a.Normalize(); err != nil {
		return nil, err
	}

	return a, nil
}

func readabilityArticle(html string, resource *url.URL, a *article.Article) {

	read, err := readability.FromReader(strings.NewReader(html), resource)
	if err != nil {
		return
	}

	// set the article fields
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
}

func distillArticle(html string, resource *url.URL, a *article.Article) {

	distill, err := distiller.ApplyForReader(strings.NewReader(html), &distiller.Options{
		OriginalURL: resource,
	})
	if err != nil {
		return
	}

	info := distill.MarkupInfo

	// set the article fields
	a.Category = info.Article.Section
	a.SourceName = info.Publisher
	a.Images = distillImages(distill, resource)
	a.Author = info.Author

	// fallback fields applied only if the fields are empty
	a.Title = lo.Ternary(a.Title == "", distill.Title, a.Title)
	a.Summary = lo.Ternary(a.Summary == "", info.Description, a.Summary)
	a.Text = lo.Ternary(a.Text == "", distill.Text, a.Text)
	a.Markup = lo.Ternary(a.Markup == "", dom.OuterHTML(distill.Node), a.Markup)
	a.Published = lo.Ternary(a.Published.IsZero(), distillPublished(distill), a.Published)
}

func distillPublished(distill *distiller.Result) time.Time {

	publishedStr := distill.MarkupInfo.Article.PublishedTime
	published, timeErr := time.Parse(time.RFC3339, publishedStr)
	if timeErr == nil {
		return time.Now()
	}
	return published
}

func distillImages(distill *distiller.Result, resource *url.URL) *article.Images {

	images := article.NewImages()

	for _, src := range distill.MarkupInfo.Images {
		image := article.NewImage(AbsoluteUrl(resource, src.URL))
		image.Width = src.Width
		image.Height = src.Height
		image.Title = src.Caption
		images.Add(image)
	}

	return images
}

// markdown converts HTML to markdown
func markupText(html string) string {
	converter := md.NewConverter("", true, nil)
	markdown, err := converter.ConvertString(html)
	return lo.Ternary(err == nil, markdown, "")
}

func legacyPublished(html string) time.Time {

	fallback := time.Now()

	q, readerErr := goquery.NewDocumentFromReader(strings.NewReader(html))
	if readerErr != nil {
		return fallback
	}

	// .field--name-created
	if el := q.Find(".field--name-created").Text(); len(el) > 0 {

		// Monday,2 January 2006 format
		published, err := monday.Parse("Monday, 2 January 2006", el, monday.LocaleRuRU)
		if err == nil {
			return published
		}
	}

	return fallback
}

func legacyAuthor(html string) (name string) {

	q, readerErr := goquery.NewDocumentFromReader(strings.NewReader(html))
	if readerErr != nil {
		return
	}

	// look at publisher info
	for _, node := range q.Find(".node-article__date").Nodes {
		if node.FirstChild != nil {
			name = strings.TrimSpace(node.FirstChild.Data)
			return
		}
	}

	return
}

func AbsoluteUrl(base *url.URL, href string) string {

	// parse the href
	rel, err := url.Parse(href)
	if err != nil {
		return ""
	}

	// resolve the base with the relative href
	abs := base.ResolveReference(rel)

	return abs.String()
}
