package manage

import (
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract/payload"
	"github.com/editorpost/spider/manage/setup"
)

// Start is a code for running spider
// as Windmill Script with extract.Article
func Start(args *config.Args, deploy *setup.Config, extractors ...payload.Extractor) error {

	c, err := Crawler(args, deploy, extractors...)

	if err != nil {
		return err
	}

	return c.Run()
}

func Crawler(args *config.Args, deploy *setup.Config, extractors ...payload.Extractor) (*collect.Crawler, error) {

	deps, err := setup.Deps(args, deploy, extractors...)
	if err != nil {
		return nil, err
	}

	crawler, err := collect.NewCrawler(args, deps)
	if err != nil {
		return nil, err
	}

	return crawler, nil
}
