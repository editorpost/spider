package extract

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/extract/pipe"
)

//goland:noinspection GoUnusedConst
const (
	// WindmillJobID is the key for the job ID
	WindmillJobID = "windmill__job_id"
	// WindmillFlowPath is the key for the flow path
	WindmillFlowPath = "windmill__flow_path"
	// WindmillFlowJobID is the key for the flow job ID
	WindmillFlowJobID = "windmill__flow_job_id"
	// WindmillJobPath is the key for the job path
	WindmillJobPath = "windmill__job_path"
)

// WindmillMeta is a meta data extractor
func WindmillMeta(p *pipe.Payload) error {
	e := vars.FromEnv()
	p.Data[WindmillJobID] = e.GetJobID()
	p.Data[WindmillJobPath] = e.GetJobPath()
	p.Data[WindmillFlowPath] = e.GetFlowPath()
	p.Data[WindmillFlowJobID] = e.GetFlowJobID()
	return nil
}
