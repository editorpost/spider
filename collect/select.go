package collect

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
)

// selections matching the query (with JS browse if EntitySelector is not found in GET response)
func (crawler *Crawler) selections(e *colly.HTMLElement) []*goquery.Selection {

	selections := e.DOM.Find(crawler.EntitySelector)

	// if the Selector is not found in the GET response,
	// but in the fallback js browser call
	if crawler.UseBrowser {
		var err error
		selections, err = crawler.browse(e.Request.URL.String())
		if err != nil {
			return nil
		}
	}

	var nodes []*goquery.Selection
	selections.Each(func(i int, s *goquery.Selection) {
		nodes = append(nodes, s)
	})

	return nodes
}
