package setup

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/extract"
	"log/slog"
)

const (
	// DefaultMongoResource is the name of the mongo resource
	DefaultMongoResource = "f/spider/resource/mongodb"
)

// Config is the setup with primitive types
type Config struct {
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
	// ProxyEnabled is the flag to enable proxy or send requests directly
	ProxyEnabled bool `json:"ProxyEnabled"`
	// ProxySources is the list of proxy sources, expected to return list of proxies URLs.
	// If empty, the default proxy sources is used.
	ProxySources []string `json:"ProxySources"`
	// MongoDbResource is the name of the mongo resource, e.g. "u/spider/mongodb"
	MongoDbResource string `json:"MongoDbResource" validate:"trim,required"`
	// VictoriaMetricsUrl // todo move to resource
	VictoriaMetricsUrl string `json:"VictoriaMetricsUrl" validate:"trim,required"`
	// VictoriaLogsUrl // todo move to resource
	VictoriaLogsUrl string `json:"VictoriaLogsUrl" validate:"trim,required"`
}

// MustMongoConfig returns the mongo config or panic
func MustMongoConfig(resource string) *mongodb.Config {

	if len(resource) == 0 {
		resource = DefaultMongoResource
	}

	cfg, err := mongodb.GetResource(resource)
	if err != nil {
		slog.Warn("failed to get mongo resource", slog.String("error", err.Error()))
		return nil
	}

	return cfg
}
