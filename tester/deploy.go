package tester

import (
	"github.com/editorpost/donq/res"
	"github.com/editorpost/spider/manage/setup"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

const DeployBucketFolder = "./testdata"

func TestDeploy(t *testing.T) setup.Deploy {

	t.Helper()

	return setup.Deploy{
		Storage: res.S3{
			Bucket:   "local",
			EndPoint: DeployBucketFolder,
		},
		Media: res.S3Public{
			S3: res.S3{
				Bucket:   "local",
				EndPoint: DeployBucketFolder,
			},
			PublicURL: "https://cdn.example.com",
		},
		Metrics: res.Metrics{
			URL: "",
		},
		Logs: res.Logs{
			URL: "",
		},
	}
}

func CleanTestBucket(t *testing.T) {
	t.Helper()
	require.NoError(t, os.RemoveAll(DeployBucketFolder))
}
