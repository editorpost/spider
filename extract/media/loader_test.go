package media_test

import (
	"fmt"
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
	return make([]byte, 1024*15) // 15KB
}

func DataAssert(t *testing.T, got []byte) {
	assert.Equal(t, DataExpected(), got)
}

func TestMain(m *testing.M) {

	// run image server
	server = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(DataExpected()) // Example of JPEG header bytes.
	}))
	defer server.Close()
	m.Run()
}

// TestDownloadImage tests the DownloadImage function.
func TestDownloader_Download(t *testing.T) {
	// Set up a test server that serves an example image.
	testImage := DataExpected()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(testImage)
	}))
	defer ts.Close()

	// Use the test server Endpoint in place of the real image Endpoint.
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

func BenchmarkDownloader_Download(b *testing.B) {
	// Set up a test server that serves an example image.
	testImage := []byte{0xFF, 0xD8, 0xFF} // Example of JPEG header bytes.
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		_, _ = w.Write(testImage)
	}))
	defer ts.Close()

	// Use the test server Endpoint in place of the real image Endpoint.
	downloader := media.NewLoader(nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := downloader.Fetch(ts.URL)
		require.NoError(b, err)
	}
}

func TestDownloader_Copy(t *testing.T) {

	// Use the test server Endpoint in place of the real image Endpoint.
	storage := NewMockStorage()
	downloader := media.NewLoader(storage)

	// Perform the download and upload.
	path := "static/media/test.jpg"
	require.NoError(t, downloader.Download(server.URL, path))

	// Assert the data was uploaded correctly.
	uploadedData, exists := storage.data[path]
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
func (ms *MockStorage) Save(data []byte, name string) (err error) {
	if ms.data == nil {
		ms.data = make(map[string][]byte)
	}
	ms.data[name] = data
	return nil
}

func TestFileExtension(t *testing.T) {
	type args struct {
		uri string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr assert.ErrorAssertionFunc
	}{
		{
			name:    "Empty URI",
			args:    args{uri: ""},
			want:    "",
			wantErr: assert.NoError,
		},
		{
			name:    "Empty URI",
			args:    args{uri: "https://example.com/expected.jpg?some=unwanted.png"},
			want:    ".jpg",
			wantErr: assert.NoError,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := media.FileExtension(tt.args.uri)
			if !tt.wantErr(t, err, fmt.Sprintf("FileExtension(%v)", tt.args.uri)) {
				return
			}
			assert.Equalf(t, tt.want, got, "FileExtension(%v)", tt.args.uri)
		})
	}
}
