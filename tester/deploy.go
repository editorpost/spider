package tester

import (
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/store"
	"testing"
)

func TestDeploy(t *testing.T) *setup.Deploy {

	t.Helper()

	return &setup.Deploy{
		Bucket: store.Bucket{
			Name:      "local",
			Endpoint:  "./testdata",
			PublicURL: "http://localhost:9000",
		},
		VictoriaMetricsUrl: "http://metrics.spider.svc:8429/api/v1/import/prometheus",
		VictoriaLogsUrl:    "http://logs.spider.svc:9428",
	}
}
