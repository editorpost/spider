package setup

import (
	"fmt"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/store"
)

type Deploy struct {
	// MongoDSN is the name of the mongo resource, e.g. "u/spider/mongodb"
	MongoDSN string `json:"MongoDSN" validate:"trim"`
	// VictoriaMetricsUrl // todo move to resource
	VictoriaMetricsUrl string `json:"VictoriaMetricsUrl" validate:"trim"`
	// VictoriaLogsUrl // todo move to resource
	VictoriaLogsUrl string `json:"VictoriaLogsUrl" validate:"trim"`
}

func Deps(args *config.Args, deploy *Deploy, extractors ...extract.PipeFn) (*config.Deps, error) {

	if err := args.Normalize(); err != nil {
		return nil, err
	}

	if deploy.VictoriaLogsUrl != "" {
		VictoriaLogs(deploy.VictoriaLogsUrl, "info", args.SpiderID)
	}

	// prepend windmill
	extractors = append([]extract.PipeFn{extract.WindmillMeta}, extractors...)

	deps := &config.Deps{}

	// database
	if deploy.MongoDSN != "" {

		collectStore, err := store.NewCollectStore(args.SpiderID, deploy.MongoDSN)
		if err != nil {
			return nil, err
		}
		deps.Storage = collectStore

		// connect storage
		extractStore, err := store.NewExtractStore(args.SpiderID, deploy.MongoDSN)
		if err != nil {
			return nil, fmt.Errorf("failed to create extract store: %v", err)
		}

		// provide save extractor func
		extractors = append(extractors, extractStore.Save)
	}

	deps.Extractor = extract.Pipe(extractors...)

	// metrics
	if deploy.VictoriaMetricsUrl != "" {
		metrics, err := NewMetrics(vars.FromEnv().JobID, args.SpiderID, deploy.VictoriaMetricsUrl)
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
