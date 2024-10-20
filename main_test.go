package main

import (
	"flag"
	"github.com/editorpost/spider/manage/setup"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestFlags(t *testing.T) {

	require.NoError(t, flag.Set("spider", validSpider))
	require.NoError(t, flag.Set("cmd", "start"))

	cmd, s, _ := Flags()

	assert.Equal(t, "start", cmd)

	// collect values assertion:
	collect := s.Collect
	assert.Equal(t, 1, collect.Depth)
	assert.Equal(t, "https://thailand-news.ru/news/turizm/", collect.StartURL)
	assert.Equal(t, "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/58.0.3029.110 Safari/537.3", collect.UserAgent)
	assert.Equal(t, "https://thailand-news.ru/news/{turizm,puteshestviya}{any}", collect.AllowedURL)
	assert.Equal(t, "https://thailand-news.ru/news/{turizm,puteshestviya}/{some}", collect.ExtractURL)
	assert.Equal(t, false, collect.UseBrowser)
	assert.Equal(t, 3, collect.ExtractLimit)
	assert.Equal(t, false, collect.ProxyEnabled)
	assert.Len(t, collect.ProxySources, 0)
	assert.Equal(t, ".node-article--full", collect.ExtractSelector)

	// extract values assertion:
	extract := s.Extract
	assert.Equal(t, false, extract.Media.Enabled)
	assert.Len(t, extract.Fields, 0)
	assert.Len(t, extract.Entities, 1)

	// spider values assertion:
	assert.Equal(t, "df265b45-00bc-4aa6-bad2-a83018ff42ca", s.ID)

	// deploy
	deploy := s.Deploy

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

func TestInvalidSpider(t *testing.T) {

	_, err := setup.SpiderFromJSON([]byte(validSpider))
	require.NoError(t, err)
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

// the code from real case for v0.3.3 version
var validSpider = `{
  "ID": "df265b45-00bc-4aa6-bad2-a83018ff42ca",
  "Name": "Tourism",
  "Tags": [],
  "edges": {},
  "Collect": {
    "Depth": 1,
    "StartURL": "https://thailand-news.ru/news/turizm/",
    "Scheduled": false,
    "Schedules": null,
    "UserAgent": "",
    "AllowedURL": "https://thailand-news.ru/news/{turizm,puteshestviya}{any}",
    "ExtractURL": "https://thailand-news.ru/news/{turizm,puteshestviya}/{some}",
    "UseBrowser": false,
    "ExtractLimit": 3,
    "ProxyEnabled": false,
    "ProxySources": [],
    "ExtractSelector": ".node-article--full"
  },
  "Created": "2024-08-04T16:00:28.096157Z",
  "Extract": {
    "Media": {
      "Enabled": false
    },
    "Fields": [],
    "Entities": [
      "article"
    ]
  },
  "Updated": "2024-08-04T20:12:27.096488Z",
  "ProjectID": "13b4ab82-bf02-40dd-a27e-d00682062872",
  "Deploy": {
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
}`
