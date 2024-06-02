package extract

import (
	"github.com/editorpost/donq/pkg/vars"
)

const (
	// DefaultMongoResource is the name of the mongo resource
	DefaultMongoResource = "f/spider/resource/mongodb"
	// WindmillJobID is the key for the job ID
	WindmillJobID = "windmill__job_id"
	// WindmillFlowPath is the key for the flow path
	WindmillFlowPath = "windmill__flow_path"
	// WindmillFlowJobID is the key for the flow job ID
	WindmillFlowJobID = "windmill__flow_job_id"
)

// WindmillMeta is a meta data extractor
func WindmillMeta(p *Payload) error {

	e := vars.FromEnv()
	// windmill flow
	p.Data[WindmillJobID] = e.GetRootFlowJobID()
	p.Data[WindmillFlowPath] = e.GetFlowPath()
	p.Data[WindmillFlowJobID] = e.GetFlowJobID()

	return nil
}
