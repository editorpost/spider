package collect_test

import (
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"spyder/collect"
	"testing"
)

func TestRun(t *testing.T) {

	srv := ServeFile(t, "colly_test.html")
	defer srv.Close()

	extracted := false

	task := collect.Task{
		StartURL: srv.URL,
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article",
		Extract: func(e *colly.HTMLElement) error {
			println(e.Text)
			extracted = true
			return nil
		},
	}
	err := collect.Run(task)
	require.NoError(t, err)
	assert.True(t, extracted)
}

func TestServeFile(t *testing.T) {

	srv := ServeFile(t, "colly_test.html")
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

	html := string(b)

	// check the response body
	require.NotNil(t, html)

	// check html contains string "Hello, World!"
	require.Contains(t, html, "Hello, World!")
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
