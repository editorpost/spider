package windmill

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/manage/setup"
)

const (
	// DefaultMongoResource is the name of the mongo resource
	DefaultMongoResource = "f/spider/resource/mongodb"
	DefaultSpiderDeploy  = "f/spider/resource/deploy"
)

// MongoConfig returns the mongo config or panic
func MongoConfig(resource string) (*mongodb.Config, error) {

	if len(resource) == 0 {
		resource = DefaultMongoResource
	}

	cfg, err := mongodb.GetResource(resource)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

// SetupConfig returns the config or panic
func SetupConfig(deploy *setup.Config) error {

	if err := vars.FromResource(DefaultSpiderDeploy, deploy); err != nil {
		return err
	}

	mongo := &mongodb.Config{}
	if err := vars.FromResource(DefaultMongoResource, mongo); err != nil {
		return err
	}

	// vars.FromResource guarantees that deploy is not nil or error
	//goland:noinspection GoDfaNilDereference
	deploy.MongoDSN = mongo.DSN

	return nil
}
