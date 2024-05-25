package manage

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/donq/pkg/script"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
	"log/slog"
)

// StartWith is an example code for running spider
// as Windmill Script with extract.Article
//
//goland:noinspection GoUnusedExportedFunction
func StartWith(input any) error {

	args := &Args{}
	if err := script.ParseArgs(input, args); err != nil {
		return err
	}

	return Start(args)
}

// Start is an example code for running spider
// as Windmill Script with extract.Article
func Start(args *Args) error {

	// fallback to user extractor
	extractor := args.EntityExtract
	if extractor == nil {
		extractor = func(*extract.Payload) error {
			return nil
		}
	}

	proxies := MustProxyPool(args)
	mongoCfg := MustMongoConfig(args.MongoDbResource)

	// create the crawler
	crawler := &collect.Crawler{
		StartURL:       args.StartURL,
		AllowedURL:     args.AllowedURL,
		EntityURL:      args.EntityURL,
		UseBrowser:     args.UseBrowser,
		Depth:          args.Depth,
		EntitySelector: args.EntitySelector,
		UserAgent:      "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/123.0.0.0 Safari/537.36",
		Extractor:      MustExtractor(args.Name, mongoCfg, extractor),
		Collector:      MustCollector(args.Name, mongoCfg),
		RoundTripper:   proxies.Transport(),
	}

	return crawler.Start()
}

// MustExtractor creates Pipe with given extractor called before Save
func MustExtractor(dbName string, cfg *mongodb.Config, extractor extract.PipeFn) extract.ExtractFn {

	// note: the `Save` must provide `created` and `updated` fields behavior
	storage, err := store.NewExtractStore(dbName, cfg)
	if err != nil {
		slog.Error("failed to create extract store", slog.String("error", err.Error()))
		panic(err)
	}

	return extract.Pipe(WindmillMeta, extract.Html, extractor, storage.Save)
}

// MustCollector creates a new collector store
func MustCollector(dbName string, cfg *mongodb.Config) *store.CollectStore {

	s, err := store.NewCollectStore(dbName, cfg)
	if err != nil {
		slog.Error("failed to create collect store", slog.String("error", err.Error()))
		panic(err)
	}

	return s
}

func MustProxyPool(args *Args) *proxy.Pool {

	// start the proxy pool
	pool := proxy.NewPool(args.StartURL)

	// provide user defined proxy sources
	// or used default public sources
	if len(args.ProxySources) > 0 {
		pool.Loader = func() ([]string, error) {
			return proxy.LoadStringLists(args.ProxySources)
		}
	}

	if err := pool.Start(); err != nil {
		panic(err)
	}

	return pool
}
