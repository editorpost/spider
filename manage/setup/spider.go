package setup

import (
	"encoding/json"
	"fmt"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/donq/res"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/collect/proxy"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/store"
	"github.com/google/uuid"
)

// Deploy provides the configuration for the Spider infrastructure.
type Deploy struct {
	Storage  res.S3         `json:"Storage"`
	Media    res.S3Public   `json:"Media"`
	Database res.Postgresql `json:"Database"`
	Metrics  res.Metrics    `json:"Metrics"`
	Logs     res.Logs       `json:"Logs"`
}

// Spider aggregates configs and create collect.Crawler.
type Spider struct {
	Collect  *config.Config
	Extract  *extract.Config
	pipe     *pipe.Pipeline
	shutdown []func() error
}

func NewSpider(args *config.Config, cfg *extract.Config) (*Spider, error) {

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

	if err = s.Collect.Normalize(); err != nil {
		return nil, err
	}

	if err = s.withPipeline(); err != nil {
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

	extractors, err := extract.Extractors(s.Extract.Fields, s.Extract.Extract...)
	if err != nil {
		return err
	}

	s.pipe = pipe.NewPipeline(extractors...)
	return nil
}

func (s *Spider) NewCrawler(deploy Deploy) (*collect.Crawler, error) {

	deps := &config.Deps{}

	s.withVictoriaLogs(deploy.Logs)

	if err := s.withVictoriaMetrics(deploy.Metrics, deps); err != nil {
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

func (s *Spider) withVictoriaLogs(config res.Logs) {

	if len(config.URL) == 0 {
		return
	}

	VictoriaLogs(config.URL, "info", s.Collect.ID)
}

func (s *Spider) withVictoriaMetrics(config res.Metrics, deps *config.Deps) (err error) {

	if len(config.URL) == 0 {
		return
	}

	deps.Monitor, err = NewMetrics(vars.FromEnv().JobID, s.Collect.ID, config.URL)
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

func (s *Spider) withStorage(deploy Deploy, deps *config.Deps) error {

	if deploy.Storage.Bucket == "" {
		return nil
	}

	if err := s.withCollectStore(deploy.Storage, deps); err != nil {
		return err
	}

	if err := s.withExtractStore(deploy.Storage); err != nil {
		return err
	}

	if err := s.withMedia(deploy.Media); err != nil {
		return err
	}

	return nil
}

func (s *Spider) withCollectStore(bucket res.S3, deps *config.Deps) error {

	storage, upload, err := store.NewCollectStorage(s.Collect.ID, bucket)
	if err != nil {
		return err
	}

	// upload visited urls to S3
	s.onShutdown(upload)
	deps.Storage = storage

	return err
}

func (s *Spider) withExtractStore(bucket res.S3) (err error) {

	extractStore, err := store.NewExtractStorage(s.Collect.ID, bucket)
	if err != nil {
		return fmt.Errorf("failed to create extract S3 storage: %w", err)
	}

	// provide save extractor func
	s.pipe.Finisher(extractStore.Save)

	return nil
}

func (s *Spider) withMedia(bucket res.S3Public) error {

	if !s.Extract.Media.Enabled {
		return nil
	}

	bucketStore, err := store.NewMediaStorage(s.Collect.ID, bucket.S3)
	if err != nil {
		return err
	}

	// public url prefix for media files, e.g. http://my-proxy:8080
	// join public url with bucket folder, e.g. spider/%/media/123.jpg
	// to simplify further proxying the bucket, e.g. http://my-proxy:8080/spider/%/media/123.jpg
	folder := store.GetMediaStorageFolder(s.Collect.ID)
	publicURL := fmt.Sprintf("%s/%s", bucket.PublicURL, folder)
	uploader := media.NewMedia(publicURL, media.NewLoader(bucketStore))

	s.pipe.Starter(uploader.Claims)
	s.pipe.Finisher(uploader.Upload)

	return nil
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

func NewDeploy(js string) (Deploy, error) {

	deploy := Deploy{}

	if err := json.Unmarshal([]byte(js), &deploy); err != nil {
		return deploy, err
	}

	return deploy, nil
}
