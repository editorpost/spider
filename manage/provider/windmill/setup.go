package windmill

import (
	"encoding/json"
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/manage/setup"
	wmill "github.com/windmill-labs/windmill-go-client"
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

// DeployConfig returns the deployment config or panic
func DeployConfig(resource string) (*setup.Deploy, error) {

	if len(resource) == 0 {
		resource = DefaultSpiderDeploy
	}

	data, err := wmill.GetResource(resource)
	if err != nil {
		return nil, err
	}

	b, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	deploy := &setup.Deploy{}
	err = json.Unmarshal(b, deploy)
	if err != nil {
		return nil, err
	}

	return deploy, nil
}
