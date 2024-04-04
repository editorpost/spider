package collect

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"log/slog"
	"strings"
)

// browse the URL this chromedp.Navigate, wait dom loaded and return the rendered HTML
func (crawler *Crawler) browse(reqURL string) (*goquery.Selection, error) {

	resp, err := crawler.chrome(reqURL)
	if err != nil {
		slog.Error("browser failed",
			slog.String("error", err.Error()),
			slog.String("url", reqURL),
		)
		return nil, err
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(resp))
	if err != nil {
		slog.Error("browser failed",
			slog.String("error", err.Error()),
			slog.String("url", reqURL),
		)
		return nil, err
	}

	return doc.Find(crawler.Query), nil
}

// browse the URL this chromedp.Navigate, wait dom loaded and return the rendered HTML
func (crawler *Crawler) chrome(reqURL string) (string, error) {

	// Initialize a new browser context
	ctx, cancel := chromedp.NewContext(context.Background())
	defer cancel()

	chromedp.UserAgent(crawler.UserAgent)

	// Navigate to the URL and fetch the rendered HTML
	var htmlContent string
	err := chromedp.Run(ctx,
		&emulation.SetUserAgentOverrideParams{
			UserAgent:      crawler.UserAgent,
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
