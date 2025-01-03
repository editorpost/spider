package console_test

import (
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidate(t *testing.T) {

	srv := tester.NewServer("../../tester/fixtures")
	defer srv.Close()

	spider := tester.NewSpiderWith(t, srv)
	require.NotNil(t, spider)

	spider.Deploy = tester.TestDeploy(t)
	err := console.Validate(spider)
	require.NoError(t, err)
}
