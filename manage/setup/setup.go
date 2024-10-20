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
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/store"
	"github.com/google/uuid"
)

// Deploy provides the configuration for the Spider infrastructure.
type Deploy struct {
	Paths    store.Paths    `json:"Paths"`
	Storage  res.S3         `json:"Storage"`
	Media    res.S3Public   `json:"Media"`
	Database res.Postgresql `json:"Database"`
	Metrics  res.Metrics    `json:"Metrics"`
	Logs     res.Logs       `json:"Logs"`
}

// Spider aggregates configs and create collect.Crawler.
type Spider struct {
	ID       string
	Collect  *config.Config
	Extract  *extract.Config
	Deploy   *Deploy
	pipe     *pipe.Pipeline
	shutdown []func() error
}

func NewSpider(id string, args *config.Config, cfg *extract.Config, deploy *Deploy) (*Spider, error) {

	if id == "" {
		return nil, fmt.Errorf("spider ID is empty")
	}

	if _, err := uuid.Parse(id); err != nil {
		return nil, fmt.Errorf("spider ID is invalid: %w", err)
	}

	s := &Spider{
		ID:      id,
		Collect: args,
		Extract: cfg,
		Deploy:  deploy,
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

	extractors, err := extract.Extractors(s.Extract.Fields, s.Extract.Entities...)
	if err != nil {
		return err
	}

	s.pipe = pipe.NewPipeline(extractors...)
	return nil
}

func (s *Spider) NewCrawler() (*collect.Crawler, error) {

	// set the extractor function to the deps
	deps, err := s.NewDeps()
	if err != nil {
		return nil, err
	}

	return collect.NewCrawler(s.Collect, deps)
}

func (s *Spider) NewDeps() (*config.Deps, error) {

	s.withVictoriaLogs()

	deps := &config.Deps{}

	err := WithDepsFn(
		deps,
		s.withVictoriaMetrics,
		s.withProxy,
		s.withStorage,
	)

	if err != nil {
		return nil, err
	}

	// add windmill run metadata to the pipeline payload
	s.pipe.Starter(extract.WindmillMeta)

	// set the extractor function to the deps
	deps.Extractor = s.pipe.Extract

	return deps, nil
}

func (s *Spider) withVictoriaLogs() {

	if len(s.Deploy.Logs.URL) == 0 {
		return
	}

	VictoriaLogs(s.Deploy.Logs.URL, "info", s.ID)
}

func (s *Spider) withVictoriaMetrics(deps *config.Deps) (err error) {

	if len(s.Deploy.Metrics.URL) == 0 {
		return
	}

	deps.Monitor, err = NewMetrics(vars.FromEnv().JobID, s.ID, s.Deploy.Metrics.URL)
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

// Normalize the Spider configuration
func (s *Spider) Normalize() error {

	if s.ID == "" {
		return fmt.Errorf("spider ID is empty")
	}

	if _, err := uuid.Parse(s.ID); err != nil {
		return fmt.Errorf("spider ID is invalid: %w", err)
	}

	if s.Collect == nil {
		return fmt.Errorf("collect config is empty")
	}

	if s.Extract == nil {
		return fmt.Errorf("extract config is empty")
	}

	if err := s.Collect.Normalize(); err != nil {
		return err
	}

	if err := s.Extract.Normalize(); err != nil {
		return err
	}

	if err := s.withPipeline(); err != nil {
		return err
	}

	if s.Deploy == nil {
		s.Deploy = &Deploy{}
	}

	if s.Deploy.Paths.Collect == "" || s.Deploy.Paths.Payload == "" {
		s.Deploy.Paths = store.DefaultStoragePaths()
	}

	return nil
}

func SpiderFromJSON(data []byte) (*Spider, error) {

	s := &Spider{}
	if err := json.Unmarshal(data, s); err != nil {
		return nil, err
	}

	if err := s.Normalize(); err != nil {
		return nil, err
	}

	return s, nil
}

func NewDeploy(js string) (*Deploy, error) {

	deploy := &Deploy{}

	if err := json.Unmarshal([]byte(js), deploy); err != nil {
		return deploy, err
	}

	return deploy, nil
}
