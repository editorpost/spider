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
	downloader := media.NewDownloader(nil, &http.Transport{})

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
	downloader := media.NewDownloader(nil, &http.Transport{})

	// Download the image multiple times to ensure the buffers are reused.
	for i := 0; i < 10; i++ {
		buf, err := downloader.Download(ts.URL)
		require.NoError(t, err)
		defer downloader.ReleaseBuffer(buf)
		assert.Equal(t, string(buf.Bytes()), string(data))
	}
}

func BenchmarkDownloader_Download(b *testing.B) {
	// Set up a test server that serves an example image.
	testImage := []byte{0xFF, 0xD8, 0xFF} // Example of JPEG header bytes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(testImage)
	}))
	defer ts.Close()

	// Use the test server URL in place of the real image URL.
	downloader := media.NewDownloader(nil, &http.Transport{})

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := downloader.Download(ts.URL)
		require.NoError(b, err)
	}
}

func TestDownloader_Copy(t *testing.T) {
	// Set up a test server that serves an example image.
	testImage := []byte{0xFF, 0xD8, 0xFF} // Example of JPEG header bytes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(testImage)
	}))
	defer ts.Close()

	// Use the test server URL in place of the real image URL.
	storage := NewMockStorage()
	downloader := media.NewDownloader(storage, &http.Transport{})

	// Perform the download and upload.
	_, err := downloader.Upload(ts.URL)
	require.NoError(t, err)

	// Generate the expected upload path.
	uploadPath, err := media.StorageHash(ts.URL)
	require.NoError(t, err)

	// Assert the data was uploaded correctly.
	uploadedData, exists := storage.data[uploadPath]
	require.True(t, exists)
	assert.Equal(t, testImage, uploadedData)
}

// MockStorage is a mock implementation of the Store interface for testing.
type MockStorage struct {
	data map[string][]byte
}

// NewMockStorage creates a new MockStorage instance.
func NewMockStorage() *MockStorage {
	return &MockStorage{
		data: make(map[string][]byte),
	}
}

// Upload mocks the upload of data to a storage system.
func (ms *MockStorage) Upload(data []byte, path string) error {
	if ms.data == nil {
		ms.data = make(map[string][]byte)
	}
	ms.data[path] = data
	return nil
}
