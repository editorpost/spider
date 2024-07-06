package windmill

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/manage/setup"
)

const (
	DefaultSpiderDeploy = "f/spider/resource/deploy"
)

// DeployResource returns the config or panic
func DeployResource(deploy *setup.Deploy) error {

	if err := vars.FromResource(DefaultSpiderDeploy, deploy); err != nil {
		return err
	}
	return nil
}
