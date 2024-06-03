package collect

import (
	"context"
	"github.com/editorpost/spider/collect/config"
	"github.com/gocolly/colly/v2"
	"github.com/gocolly/colly/v2/queue"
	"log/slog"
	"time"
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

	slog.Info("collector running")

	crawler.collect.Wait()

	slog.Info("collector finishing")
	time.Sleep(500 * time.Millisecond) // todo: remove after ensuring graceful shutdown works as expected

	return nil
}
