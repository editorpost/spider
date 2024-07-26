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
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/store"
	"github.com/google/uuid"
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
	// VictoriaMetricsUrl
	VictoriaMetricsUrl string `json:"VictoriaMetricsUrl" validate:"trim"`
	// VictoriaLogsUrl
	VictoriaLogsUrl string `json:"VictoriaLogsUrl" validate:"trim"`
}

type Spider struct {
	Collect  *config.Args
	Extract  *extract.Config
	pipe     *pipe.Pipeline
	shutdown []func() error
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

	return NewSpider(s.Collect, s.Extract)
}

func NewSpider(args *config.Args, cfg *extract.Config) (*Spider, error) {

	if args.ID == "" {
		return nil, fmt.Errorf("spider ID is empty")
	}

	_, err := uuid.Parse(args.ID)
	if err != nil {
		return nil, fmt.Errorf("spider ID is invalid: %w", err)
	}

	s := &Spider{
		Collect: args,
		Extract: cfg,
	}

	if err := s.Collect.Normalize(); err != nil {
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
func (s *Spider) Pipeline() *pipe.Pipeline {
	return s.pipe
}

// WithPipeline sets up the extractors for the spider.
func (s *Spider) withPipeline() error {

	if s.pipe != nil {
		return nil
	}

	extractors, err := extract.Extractors(s.Extract.ExtractFields, s.Extract.ExtractEntities...)
	if err != nil {
		return err
	}

	s.pipe = pipe.NewPipeline(extractors...)
	return nil
}

func (s *Spider) NewCrawler(deploy *Deploy) (*collect.Crawler, error) {

	deps := &config.Deps{}

	s.withVictoriaLogs(deploy.VictoriaLogsUrl)

	if err := s.withVictoriaMetrics(deploy.VictoriaMetricsUrl, deps); err != nil {
		return nil, err
	}

	if err := s.withProxy(deps); err != nil {
		return nil, err
	}

	if err := s.withStorage(deploy, deps); err != nil {
		return nil, err
	}

	s.pipe.Starter(extract.WindmillMeta)
	deps.Extractor = s.pipe.Extract

	return collect.NewCrawler(s.Collect, deps)
}

func (s *Spider) withVictoriaLogs(uri string) {
	if uri != "" {
		VictoriaLogs(uri, "info", s.Collect.ID)
	}
}

func (s *Spider) withVictoriaMetrics(uri string, deps *config.Deps) (err error) {

	if len(uri) == 0 {
		return nil
	}

	deps.Monitor, err = NewMetrics(vars.FromEnv().JobID, s.Collect.ID, uri)
	return err
}

func (s *Spider) withProxy(deps *config.Deps) error {

	if !s.Collect.ProxyEnabled {
		return nil
	}

	proxies, err := proxy.StartPool(s.Collect.StartURL, s.Collect.ProxySources...)
	if err != nil {
		return err
	}

	deps.RoundTripper = proxies.Transport()

	return nil
}

func (s *Spider) withStorage(deploy *Deploy, deps *config.Deps) error {

	// s3
	if deploy.Bucket.Name == "" {
		return nil
	}

	if err := s.withCollectStore(deploy, deps); err != nil {
		return err
	}

	if err := s.withExtractStore(deploy); err != nil {
		return err
	}

	if err := s.withMedia(deploy); err != nil {
		return err
	}

	return nil
}

func (s *Spider) withCollectStore(deploy *Deploy, deps *config.Deps) error {

	storage, upload, err := store.NewCollectStorage(s.Collect.ID, deploy.Bucket)
	if err != nil {
		return err
	}

	// upload visited urls to S3
	s.onShutdown(upload)
	deps.Storage = storage

	return err
}

func (s *Spider) withExtractStore(deploy *Deploy) (err error) {

	extractStore, err := store.NewExtractStorage(s.Collect.ID, deploy.Bucket)
	if err != nil {
		return fmt.Errorf("failed to create extract S3 storage: %w", err)
	}

	// provide save extractor func
	s.pipe.Finisher(extractStore.Save)

	return nil
}

func (s *Spider) withMedia(deploy *Deploy) error {

	if !s.Extract.ExtractMedia {
		return nil
	}

	bucketStore, err := store.NewMediaStorage(s.Collect.ID, deploy.Bucket)
	if err != nil {
		return err
	}

	// public url prefix for media files, e.g. http://my-proxy:8080
	// join public url with bucket folder, e.g. spider/%/media/123.jpg
	// to simplify further proxying the bucket, e.g. http://my-proxy:8080/spider/%/media/123.jpg
	folder := store.GetMediaStorageFolder(s.Collect.ID)
	publicURL := fmt.Sprintf("%s/%s", deploy.Bucket.PublicURL, folder)
	uploader := media.NewMedia(publicURL, media.NewLoader(bucketStore))

	s.pipe.Starter(uploader.Claims)
	s.pipe.Finisher(uploader.Upload)

	return nil
}
