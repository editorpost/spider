package setup

import (
	"fmt"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
)

// Config is the configuration for the spider
// JSON:
//
//		{
//		  "MongoDSN": "mongodb://localhost:27017",
//		  "VictoriaMetricsUrl": "http://localhost:8428",
//		  "VictoriaLogsUrl": "http://localhost:8429"
//		}
//	 db.updateUser("spider", {roles: [{ role: "readAnyDatabase", db: "admin"}
type Config struct {
	// MongoDSN is connection string to MongoDB
	MongoDSN string `json:"MongoDSN" validate:"trim"`
	// VictoriaMetricsUrl // todo move to resource
	VictoriaMetricsUrl string `json:"VictoriaMetricsUrl" validate:"trim"`
	// VictoriaLogsUrl // todo move to resource
	VictoriaLogsUrl string `json:"VictoriaLogsUrl" validate:"trim"`
}

func Deps(args *config.Args, deploy *Config, extractors ...extract.Extractor) (*config.Deps, error) {

	if err := args.Normalize(); err != nil {
		return nil, err
	}

	if deploy.VictoriaLogsUrl != "" {
		VictoriaLogs(deploy.VictoriaLogsUrl, "info", args.ID)
	}

	// prepend windmill
	extractors = append([]extract.Extractor{extract.WindmillMeta}, extractors...)

	deps := &config.Deps{}

	// database
	if deploy.MongoDSN != "" {

		collectStore, err := store.NewCollectStore(args.ID, deploy.MongoDSN)
		if err != nil {
			return nil, err
		}
		deps.Storage = collectStore

		// connect storage
		extractStore, err := store.NewExtractStore(args.ID, deploy.MongoDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to create extract store: %v", err)
		}

		// provide save extractor func
		extractors = append(extractors, extractStore.Save)
	}

	deps.Extractor = extract.Pipe(extractors...)

	// metrics
	if deploy.VictoriaMetricsUrl != "" {
		metrics, err := NewMetrics(vars.FromEnv().JobID, args.ID, deploy.VictoriaMetricsUrl)
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
