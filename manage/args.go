package manage

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
	"log/slog"
)

const (
	// DefaultMongoResource is the name of the mongo resource
	DefaultMongoResource = "f/spider/resource/mongodb"
	// WindmillJobID is the key for the job ID
	WindmillJobID = "windmill__job_id"
	// WindmillFlowPath is the key for the flow path
	WindmillFlowPath = "windmill__flow_path"
	// WindmillFlowJobID is the key for the flow job ID
	WindmillFlowJobID = "windmill__flow_job_id"
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
	// ProxySourceURLs is the list of proxy sources, expected to return list of proxies URLs
	// by default used public proxy sources
	ProxySources []string `json:"ProxySources"`
	// MongoDbResource is the name of the mongo resource, e.g. "u/spider/mongodb"
	MongoDbResource string `json:"MongoDbResource" validate:"trim,required"`
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
