package manage

import (
	"fmt"
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/donq/pkg/script"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/store"
	"github.com/gocolly/colly/v2/storage"
)

// StartWith is an example code for running spider
// as Windmill Script with extract.Article
//
//goland:noinspection GoUnusedExportedFunction
func StartWith(input any) error {

	args := &setup.Config{}
	if err := script.ParseArgs(input, args); err != nil {
		return err
	}

	return Start(args)
}

// Start is a code for running spider
// as Windmill Script with extract.Article
func Start(args *setup.Config) error {

	// defaults
	var collectStore storage.Storage = &storage.InMemoryStorage{}

	// prepend windmill
	extractors := append([]extract.PipeFn{extract.WindmillMeta}, args.EntityExtract)

	// database
	if args.MongoDbResource != "" {

		// load mongo windmill resource
		mongo, err := mongodb.GetResource(args.MongoDbResource)
		if err != nil {
			return err
		}

		collectStore, err = store.NewCollectStore(args.Name, mongo)
		if err != nil {
			return err
		}

		// connect storage
		extractStore, err := store.NewExtractStore(args.Name, mongo)
		if err != nil {
			return fmt.Errorf("failed to create extract store: %v", err)
		}

		// provide save extractor func
		extractors = append(extractors, extractStore.Save)
	}

	// crawler setup
	crawler := &collect.Crawler{
		StartURL:       args.StartURL,
		AllowedURL:     args.AllowedURL,
		EntityURL:      args.EntityURL,
		UseBrowser:     args.UseBrowser,
		Depth:          args.Depth,
		EntitySelector: args.EntitySelector,
		UserAgent:      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		Extractor:      extract.Pipe(extractors...),
		Storage:        collectStore,
		JobID:          vars.FromEnv().JobID,
		SpiderID:       args.Name,
		Monitor:        &collect.MetricsFallback{},
	}

	// metrics
	if args.VictoriaMetricsUrl != "" {
		metrics, err := setup.NewMetrics(vars.FromEnv().JobID, args.Name, args.VictoriaMetricsUrl)
		if err != nil {
			return err
		}
		crawler.Monitor = metrics
	}

	// proxy
	if args.ProxyEnabled {
		proxies, err := proxy.StartPool(args.StartURL, args.ProxySources...)
		if err != nil {
			return err
		}
		crawler.RoundTripper = proxies.Transport()
	}

	return crawler.Start()
}
