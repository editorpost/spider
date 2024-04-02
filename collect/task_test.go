package collect_test

import (
	"github.com/PuerkitoBio/goquery"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"spider/collect"
	"spider/extract"
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
		Extract: func(*goquery.Selection, *url.URL) error {
			dispatched = true
			return nil
		},
		Storage: storage.NewCollectorStore("spider", "mongodb://localhost:27018"),
	}

	err := task.Start()
	require.NoError(t, err)
	assert.True(t, dispatched)
}

func TestSave(t *testing.T) {

	srv := ServeFile(t, "task_test.html")
	defer srv.Close()

	dispatched := false

	// todo replace task id
	_ = os.Setenv("CRAWLAB_TASK_ID", "65f3b30ffc2e89f5abbd078a")
	_ = os.Setenv("CRAWLAB_MONGO_HOST", "localhost")
	_ = os.Setenv("CRAWLAB_MONGO_PORT", "27018")
	_ = os.Setenv("CRAWLAB_MONGO_DATABASE", "spider")
	_ = os.Setenv("CRAWLAB_COLLECTION", "results_test")
	_ = os.Setenv("CRAWLAB_MONGO_USERNAME", "root")
	_ = os.Setenv("CRAWLAB_MONGO_PASSWORD", "nopass")
	_ = os.Setenv("CRAWLAB_MONGO_AUTHSOURCE", "admin")
	// _ = os.Setenv("CRAWLAB_DATA_SOURCE", "mongodb://root:nopass@mongo:27018/crawlab?authSource=admin")
	_ = os.Setenv("CRAWLAB_NODE_MASTER", "Y")
	_ = os.Setenv("CRAWLAB_MONGO_DB", "crawlab")

	// leave empty to use source from env
	// _ = os.Setenv("CRAWLAB_DATA_SOURCE", "")
	_ = os.Setenv("CRAWLAB_DATA_SOURCE", "mongodb://root:nopass@localhost:27018/spider?authSource=admin")

	dispatcher := func(payload *extract.Payload) error {
		dispatched = true
		return nil
	}

	collectorStore := storage.NewCollectorStore("spider", "mongodb://localhost:27018")
	err := collectorStore.Init()
	require.NoError(t, err)

	extractorStore := storage.NewExtractStoreFromEnv()

	saveExtracted := func(p *extract.Payload) error {

		p.Data["_tid"] = os.Getenv("CRAWLAB_TASK_ID") // todo: investigate why used explicit value. Testing?
		p.Data["url"] = p.URL
		p.Data["html"], err = p.Selection.Html()
		if err != nil {
			return err
		}

		return extractorStore.Save(p.Data)
	}

	task := collect.Task{
		StartURL: srv.URL,
		MatchURL: ".*",
		Depth:    1,
		Query:    ".article--ssr",
		Extract:  extract.Pipe(extract.Article, saveExtracted, dispatcher),
		Storage:  collectorStore,
	}

	err = task.Start()
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
		Extract: func(*goquery.Selection, *url.URL) error {
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
