package media_test

import (
	"github.com/editorpost/spider/extract/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

// TestDownloadImage tests the DownloadImage function.
func TestDownload(t *testing.T) {
	// Set up a test server that serves an example image.
	testImage := []byte{0xFF, 0xD8, 0xFF} // Example of JPEG header bytes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(testImage)
	}))
	defer ts.Close()

	// Use the test server URL in place of the real image URL.
	transport := &http.Transport{}
	data, err := media.Download(ts.URL, transport)
	require.NoError(t, err)
	assert.Equal(t, testImage, data)
}
