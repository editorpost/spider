package article

import (
	dto "github.com/editorpost/article"
	"github.com/editorpost/spider/extract/media"
	"github.com/go-shiori/dom"
	distiller "github.com/markusmobius/go-domdistiller"
	"github.com/samber/lo"
	"net/url"
	"strings"
	"time"
)

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
