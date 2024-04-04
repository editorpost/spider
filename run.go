package spider

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
	"log/slog"
)

type WindmillArgs struct {
	StartURL string
	MatchURL string
	Depth    int
	Query    string
}

// WindmillArticleExample is an example code for running spider
// as Windmill Script with extract.Article
func WindmillArticleExample() error {

	crawler := &collect.Crawler{
		StartURL:  "http://example.com",
		MatchURL:  ".*",
		Depth:     1,
		Query:     ".article--ssr",
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
