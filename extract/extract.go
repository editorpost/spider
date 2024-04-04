package extract

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
	"log/slog"
	"net/url"
	"strings"
)

const (
	EntityArticle = "article"
	UrlField      = "url"
	HtmlField     = "html"
)

type (
	PipeFn    func(*Payload) error
	ExtractFn func(*goquery.Selection, *url.URL) error

	Payload struct {
		HTML      string
		Selection *goquery.Selection
		URL       *url.URL
		Data      map[string]any
	}
)

// Pipe is a function to process the html node and url.
// Order of the pipes is important.
func Pipe(pipes ...PipeFn) ExtractFn {

	return func(s *goquery.Selection, u *url.URL) error {

		if s == nil {
			return errors.New("document is nil")
		}

		str, err := s.Html()
		if err != nil {
			return err
		}

		payload := &Payload{
			HTML:      str,
			Selection: s,
			URL:       u,
			Data:      make(map[string]any),
		}

		for _, pipe := range pipes {
			err := pipe(payload)
			if err != nil {
				return err
			}
		}

		return nil
	}
}

func Article(p *Payload) error {

	article, err := readability.FromReader(strings.NewReader(p.HTML), p.URL)
	if err != nil {
		slog.Error("extract failed", err)
		return err
	}

	p.Data["entity__type"] = EntityArticle
	p.Data["entity__title"] = article.Title
	p.Data["entity__byline"] = article.Byline
	p.Data["entity__content"] = article.Content
	p.Data["entity__text"] = article.TextContent
	p.Data["entity__length"] = article.Length
	p.Data["entity__excerpt"] = article.Excerpt
	p.Data["entity__site"] = article.SiteName
	p.Data["entity__image"] = article.Image
	p.Data["entity__favicon"] = article.Favicon
	p.Data["entity__language"] = article.Language
	p.Data["entity__published"] = article.PublishedTime
	p.Data["entity__modified"] = article.ModifiedTime

	slog.Debug("extract success", slog.String("title", article.Title))

	return nil
}

func Crawler(p *Payload) (err error) {

	p.Data[UrlField] = p.URL.String()
	p.Data[HtmlField], err = p.Selection.Html()

	return err
}
