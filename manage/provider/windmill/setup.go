package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/manage/setup"
)

const (
	DefaultSpiderDeploy = "f/spider/resource/deploy"
)

// LoadDeployResource returns the config or panic
func LoadDeployResource(deploy *setup.Deploy) error {

	if err := vars.FromResource(DefaultSpiderDeploy, deploy); err != nil {
		return err
	}
	return nil
}
