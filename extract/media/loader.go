package media

import (
	"bytes"
	"errors"
	"github.com/editorpost/spider/extract/pipe"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
)

type (
	// Store media data to destinations like S3.
	Store interface {
		Save(data []byte, filename string) error
	}
	// Loader copy data from url to store.
	Loader struct {
		pool         sync.Pool
		store        Store
		client       *http.Client
		skipLessThan int
	}
)

func NewLoader(store Store) *Loader {
	return &Loader{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		store:        store,
		client:       &http.Client{},
		skipLessThan: 1024 * 15, // 15KB
	}
}

// SetClient sets the HTTP client used to download the media.
// Proxy pool might be used to download media from different sources.
func (dl *Loader) SetClient(client *http.Client) {
	dl.client = client
}

// Download fetches the media from the specified Endpoint and uploads it to the store.
// Return http.ErrShortBody if the media is less than defined size in bytes.
func (dl *Loader) Download(src, dst string) error {

	// download
	buf, err := dl.Fetch(src)
	if err != nil {
		return err
	}
	defer dl.ReleaseBuffer(buf)

	// Skip small images, less than 1KB.
	if buf.Len() < dl.skipLessThan {
		return http.ErrShortBody
	}

	// upload the data
	return dl.store.Save(buf.Bytes(), dst)
}

// Fetch data from the specified Endpoint and return a buffer with the data.
func (dl *Loader) Fetch(imageURL string) (*bytes.Buffer, error) {
	// Parse the Endpoint to ensure it's valid.
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return nil, err
	}

	// Send the GET request.
	resp, err := dl.client.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Check if the response status is OK.
	if resp.StatusCode != http.StatusOK {
		return nil, http.ErrMissingFile
	}

	// Get a buffer from the pool.
	buf, _ := dl.pool.Get().(*bytes.Buffer)
	buf.Reset()

	// Read the image data into the buffer.
	_, err = io.Copy(buf, resp.Body)
	if err != nil {
		dl.pool.Put(buf)
		return nil, err
	}

	return buf, nil
}

// ReleaseBuffer returns the buffer back to the pool.
func (dl *Loader) ReleaseBuffer(buf *bytes.Buffer) {
	dl.pool.Put(buf)
}

// Filename generates a unique filename hash for the media based on the source Endpoint.
func Filename(srcURL string) (string, error) {

	if len(srcURL) == 0 {
		return "", errors.New("empty source Endpoint for filename")
	}

	// Generate upload path from the source Endpoint using FNV hash.
	name, err := pipe.Hash(srcURL)
	if err != nil {
		return "", err
	}

	ext, err := FileExtension(srcURL)
	if err != nil {
		return "", err
	}

	return name + ext, nil
}

func FileExtension(uri string) (string, error) {
	u, err := url.Parse(uri)
	if err != nil {
		return "", err
	}
	return filepath.Ext(u.Path), nil
}
