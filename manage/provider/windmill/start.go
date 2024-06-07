package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage"
	"github.com/editorpost/spider/manage/setup"
)

// Start is an example code for running spider
// as Windmill Script with extract.Article
//
//goland:noinspection GoUnusedExportedFunction
func Start(argsJSON any, extractors ...extract.PipeFn) (err error) {

	args := &config.Args{}
	deploy := &setup.Config{}

	if err = vars.FromJSON(argsJSON, args); err != nil {
		return err
	}

	if err = SetupConfig(deploy); err != nil {
		return err
	}

	return manage.Start(args, deploy, extractors...)
}
