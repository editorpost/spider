package collect_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/html"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"spider/collect"
	"spider/storage"
	"testing"
)

func TestCollect(t *testing.T) {

	srv := ServeFile(t, "task_test.html")
	defer srv.Close()

	dispatched := false

	task := collect.Task{
		StartURL: srv.URL,
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article--ssr",
		Extract: func(*html.Node, *url.URL) error {
			dispatched = true
			return nil
		},
		Storage: storage.NewStorage("spider", "mongodb://localhost:27018"),
	}

	err := task.Start()
	require.NoError(t, err)
	assert.True(t, dispatched)
}

func TestJSCollect(t *testing.T) {

	srv := ServeFile(t, "task_test.html")
	defer srv.Close()

	dispatched := false

	task := collect.Task{
		StartURL: srv.URL,
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article--js",
		Extract: func(*html.Node, *url.URL) error {
			dispatched = true
			return nil
		},
	}

	require.NoError(t, task.Start())
	assert.True(t, dispatched)
}

func TestServeFile(t *testing.T) {

	srv := ServeFile(t, "task_test.html")
	defer srv.Close()

	// create a new request
	req, err := http.NewRequest("GET", srv.URL, nil)
	require.NoError(t, err)

	// create http client
	client := srv.Client()

	// send the request
	resp, err := client.Do(req)
	require.NoError(t, err)

	// check the response
	require.Equal(t, http.StatusOK, resp.StatusCode)

	// read the response body
	b, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalln(err)
	}

	body := string(b)

	// check the response body
	require.NotNil(t, body)

	// check html contains string "Hello, World!"
	require.Contains(t, body, "Hello, World!")
}

// ServeFile serves the file at the given path and returns the URL
func ServeFile(t *testing.T, path string) *httptest.Server {

	t.Helper()

	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}

	// create http server and serve the file

	srv, err := NewServer(b)
	if err != nil {
		t.Fatal(err)
	}

	// return the server URL
	return srv
}

// NewServer creates a new server that serves the given content
func NewServer(content []byte) (*httptest.Server, error) {

	// create a new server
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(content)
	}))

	return srv, nil
}
