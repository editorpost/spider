package windmill_test

import (
	"encoding/json"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/provider/windmill"
	"github.com/editorpost/spider/manage/setup"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestTrial(t *testing.T) {

	s := NewSpider(t)
	extractors, err := extract.Extractors(s.ExtractFields, s.ExtractEntities...)
	require.NoError(t, err)
	require.NoError(t, windmill.Trial(s.Args, extractors...))
}

func NewSpider(t *testing.T) *setup.Spider {
	s := &setup.Spider{}
	require.NoError(t, json.Unmarshal([]byte(spiderJSon), s))
	return s
}

const spiderJSon = `
{
    "ID": "66781dbbc739b522cfff973e",
    "Name": "Dold",
    "Depth": 1,
    "Created": "2024-06-23 13:06:03.299508323 +0000 UTC",
    "Updated": "2024-06-24 07:57:04.81246902 +0000 UTC",
    "StartURL": "https://www.dold.com/en/products/relay-modules/monitoring-devices/residual-current-monitors/",
    "ProjectID": "6676117acd41b43c99836979",
    "UserAgent": "",
    "AllowedURL": "https://www.dold.com/en/products/relay-modules/monitoring-devices/residual-current-monitors",
    "ExtractURL": "https://www.dold.com/en/products/relay-modules/monitoring-devices/residual-current-monitors/{some}",
    "UseBrowser": false,
    "ExtractLimit": 3,
    "ProxyEnabled": false,
    "ProxySources": [],
    "ExtractFields": [
        {
            "ID": "ea2a1a5b-e4a9-4be0-bc14-9f6f9b860b0a",
            "Name": "model",
            "Scoped": false,
            "Children": [],
            "ParentID": "",
            "Required": true,
            "Selector": "h1",
            "Multiline": false,
            "BetweenEnd": "",
            "FinalRegex": "",
            "Cardinality": 1,
            "InputFormat": "html",
            "BetweenStart": "",
            "OutputFormat": [
                "text"
            ]
        }
    ],
    "ExtractEntities": [],
    "ExtractSelector": ".product--details"
}
`
