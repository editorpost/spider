package tester_test

import (
	"github.com/editorpost/spider/tester"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestNewServer(t *testing.T) {
	srv := tester.NewServer("./fixtures")
	defer srv.Close()

	// make request to the server /article.html and check the response
	resp, err := http.Get(srv.URL + "/article.html")
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// make request to the server /not-found.html and check the response
	resp, err = http.Get(srv.URL + "/not-found.html")
	require.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
