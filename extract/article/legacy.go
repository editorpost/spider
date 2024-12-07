package article

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/goodsign/monday"
	"strings"
	"time"
)

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
