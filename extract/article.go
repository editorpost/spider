package extract

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"github.com/go-shiori/go-readability"
	"log/slog"
	"strings"
)

func Article(p *Payload) error {
	if p.Selection == nil {
		err := errors.New("selection is nil")
		slog.Warn("selection is nil", err)
		return err
	}

	htmlStr, err := p.Selection.Html()
	if err != nil {
		slog.Warn("failed to get HTML from selection", err)
		return err
	}

	// Удаление рекламных баннеров и скриптов
	doc, err := goquery.NewDocumentFromReader(strings.NewReader(htmlStr))
	if err != nil {
		slog.Error("failed to parse HTML", err)
		return err
	}

	doc.Find("script, .advertisement, .ad, .ads, iframe").Each(func(i int, s *goquery.Selection) {
		s.Remove()
	})

	cleanedHTML, err := doc.Html()
	if err != nil {
		slog.Error("failed to get cleaned HTML", err)
		return err
	}

	article, err := readability.FromReader(strings.NewReader(cleanedHTML), p.URL)
	if err != nil {
		slog.Error("failed to extract article", err)
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
