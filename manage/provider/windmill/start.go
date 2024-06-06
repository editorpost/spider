package windmill

import (
	"github.com/editorpost/donq/pkg/script"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage"
)

// Start is an example code for running spider
// as Windmill Script with extract.Article
//
//goland:noinspection GoUnusedExportedFunction
func Start(argsMap any, extractors ...extract.PipeFn) error {

	args := &config.Args{}
	err := script.ParseArgs(argsMap, args)
	if err != nil {
		return err
	}

	deploy, err := DeployConfig(DefaultSpiderDeploy)
	if err != nil {
		return err
	}

	mongo, err := MongoConfig(DefaultMongoResource)
	if err != nil {
		return err
	}

	deploy.MongoDSN = mongo.DSN

	return manage.Start(args, deploy, extractors...)
}
