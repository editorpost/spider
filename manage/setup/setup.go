package setup

import (
	"fmt"
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
)

const (
	// DefaultMongoResource is the name of the mongo resource
	DefaultMongoResource = "f/spider/resource/mongodb"
)

type Deploy struct {
	// JobID is the unique identifier for the job
	JobID string `json:"JobID"`
	// SpiderID is the unique identifier for the spider
	SpiderID string `json:"SpiderID"`
	// MongoDbResource is the name of the mongo resource, e.g. "u/spider/mongodb"
	MongoDbResource string `json:"MongoDbResource" validate:"trim,required"`
	// VictoriaMetricsUrl // todo move to resource
	VictoriaMetricsUrl string `json:"VictoriaMetricsUrl" validate:"trim,required"`
	// VictoriaLogsUrl // todo move to resource
	VictoriaLogsUrl string `json:"VictoriaLogsUrl" validate:"trim,required"`
}

// Crawler setup
func Deps(args *config.Args, deploy *Deploy, extractors ...extract.PipeFn) (*config.Deps, error) {

	if err := args.Normalize(); err != nil {
		return nil, err
	}

	if deploy.VictoriaLogsUrl != "" {
		VictoriaLogs(deploy.VictoriaLogsUrl, "info", deploy.SpiderID)
	}

	// prepend windmill
	extractors = append([]extract.PipeFn{extract.WindmillMeta}, extractors...)

	deps := &config.Deps{}

	// database
	if deploy.MongoDbResource != "" {

		// load mongo windmill resource
		mongo, err := mongodb.GetResource(deploy.MongoDbResource)
		if err != nil {
			return nil, err
		}

		collectStore, err := store.NewCollectStore(deploy.SpiderID, mongo)
		if err != nil {
			return nil, err
		}
		deps.Storage = collectStore

		// connect storage
		extractStore, err := store.NewExtractStore(deploy.SpiderID, mongo)
		if err != nil {
			return nil, fmt.Errorf("failed to create extract store: %v", err)
		}

		// provide save extractor func
		extractors = append(extractors, extractStore.Save)
	}

	deps.Extractor = extract.Pipe(extractors...)

	// metrics
	if deploy.VictoriaMetricsUrl != "" {
		metrics, err := NewMetrics(vars.FromEnv().JobID, deploy.SpiderID, deploy.VictoriaMetricsUrl)
		if err != nil {
			return nil, err
		}
		deps.Monitor = metrics
	}

	// proxy
	if args.ProxyEnabled {
		proxies, err := proxy.StartPool(args.StartURL, args.ProxySources...)
		if err != nil {
			return nil, err
		}
		deps.RoundTripper = proxies.Transport()
	}

	return deps, nil
}
