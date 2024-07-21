package tester

import (
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestServer struct {
	*httptest.Server
}

// NewServer creates a new server that serves the given content
func NewServer(dir string) *TestServer {
	srv := httptest.NewServer(http.StripPrefix("/", http.FileServer(http.Dir(dir))))
	return &TestServer{srv}
}

func (srv *TestServer) HTML(t *testing.T, path string) string {

	resp, err := http.Get(srv.URL + "/" + path)
	require.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	body := make([]byte, 0)
	_, err = resp.Body.Read(body)
	require.NoError(t, err)

	return string(body)
}

func (srv *TestServer) Document(t *testing.T, path string) *colly.HTMLElement {
	return GetDocument(t, srv.HTML(t, path))
}
