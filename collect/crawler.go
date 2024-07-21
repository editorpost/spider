package collect

import (
	"context"
	"github.com/editorpost/spider/collect/config"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"log/slog"
)

// Crawler for scraping a website
type Crawler struct {
	args      *config.Args
	deps      *config.Deps
	queue     *queue.Queue
	collect   *colly.Collector
	chromeCtx context.Context
}

func NewCrawler(args *config.Args, deps *config.Deps) (*Crawler, error) {

	if err := args.Normalize(); err != nil {
		return nil, err
	}

	crawler := &Crawler{
		args: args,
		deps: deps.Normalize(),
	}

	if _, err := crawler.collector(); err != nil {
		return nil, err
	}

	return crawler, nil
}

// Run the scraping Crawler.
func (crawler *Crawler) Run() error {

	if crawler.args.UseBrowser {
		// create chrome allocator context
		cancel := crawler.setupChrome()
		// disable async in browser mode
		crawler.collect.Async = false
		defer cancel()
	}

	slog.Info("collector starting", crawler.args.Log())

	if err := crawler.queue.AddURL(crawler.args.StartURL); err != nil {
		return err
	}

	if err := crawler.queue.Run(crawler.collect); err != nil {
		return err
	}

	crawler.collect.Wait()

	return nil
}

// Stop the scraping Crawler (takes a while to finish).
func (crawler *Crawler) Stop() {
	crawler.queue.Stop()
}
