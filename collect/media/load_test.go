package media_test

import (
	"github.com/editorpost/spider/collect/media"
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

// TestDownloadImage tests the DownloadImage function.
func TestDownloader_Download(t *testing.T) {
	// Set up a test server that serves an example image.
	testImage := []byte{0xFF, 0xD8, 0xFF} // Example of JPEG header bytes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(testImage)
	}))
	defer ts.Close()

	// Use the test server URL in place of the real image URL.
	downloader := media.NewDownloader(&http.Transport{})

	buf, err := downloader.Download(ts.URL)
	require.NoError(t, err)
	defer downloader.ReleaseBuffer(buf)
	assert.Equal(t, string(buf.Bytes()), string(testImage))
}

func TestDownloader_BuffersMemoryAllocation(t *testing.T) {
	// Set up a test server that serves an example image.
	data := []byte{0xFF, 0xD8, 0xFF} // Example of JPEG header bytes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(data)
	}))
	defer ts.Close()

	// Use the test server URL in place of the real image URL.
	downloader := media.NewDownloader(&http.Transport{})

	// Download the image multiple times to ensure the buffers are reused.
	for i := 0; i < 10; i++ {
		buf, err := downloader.Download(ts.URL)
		require.NoError(t, err)
		defer downloader.ReleaseBuffer(buf)
		assert.Equal(t, string(buf.Bytes()), string(data))
	}
}
