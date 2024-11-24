package tester_test

import (
	"encoding/json"
	"github.com/editorpost/spider/manage/provider/windmill"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestNewSpider(t *testing.T) {

	srv := tester.NewServer("./fixtures")
	defer srv.Close()

	s := tester.NewSpiderWith(t, srv)
	require.NotNil(t, s)

	// check spider
	require.NoError(t, windmill.Check(s))

	// read results.json
	f, err := os.ReadFile(windmill.JobResultFile)
	require.NoError(t, err)
	require.NotNil(t, f)

	// unmarshal to map
	var data struct{ CheckID string }
	require.NoError(t, json.Unmarshal(f, &data))
	assert.NotEmpty(t, data.CheckID)

	// remove results.json
	require.NoError(t, os.Remove(windmill.JobResultFile))
}
