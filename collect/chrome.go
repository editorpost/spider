package collect

import (
	"context"
	"github.com/PuerkitoBio/goquery"
	"github.com/chromedp/cdproto/emulation"
	"github.com/chromedp/chromedp"
	"log/slog"
	"strings"
)

// Browse the Endpoint this chromedp.Navigate, wait dom loaded and return the rendered HTML
func (crawler *Crawler) Browse(reqURL string) (*goquery.Selection, error) {

	resp, err := crawler.browseChrome(reqURL)
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

	return doc.Find(crawler.args.ExtractSelector), nil
}

// Browse the Endpoint this chromedp.Navigate, wait dom loaded and return the rendered HTML
func (crawler *Crawler) browseChrome(reqURL string) (string, error) {

	// Initialize a new browser context
	ctx, cancel := chromedp.NewContext(crawler.chromeCtx)
	defer cancel()

	chromedp.UserAgent(crawler.args.UserAgent)

	// Navigate to the Endpoint and fetch the rendered HTML
	var htmlContent string
	err := chromedp.Run(ctx,
		&emulation.SetUserAgentOverrideParams{
			// todo: why UserAgent set twice?
			UserAgent:      crawler.args.UserAgent,
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

func (crawler *Crawler) setupChrome() context.CancelFunc {

	slog.Info("chrome collector")

	opts := []chromedp.ExecAllocatorOption{
		chromedp.NoFirstRun,
		chromedp.NoDefaultBrowserCheck,
		chromedp.Headless,
		chromedp.DisableGPU,
	}

	// create context
	ctx, cancelAlloc := chromedp.NewExecAllocator(context.Background(), opts...)
	crawler.chromeCtx = ctx

	// create context
	_, cancel := chromedp.NewContext(crawler.chromeCtx) // create new tab

	return func() {
		cancelAlloc()
		cancel()
	}
}
