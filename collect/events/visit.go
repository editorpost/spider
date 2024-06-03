package events

import (
	"github.com/editorpost/spider/collect/config"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"strings"
)

// visit links found in the DOM
func (crawler *Dispatch) visit() func(e *colly.HTMLElement) {

	return func(e *colly.HTMLElement) {

		// absolute url
		link := e.Request.AbsoluteURL(e.Attr("href"))

		// skip empty and anchor links
		if link == "" || strings.HasPrefix(link, "#") {
			return
		}

		// skip images, scripts, etc.
		if !config.ContentLikeURL(link) {
			return
		}

		// visit the link
		if err := crawler.queue.AddURL(link); err != nil {
			slog.Warn("crawler queue", slog.String("error", err.Error()))
		}
	}
}
