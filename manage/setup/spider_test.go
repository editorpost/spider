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
		"ID": "test",
		"StartURL": "http://example.com",
		"ExtractLimit": 10,
		"ExtractEntities": ["person", "organization"],
		"ExtractFields": [
			{
				"Name": "name",
				"Selector": "h1"
			}
		]
	}`

	s := &setup.Spider{}
	require.NoError(t, json.Unmarshal([]byte(jsonStr), s))
	require.Equal(t, "test", s.ID)
	require.Equal(t, "http://example.com", s.StartURL)
	require.Equal(t, 10, s.ExtractLimit)
	require.Len(t, s.ExtractEntities, 2)
	require.Len(t, s.ExtractFields, 1)
}
