package extract

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/manage"
)

// WindmillMeta is a meta data extractor
func WindmillMeta(p *Payload) error {

	e := vars.FromEnv()
	// windmill flow
	p.Data[manage.WindmillJobID] = e.GetRootFlowJobID()
	p.Data[manage.WindmillFlowPath] = e.GetFlowPath()
	p.Data[manage.WindmillFlowJobID] = e.GetFlowJobID()

	return nil
}
