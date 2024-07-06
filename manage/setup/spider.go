package setup

import (
	"encoding/json"
	"fmt"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/store"
)

const (
	SpiderMediaPath = "spiders/%s/media/"
)

// Deploy is the configuration for the spider
// JSON:
//
//		{
//	     "MongoDSN": "mongodb://spider:pass@mongo-server-rs.spider.svc/spider?ssl=false",
//	     "VictoriaLogsUrl": "http://spider-victoria-logs-single-server.spider.svc:9428",
//	     "VictoriaMetricsUrl": "http://vmsingle-spider.spider.svc:8429/api/v1/import/prometheus",
//			"Bucket": {
//				"Name": "ep-spider",
//				"Endpoint": "https://s3.ap-southeast-1.wasabisys.com",
//				"Region": "ap-southeast-1",
//				"PublicURL": "http://localhost:9000",
//				"Access": "",
//				"Secret": "",
//			}
//		}
//
//		db.updateUser("spider", {roles: [{ role: "readAnyDatabase", db: "admin"}
type Deploy struct {
	// Bucket is the name of the bucket to store data
	Bucket store.Bucket `json:"Bucket"`
	// MongoDSN is connection string to MongoDB
	MongoDSN string `json:"MongoDSN" validate:"trim"`
	// VictoriaMetricsUrl
	VictoriaMetricsUrl string `json:"VictoriaMetricsUrl" validate:"trim"`
	// VictoriaLogsUrl
	VictoriaLogsUrl string `json:"VictoriaLogsUrl" validate:"trim"`
}

type Spider struct {
	*config.Args
	*extract.Config
	pipe *payload.Pipeline
}

// UnmarshalJSON is the custom unmarshalling for Spider
func (s *Spider) UnmarshalJSON(data []byte) error {
	type Alias Spider
	aux := &struct {
		*Alias
	}{
		Alias: (*Alias)(s),
	}
	if err := json.Unmarshal(data, &aux); err != nil {
		return err
	}
	return nil
}

func NewSpiderFromJSON(data []byte) (*Spider, error) {

	s := &Spider{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	return NewSpider(s.Args, s.Config)
}

func NewSpider(args *config.Args, cfg *extract.Config) (*Spider, error) {

	s := &Spider{
		Args:   args,
		Config: cfg,
	}

	if err := s.Args.Normalize(); err != nil {
		return nil, err
	}

	if err := s.withPipeline(); err != nil {
		return nil, err
	}

	return s, nil

}

// Pipeline for a crawler hooks with extractor functions.
// Might have prepopulated extractors or extractors added later.
// Empty by default.
func (s *Spider) Pipeline() *payload.Pipeline {
	return s.pipe
}

// WithPipeline sets up the extractors for the spider.
func (s *Spider) withPipeline() error {

	if s.pipe != nil {
		return nil
	}

	extractors, err := extract.Extractors(s.ExtractFields, s.ExtractEntities...)
	if err != nil {
		return err
	}

	s.pipe = payload.NewPipeline(extractors...)
	return nil
}

func (s *Spider) NewCrawler(deploy *Deploy) (*collect.Crawler, error) {

	if deploy.VictoriaLogsUrl != "" {
		VictoriaLogs(deploy.VictoriaLogsUrl, "info", s.Args.ID)
	}

	s.pipe.Starter(extract.WindmillMeta)

	if err := s.withMedia(deploy); err != nil {
		return nil, err
	}

	deps := &config.Deps{}

	if err := s.withDatabase(deploy, deps); err != nil {
		return nil, err
	}

	deps.Extractor = s.pipe.Extract

	// metrics
	if deploy.VictoriaMetricsUrl != "" {
		metrics, err := NewMetrics(vars.FromEnv().JobID, s.Args.ID, deploy.VictoriaMetricsUrl)
		if err != nil {
			return nil, err
		}
		deps.Monitor = metrics
	}

	// proxy
	if s.Args.ProxyEnabled {
		proxies, err := proxy.StartPool(s.Args.StartURL, s.Args.ProxySources...)
		if err != nil {
			return nil, err
		}
		deps.RoundTripper = proxies.Transport()
	}

	return collect.NewCrawler(s.Args, deps)
}

func (s *Spider) withMedia(deploy *Deploy) error {

	if !s.ExtractMedia {
		return nil
	}

	if deploy.Bucket.Name != "" {
		return nil
	}

	s3, err := store.NewS3Client(deploy.Bucket)

	if err != nil {
		return err
	}

	bucketStore := store.NewBucketStore(deploy.Bucket.Name, s3)

	path := fmt.Sprintf(SpiderMediaPath, s.Args.ID)
	publicURL := fmt.Sprintf("%s/%s", deploy.Bucket.PublicURL, path)
	uploader := media.NewMedia(publicURL, path, media.NewLoader(bucketStore))

	s.pipe.Starter(uploader.Claims)
	s.pipe.Finisher(uploader.Upload)

	return nil
}

func (s *Spider) withDatabase(deploy *Deploy, deps *config.Deps) (err error) {

	// database
	if deploy.MongoDSN == "" {
		return nil
	}

	deps.Storage, err = store.NewCollectStore(s.Args.ID, deploy.MongoDSN)
	if err != nil {
		return err
	}

	extractStore, err := store.NewExtractStore(s.Args.ID, deploy.MongoDSN)
	if err != nil {
		return fmt.Errorf("failed to create extract store: %w", err)
	}

	// provide save extractor func
	s.pipe.Finisher(extractStore.Save)

	return nil
}
