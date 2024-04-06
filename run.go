package spider

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/donq/pkg/script"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
	"log/slog"
)

const (
	// DefaultMongoResource is the name of the mongo resource
	DefaultMongoResource = "f/spider/resource/mongodb"
)

// Args is a minimal required input arguments for the spider
type Args struct {
	// Name is the name of the spider and mongo collection
	Name string `json:"Name" validate:"trim,required"`
	// StartURL is the URL to start crawling, e.g. http://example.com
	StartURL string `json:"StartURL" validate:"trim,required"`
	// AllowedURL is the regex to match the URLs, e.g. "https://example.com?.+"
	AllowedURL string `json:"AllowedURL" validate:"trim,required"`
	// EntityURL is the URL to extract, e.g. "https://example.com/articles/((?:[^/]+/)*[^/]+)/.+"
	EntityURL string `json:"EntityURL" validate:"trim"`
	// Extractor is the function to process matched the data
	EntityExtract func(*extract.Payload) error
	// EntitySelector CSS to match the entities to extract, e.g. ".article--ssr"
	EntitySelector string `json:"EntitySelector" validate:"trim,required"`
	// UseBrowser is a flag to use browser for rendering the page
	UseBrowser bool `json:"UseBrowser"`
	// Depth is the number of levels to follow the links
	Depth int `json:"Depth"`
	// MongoDbResource is the name of the mongo resource, e.g. "u/spider/mongodb"
	MongoDbResource string `json:"MongoDbResource" validate:"trim,required"`
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

	extractor := args.EntityExtract

	if extractor == nil {
		extractor = func(*extract.Payload) error {
			return nil
		}
	}

	mongoCfg := MustMongoConfig(args.MongoDbResource)

	crawler := &collect.Crawler{
		StartURL:       args.StartURL,
		AllowedURL:     args.AllowedURL,
		EntityURL:      args.EntityURL,
		UseBrowser:     args.UseBrowser,
		Depth:          args.Depth,
		EntitySelector: args.EntitySelector,
		Extractor:      MustExtractor(mongoCfg, args.Name, extractor),
		Collector:      nil, // use colly default in-memory storage
	}

	return crawler.Start()
}

// MustExtractor creates Pipe with given extractor called before Save
func MustExtractor(cfg *mongodb.Config, dbName string, extractor extract.PipeFn) extract.ExtractFn {

	// note: the `Save` must provide `created` and `updated` fields behavior
	storage, err := store.NewExtractStore(dbName, cfg)
	if err != nil {
		slog.Error("failed to create extract store", slog.String("error", err.Error()))
		panic(err)
	}

	return extract.Pipe(WindmillMeta, extract.Html, extractor, storage.Save)
}

// MetricStore creates a new metric store
func MustMetricStore(cfg *mongodb.Config) *store.MetricStore {

	s, err := store.NewMetricStore(cfg)
	if err != nil {
		slog.Error("failed to create metric store", slog.String("error", err.Error()))
		panic(err)
	}

	return s
}

// MustMongoConfig returns the mongo config or panic
func MustMongoConfig(resource string) *mongodb.Config {

	if len(resource) == 0 {
		resource = DefaultMongoResource
	}

	cfg, err := mongodb.GetResource(resource)
	if err != nil {
		slog.Error("failed to get mongo resource", slog.String("error", err.Error()))
		panic(err)
	}

	return cfg
}

// WindmillMeta is a meta data extractor
func WindmillMeta(p *extract.Payload) error {

	// windmill flow
	p.Data["job_id"] = Env().GetRootFlowJobID()
	p.Data["flow_path"] = Env().GetFlowPath()
	p.Data["flow_job_id"] = Env().GetFlowJobID()

	return nil
}
