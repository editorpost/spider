package media_test

import (
	"github.com/editorpost/spider/extract/media"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

var (
	server *httptest.Server
)

func DataExpected() []byte {
	return []byte{0xFF, 0xD8, 0xFF}
}

func DataAssert(t *testing.T, got []byte) {
	assert.Equal(t, DataExpected(), got)
}

func TestMain(m *testing.M) {

	// run image server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write(DataExpected()) // Example of JPEG header bytes.
	}))
	defer server.Close()
	m.Run()
}

// TestDownloadImage tests the DownloadImage function.
func TestDownload(t *testing.T) {
	data, err := media.Download(server.URL, &http.Transport{})
	require.NoError(t, err)
	DataAssert(t, data)
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
	downloader := media.NewLoader(nil)
	downloader.SetClient(ts.Client())

	buf, err := downloader.Fetch(ts.URL)
	require.NoError(t, err)
	defer downloader.ReleaseBuffer(buf)
	DataAssert(t, buf.Bytes())
}

// TestDownloadImage tests the DownloadImage function.
func TestDownloader_SetClient(t *testing.T) {

	downloader := media.NewLoader(nil)
	downloader.SetClient(server.Client())

	buf, err := downloader.Fetch(server.URL)
	require.NoError(t, err)
	defer downloader.ReleaseBuffer(buf)
	DataAssert(t, buf.Bytes())
}

func TestDownloader_BuffersMemoryAllocation(t *testing.T) {
	// Use the test server URL in place of the real image URL.
	downloader := media.NewLoader(nil)

	// Fetch the image multiple times to ensure the buffers are reused.
	for i := 0; i < 10; i++ {
		buf, err := downloader.Fetch(server.URL)
		require.NoError(t, err)
		defer downloader.ReleaseBuffer(buf)
		DataAssert(t, buf.Bytes())
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
	downloader := media.NewLoader(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := downloader.Fetch(ts.URL)
		require.NoError(b, err)
	}
}

func TestDownloader_Copy(t *testing.T) {

	// Use the test server URL in place of the real image URL.
	storage := NewMockStorage()
	downloader := media.NewLoader(storage)

	// Perform the download and upload.
	_, err := downloader.Upload(server.URL, "test.jpg")
	require.NoError(t, err)

	// Generate the expected upload path.
	name, err := media.Filename(server.URL)
	require.NoError(t, err)

	// Assert the data was uploaded correctly.
	uploadedData, exists := storage.data[name]
	require.True(t, exists)
	DataAssert(t, uploadedData)
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

// Save mocks the upload of data to a storage system.
func (ms *MockStorage) Save(data []byte, name string) (path string, err error) {
	if ms.data == nil {
		ms.data = make(map[string][]byte)
	}
	ms.data[name] = data
	return name, nil
}
