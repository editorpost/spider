package manage

import (
	"github.com/editorpost/donq/pkg/script"
	"github.com/editorpost/spider/collect"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/setup"
)

// StartWith is an example code for running spider
// as Windmill Script with extract.Article
//
//goland:noinspection GoUnusedExportedFunction
func StartWith(input any) error {

	args := &collect.Args{}
	if err := script.ParseArgs(input, args); err != nil {
		return err
	}

	return Start(args)
}

// Start is a code for running spider
// as Windmill Script with extract.Article
func Start(args *collect.Args) error {

	deploy := &setup.Deploy{
		SpiderID: "ready-check",
	}

	extractor := func(*extract.Payload) error {
		return nil
	}

	deps, err := setup.Deps(args, deploy, extractor)
	if err != nil {
		return err
	}

	crawler, err := collect.NewCrawler(args, deps)
	if err != nil {
		return err
	}

	return crawler.Run()
}
