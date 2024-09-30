package article

import (
	md "github.com/JohannesKaufmann/html-to-markdown"
	"github.com/PuerkitoBio/goquery"
	dto "github.com/editorpost/article"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/extract/pipe"
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

// Article extracts the dto from the HTML
// and sets the dto fields to the payload
func Article(payload *pipe.Payload) error {

	htmlStr, err := ArticleSelection(payload)
	if err != nil {
		slog.Warn("failed to get HTML from selection", slog.String("err", err.Error()))
		return err
	}

	art, err := ArticleFromHTML(htmlStr, payload.URL)
	if err != nil {
		slog.Warn("failed to extract dto", slog.String("err", err.Error()))
		return err
	}

	// rewrite the markup from readability and distiller due to the different behavior
	// use provided Config.ExtractSelector markup as a default
	// todo: if extract selector is empty, then use the readability markup
	if art.Markup, err = ArticleSelectionToMarkup(payload); err != nil {
		slog.Warn("failed to get HTML from selection", slog.String("err", err.Error()))
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

// ArticleSelection combine head tags and article html code selection
// Head tags are required for readability to work properly
func ArticleSelection(payload *pipe.Payload) (dom string, err error) {

	if dom, err = payload.Selection.Html(); err != nil {
		slog.Warn("failed to get HTML from selection", slog.String("err", err.Error()))
		return "", err
	}

	// get head tags as a string and attach to rest of the html
	head, _ := payload.Doc.DOM.Find("head").Html()

	return strings.Join([]string{head, dom}, ""), nil
}

// ArticleFromHTML extracts Article
func ArticleFromHTML(html string, resource *url.URL) (*dto.Article, error) {

	a := dto.NewArticle()
	a.SourceURL = resource.String()

	// readability: title, summary, text, html, language
	readabilityArticle(html, resource, a)

	// fallback: published
	a.Published = lo.Ternary(a.Published.IsZero(), legacyPublished(html), a.Published)
	a.Author = lo.Ternary(a.Author == "", legacyAuthor(html), a.Author)

	// html to markdown
	a.Markup = HTMLToMarkdown(a.Markup)

	// todo: strip non-image links
	// a.Markup = StripMarkdown(a.Markup)

	// nil dto if it's invalid
	if err := a.Normalize(); err != nil {
		return nil, err
	}

	return a, nil
}

func ArticleSelectionToMarkup(payload *pipe.Payload) (markup string, err error) {

	if markup, err = payload.Selection.Html(); err != nil {
		return
	}

	return HTMLToMarkdown(markup), nil
}

func readabilityArticle(html string, resource *url.URL, a *dto.Article) {

	read, err := readability.FromReader(strings.NewReader(html), resource)
	if err != nil {
		return
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
}

func distillArticle(html string, resource *url.URL, a *dto.Article) {

	distill, err := distiller.ApplyForReader(strings.NewReader(html), &distiller.Options{
		OriginalURL: resource,
	})
	if err != nil {
		return
	}

	info := distill.MarkupInfo

	// set the dto fields
	a.SourceName = info.Publisher
	a.Author = info.Author

	// fallback fields applied only if the fields are empty
	a.Title = lo.Ternary(a.Title == "", distill.Title, a.Title)
	a.Summary = lo.Ternary(a.Summary == "", info.Description, a.Summary)
	a.Text = lo.Ternary(a.Text == "", distill.Text, a.Text)

	// a.Markup = lo.Ternary(a.Markup == "", dom.OuterHTML(distill.Node), a.Markup)
	a.Markup = dom.OuterHTML(distill.Node)

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

func distillImages(distill *distiller.Result, resource *url.URL) *dto.Images {

	images := dto.NewImages()

	for _, src := range distill.ContentImages {
		image := dto.NewImage(media.AbsoluteUrl(resource, src))
		images.Add(image)
	}

	return images
}

// HTMLToMarkdown converts HTML to markdown
func HTMLToMarkdown(html string) string {
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
