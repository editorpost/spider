package media

import (
	"bytes"
	"errors"
	"fmt"
	"hash/fnv"
	"io"
	"net/http"
	"net/url"
	"path/filepath"
	"sync"
)

type Store interface {
	Save(data []byte, filename string) (string, error)
}

// Downloader manages downloading and coping data to storage data and uses a pool for bytes.Buffer.
type Downloader struct {
	pool   sync.Pool
	store  Store
	client *http.Client
}

// NewDownloader creates a new Downloader.
func NewDownloader(store Store) *Downloader {
	return &Downloader{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		store:  store,
		client: &http.Client{},
	}
}

func (dl *Downloader) SetClient(client *http.Client) {
	dl.client = client
}

func (dl *Downloader) Filename(srcURL string) (string, error) {

	// Generate upload path from the source URL using FNV hash.
	name, err := StorageHash(srcURL)
	if err != nil {
		return "", err
	}

	// add file extension from srcURL to the upload pat
	name += filepath.Ext(srcURL)

	return name, nil
}

func (dl *Downloader) Upload(srcURL string) (string, error) {

	// download
	buf, err := dl.Download(srcURL)
	if err != nil {
		return "", err
	}
	defer dl.ReleaseBuffer(buf)

	// path
	name, err := dl.Filename(srcURL)
	if err != nil {
		return "", err
	}

	// upload the data
	path, err := dl.store.Save(buf.Bytes(), name)
	if err != nil {
		return "", err
	}

	return path, nil
}

// Download data from the specified URL and return a buffer with the data.
func (dl *Downloader) Download(imageURL string) (*bytes.Buffer, error) {
	// Parse the URL to ensure it's valid.
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
		return nil, errors.New("failed to download image: " + resp.Status)
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
func (dl *Downloader) ReleaseBuffer(buf *bytes.Buffer) {
	dl.pool.Put(buf)
}

// StorageHash generates an FNV hash from the source URL.
func StorageHash(sourceURL string) (string, error) {
	h := fnv.New64a()
	_, err := h.Write([]byte(sourceURL))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum64()), nil
}

// Download downloads an image from the specified URL using the provided http.Transport.
func Download(absoluteURL string, transport *http.Transport) ([]byte, error) {

	// Parse the URL to ensure it's valid.
	parsedURL, err := url.Parse(absoluteURL)
	if err != nil {
		return nil, err
	}

	// Create a new HTTP client with the specified transport.
	client := &http.Client{
		Transport: transport,
	}

	// Send the GET request.
	resp, err := client.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}
	//goland:noinspection GoUnhandledErrorResult
	defer resp.Body.Close()

	// Check if the response status is OK.
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download media: " + resp.Status)
	}

	// Read the media data.
	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return data, nil
}
