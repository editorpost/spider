package spider

import (
	"github.com/editorpost/donq/mongodb"
	"log/slog"
	"spider/collect"
	"spider/extract"
	"spider/store"
)

// WindmillExample is an example code for running spider as Windmill Script
func WindmillExample() error {

	task := &collect.Task{
		StartURL: "http://example.com",
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article--ssr",
		Extract:  Extract(extract.Article),
		Storage:  nil, // use colly default in-memory storage
	}

	return task.Start()
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

	return extract.Pipe(WindmillMeta, extract.Crawler, extractor, storage.Save)
}

// WindmillMeta is a meta data extractor
func WindmillMeta(p *extract.Payload) error {
	p.Data["job_id"] = Env().GetRootFlowJobID()
	return nil
}
