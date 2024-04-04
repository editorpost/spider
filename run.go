package spider

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
	"log/slog"
)

type Args struct {
	// StartURL is the URL to start crawling, e.g. http://example.com
	StartURL string
	// MatchURL is the regex to match the URLs, e.g. ".*"
	MatchURL string
	// Depth is the number of levels to follow the links
	Depth int
	// Query is the CSS selector to match the elements, e.g. ".article--ssr"
	Query string
}

// Start is an example code for running spider
// as Windmill Script with extract.Article
func Start(args *Args) error {

	crawler := &collect.Crawler{
		StartURL:  args.StartURL,
		MatchURL:  args.MatchURL,
		Depth:     args.Depth,
		Query:     args.Query,
		Extractor: Extract(extract.Article),
		Collector: nil, // use colly default in-memory storage
	}

	return crawler.Start()
}

// Extract creates Pipe with given extractor called before Save
func Extract(extractor extract.PipeFn) extract.ExtractFn {

	cfg, err := mongodb.GetResource("u/spider/mongo")
	if err != nil {
		slog.Error("failed to get mongo resource", slog.String("error", err.Error()))
		panic(err)
	}

	jobID := Env().GetRootFlowJobID()
	storage, err := store.NewExtractStore(jobID, cfg)
	if err != nil {
		slog.Error("failed to create extract store", slog.String("error", err.Error()))
		panic(err)
	}

	return extract.Pipe(WindmillMeta, extract.Html, extractor, storage.Save)
}

// WindmillMeta is a meta data extractor
func WindmillMeta(p *extract.Payload) error {
	p.Data["job_id"] = Env().GetRootFlowJobID()
	return nil
}
