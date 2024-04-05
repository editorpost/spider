package spider

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/donq/pkg/script"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
	"log/slog"
)

// Args is a minimal required input arguments for the spider
type Args struct {
	// Name is the name of the spider and mongo collection
	Name string `json:"Name" validate:"trim,required"`
	// StartURL is the URL to start crawling, e.g. http://example.com
	StartURL string `json:"StartURL" validate:"trim,required"`
	// AllowedURL is the regex to match the URLs, e.g. "https://example.com/articles/.+"
	AllowedURL string `json:"AllowedURL" validate:"trim,required"`
	// EntityURL is the URL to extract, e.g. "https://example.com/articles/science/.+"
	EntityURL string `json:"EntityURL" validate:"trim"`
	// UseBrowser is a flag to use browser for rendering the page
	UseBrowser bool `json:"UseBrowser"`
	// Depth is the number of levels to follow the links
	Depth int `json:"Depth"`
	// Selector CSS to match the entities to extract, e.g. ".article--ssr"
	Selector string `json:"Selector" validate:"trim,required"`
}

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
	crawler := &collect.Crawler{
		StartURL:       args.StartURL,
		AllowedURL:     args.AllowedURL,
		EntityURL:      args.EntityURL,
		UseBrowser:     args.UseBrowser,
		Depth:          args.Depth,
		EntitySelector: args.Selector,
		Extractor:      Extract(args.Name, extract.Article),
		Collector:      nil, // use colly default in-memory storage
	}

	return crawler.Start()
}

// Extract creates Pipe with given extractor called before Save
func Extract(dbName string, extractor extract.PipeFn) extract.ExtractFn {

	cfg, err := mongodb.GetResource("u/spider/mongo")
	if err != nil {
		slog.Error("failed to get mongo resource", slog.String("error", err.Error()))
		panic(err)
	}

	storage, err := store.NewExtractStore(dbName, cfg)
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
