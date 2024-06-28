package media

import (
	"bytes"
	"errors"
	"io"
	"net/http"
	"net/url"
	"sync"
)

// Downloader manages downloading data and uses a pool for bytes.Buffer.
type Downloader struct {
	pool  sync.Pool
	proxy *http.Transport
}

// NewDownloader creates a new Downloader.
func NewDownloader(proxy *http.Transport) *Downloader {
	return &Downloader{
		pool: sync.Pool{
			New: func() interface{} {
				return new(bytes.Buffer)
			},
		},
		proxy: proxy,
	}
}

// Download downloads an data from the specified URL using the provided http.Transport.
func (dl *Downloader) Download(imageURL string) (*bytes.Buffer, error) {
	// Parse the URL to ensure it's valid.
	parsedURL, err := url.Parse(imageURL)
	if err != nil {
		return nil, err
	}
	dl.pool.Get()
	// Create a new HTTP client with the specified transport.
	client := &http.Client{
		Transport: dl.proxy,
	}

	// Send the GET request.
	resp, err := client.Get(parsedURL.String())
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// Check if the response status is OK.
	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("failed to download image: " + resp.Status)
	}

	// Get a buffer from the pool.
	buf := dl.pool.Get().(*bytes.Buffer)
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
