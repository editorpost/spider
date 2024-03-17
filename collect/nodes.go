package collect

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"github.com/gocolly/colly/v2"
	"golang.org/x/net/html"
	"log/slog"
	"strings"
)

// nodes matching the query (with JS browse if Query is not found in GET response)
func (task Task) nodes(e *colly.HTMLElement) []*html.Node {

	entries := e.DOM.Find(task.Query)

	if entries.Length() == 0 {

		resp, err := task.browse(e.Request.URL.String())
		if err != nil {
			slog.Error("browser failed",
				slog.String("error", err.Error()),
				slog.String("url", e.Request.URL.String()),
			)
			return nil
		}

		doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
		if err != nil {
			slog.Error("browser failed",
				slog.String("error", err.Error()),
				slog.String("url", e.Request.URL.String()),
			)
			return nil
		}

		entries = doc.Find(task.Query)
	}

	var nodes []*html.Node
	entries.Each(func(i int, s *goquery.Selection) {
		nodes = append(nodes, s.Nodes...)
	})

	return nodes
}

// browse the URL this chromedp.Navigate, wait dom loaded and return the rendered HTML
func (task Task) browse(reqURL string) (string, error) {

	// Initialize a new browser context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	chromedp.UserAgent(task.UserAgent)

	// Navigate to the URL and fetch the rendered HTML
	var htmlContent string
	err := chromedp.Run(ctx,
		&emulation.SetUserAgentOverrideParams{
			UserAgent:      task.UserAgent,
			AcceptLanguage: "en-US,en;q=0.9",
		},
		chromedp.Navigate(reqURL),
		chromedp.WaitReady(`body`),
		chromedp.OuterHTML("html", &htmlContent),
	)
	if err != nil {
		return "", err
	}

	return htmlContent, nil
}
