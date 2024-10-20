package windmill_test

import (
	"encoding/json"
	"github.com/editorpost/spider/manage/provider/windmill"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestTrial(t *testing.T) {

	srv := tester.NewServer("../../../tester/fixtures")
	defer srv.Close()

	s := tester.NewSpiderWith(t, srv)
	require.NotNil(t, s)

	// set env vars
	require.NoError(t, os.Setenv("WM_JOB_ID", "test-job-id"))
	require.NoError(t, windmill.Check(s))

	// read results.json
	f, err := os.ReadFile(windmill.JobResultFile)
	require.NoError(t, err)
	require.NotNil(t, f)

	// unmarshal to map
	var data struct{ ID string }
	require.NoError(t, json.Unmarshal(f, &data))
	assert.Equal(t, data.ID, "test-job-id") // UUID

	// remove results.json
	require.NoError(t, os.Remove(windmill.JobResultFile))
}
