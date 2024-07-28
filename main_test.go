package main

import (
	"encoding/json"
	"flag"
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFlags(t *testing.T) {

	jsonStr := `{
		"Collect": {
			"ID": "4c62e925-2e3e-4ff6-a92b-305b85d6281a",
			"StartURL": "http://example.com",
			"ExtractLimit": 10
		},
		"Extract": {
			"Extract": ["person", "organization"],
			"Fields": [
				{
					"Name": "name",
					"Selector": "h1"
				}
			],
			"Media": {
				"Enabled": true
			}
		}
	}`

	s, err := setup.NewSpiderFromJSON([]byte(jsonStr))
	require.NoError(t, err)
	require.NoError(t, json.Unmarshal([]byte(jsonStr), s))

	require.NoError(t, flag.Set("spider", jsonStr))
	require.NoError(t, flag.Set("deploy", deployStr))
	require.NoError(t, flag.Set("cmd", "start"))

	cmd, spider, deploy := Flags()

	assert.Equal(t, "start", cmd)
	assert.Equal(t, s.Collect, spider.Collect)
	assert.Equal(t, s.Extract, spider.Extract)
	assert.Equal(t, "local", deploy.Storage.Bucket)
	assert.Equal(t, "apac", deploy.Storage.Region)
}

func TestDeployFromString(t *testing.T) {

	t.Helper()
	deploy, err := setup.NewDeploy(deployStr)
	require.NoError(t, err)

	assert.Equal(t, "local", deploy.Storage.Bucket)
	assert.Equal(t, "apac", deploy.Storage.Region)

	// media
	assert.Equal(t, true, deploy.Storage.UseSSL)
	assert.Equal(t, "./testdata", deploy.Storage.EndPoint)
	assert.Equal(t, "media-access-key", deploy.Media.AccessKey)
	assert.Equal(t, false, deploy.Media.PathStyle)
	assert.Equal(t, "https://cdn.example.com", deploy.Media.PublicURL)
	assert.Equal(t, "storage-secret-key", deploy.Storage.SecretKey)

	// storage
	assert.Equal(t, "storage-access-key", deploy.Storage.AccessKey)
	assert.Equal(t, "storage-secret-key", deploy.Storage.SecretKey)

	// database
	assert.Equal(t, "primary.postgres.svc", deploy.Database.Host)
	assert.Equal(t, 5432, deploy.Database.Port)
	assert.Equal(t, "testdb", deploy.Database.Dbname)
	assert.Equal(t, "", deploy.Database.SSLMode)
	assert.Equal(t, "", deploy.Database.Password)
	assert.Equal(t, "", deploy.Database.RootCertificatePEM)

	// logs
	assert.Equal(t, "", deploy.Logs.URL)

	// metrics
	assert.Equal(t, "http://metrics.spider.svc:8429/api/v1/import/prometheus", deploy.Metrics.URL)
}

var deployStr = `
{
  "Media": {
    "bucket": "local",
    "region": "apac",
    "useSSL": true,
    "endPoint": "` + tester.DeployBucketFolder + `",
    "accessKey": "media-access-key",
    "pathStyle": false,
    "publicURL": "https://cdn.example.com",
    "secretKey": "media-secret-key"
  },
  "Storage": {
    "bucket": "local",
    "region": "apac",
    "useSSL": true,
    "endPoint": "` + tester.DeployBucketFolder + `",
    "accessKey": "storage-access-key",
    "pathStyle": false,
    "secretKey": "storage-secret-key"
  },
  "Database": {
    "host": "primary.postgres.svc",
    "port": 5432,
    "user": "",
    "dbname": "testdb",
    "sslmode": "",
    "password": "",
    "root_certificate_pem": ""
  },
  "Logs": {
    "url": ""
  },
  "Metrics": {
    "url": "http://metrics.spider.svc:8429/api/v1/import/prometheus"
  }
}
`
