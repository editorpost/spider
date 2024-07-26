package setup_test

import (
	"encoding/json"
	"github.com/editorpost/spider/manage/setup"
	"github.com/stretchr/testify/require"
	"testing"
)

// TestSpiderFromJSONString is the test for Spider unmarshalling from JSON string
func TestSpiderFromJSONString(t *testing.T) {

	jsonStr := `{
		"Collect": {
			"ID": "test",
			"StartURL": "http://example.com",
			"ExtractLimit": 10
		},
		"Extract": {
			"ExtractEntities": ["person", "organization"],
			"ExtractFields": [
				{
					"Name": "name",
					"Selector": "h1"
				}
			]
		}
	}`

	s := &setup.Spider{}
	require.NoError(t, json.Unmarshal([]byte(jsonStr), s))
	require.Equal(t, "test", s.Collect.ID)
	require.Equal(t, "http://example.com", s.Collect.StartURL)
	require.Equal(t, 10, s.Collect.ExtractLimit)
	require.Len(t, s.Extract.ExtractEntities, 2)
	require.Len(t, s.Extract.ExtractFields, 1)
}
