package extract_test

import (
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/require"
	"hash/fnv"
	"io"
	"log/slog"
	"net/http"
	"net/url"
	"path"
	"path/filepath"
	"strings"
	"sync"
	"testing"
)

type (
	// Store media data to destinations like S3.
	Store interface {
		Save(data []byte, filename string) (string, error)
	}
	// Loader copy data from publicURL to store.
	Loader struct {
		uploads sync.Map
	}
)

func NewLoader() *Loader {
	return &Loader{}
}

// Upload fetches the media from the specified URL and uploads it to the store.
func (dl *Loader) Upload(src, dst string) (string, error) {
	// Check if the media is already uploaded.
	if _, exists := dl.uploads.Load(src); exists {
		return "", nil
	}

	// Simulate the upload process.
	dl.uploads.Store(src, dst)

	return dst, nil
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

// Filename generates a unique filename hash for the media based on the source URL.
func Filename(srcURL string) (string, error) {

	if len(srcURL) == 0 {
		return "", errors.New("empty source URL for filename")
	}

	// Generate upload storagePath from the source URL using FNV hash.
	name, err := StorageHash(srcURL)
	if err != nil {
		return "", err
	}

	// add file extension from srcURL to the upload pat
	name += filepath.Ext(srcURL)

	return name, nil
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

// Claim claim to get src media and save it to dst.
type Claim struct {
	// Src is the URL of the media to download.
	Src string `json:"Src"`
	// Dst is the storagePath to save the downloaded media.
	Dst string `json:"Dst"`
	// Requested is true if the media is requested to download.
	Requested bool `json:"Requested"`
	// Done is true if the media is downloaded to destination.
	Done bool `json:"Done"`
}

type Claims struct {
	// dstURL is a prefix of public storagePath of the replaced media publicURL.
	dstURL string
	// claims keyed with source publicURL
	claims map[string]Claim
}

// NewClaims creates a new Claim for each image and replace src storagePath in document and selection.
// Replacement storagePath from media.Filename. Replaces src publicURL in selection.
func NewClaims(uri string) *Claims {

	claims := &Claims{
		dstURL: uri,
		claims: make(map[string]Claim),
	}

	return claims
}

// ExtractAndReplace Claim for each img tag and replace src storagePath in selection.
func (list *Claims) ExtractAndReplace(uri string, selection *goquery.Selection) {

	selection.Find("img").Each(func(i int, el *goquery.Selection) {

		// has src
		src, exists := el.Attr("src")
		if !exists {
			return
		}

		// already claimed
		if _, exists = list.claims[src]; exists {
			return
		}

		// already replaced
		if strings.HasPrefix(src, list.dstURL) {
			return
		}

		// filename as src publicURL hash
		dst, err := Filename(src)
		if err != nil {
			slog.Error("failed to hash filename", slog.String("src", src), slog.String("err", err.Error()))
			return
		}

		// full publicURL
		dst = path.Join(uri, dst)

		// replace publicURL in selection
		el.SetAttr("src", dst)

		// add claim
		list.Add(Claim{
			Src: src,
			Dst: dst,
		})
	})
}

func (list *Claims) Add(c Claim) *Claims {
	list.claims[c.Src] = c
	return list
}

// Request Claim for uploading by Dst publicURL.
func (list *Claims) Request(byDestinationURL string) *Claims {

	for _, claim := range list.claims {
		if claim.Dst == byDestinationURL {
			claim.Requested = true
		}
	}

	return list
}

// Done Claim marks Claim as uploaded.
func (list *Claims) Done(byDestinationURL string) *Claims {

	for _, claim := range list.claims {
		if claim.Dst == byDestinationURL {
			claim.Done = true
		}
	}

	return list
}

// Requested returns a list of requested claims.
func (list *Claims) Requested() []Claim {

	requested := make([]Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		if claim.Requested {
			requested = append(requested, claim)
		}
	}
	return requested
}

// Uploaded returns a list of uploaded claims.
func (list *Claims) Uploaded() []Claim {

	uploaded := make([]Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		if claim.Done {
			uploaded = append(uploaded, claim)
		}
	}
	return uploaded
}

func (list *Claims) All() []Claim {

	all := make([]Claim, 0, len(list.claims))
	for _, claim := range list.claims {
		all = append(all, claim)
	}
	return all
}

type (
	Extractor func(*Payload) error

	Payload struct {
		// Doc is full document
		Doc *colly.HTMLElement `json:"-"`
		// Selection of entity in document
		Selection *goquery.Selection `json:"-"`
		// URL of the document
		URL *url.URL `json:"-"`
		// Data is a map of extracted data
		Data map[string]any `json:"Data"`
		// Claims is a list of media to upload
		Claims *Claims `json:"Claims"`
	}
)
type Media struct {
	pubURL string
	loader *Loader
}

// NewMedia creates a new media extractor.
// All images urls replaced by predefined public publicURL.
// Claims requested by extract.Fn will be downloaded from the source and uploaded to the destination.
func NewMedia(pubURL string, loader *Loader) *Media {
	return &Media{
		pubURL: pubURL,
		loader: loader,
	}
}

// Claims extracts all images urls from `src` attribute in the document.
// Creates a claim for each image with the source and desired destination publicURL.
// DOM source urls are replaced with the destination publicURL.
//
// Url replacement affected to every next extract.Fn in the pipeline.
// Any extract.Fn might mark Claim as required to upload.
// For example:
//
//	func(payload *Context) error {
//		uri := getImageUrlFromDocumentToUpload()
//		payload.Claims.Request(uri)
//		return nil
//	}
func (m *Media) Claims(payload *Payload) error {

	payload.Claims = NewClaims(m.pubURL)
	payload.Claims.ExtractAndReplace(payload.URL.String(), payload.Doc.DOM)

	return nil
}

// Upload creates Fn to upload requested media from claims.
func (m *Media) Upload() (Extractor, error) {

	return func(payload *Payload) error {

		// skip if no requested claims
		requested := payload.Claims.Requested()
		if len(requested) == 0 {
			return nil
		}

		// download source and upload to destination
		for _, claim := range requested {
			if _, err := m.loader.Upload(claim.Src, claim.Dst); err != nil {
				slog.Error("failed to download media", slog.String("publicURL", claim.Src), slog.String("err", err.Error()))
				continue
			}
			payload.Claims.Done(claim.Dst)
		}

		// set uploaded media mapping to payload
		payload.Data["media"] = payload.Claims.Uploaded()

		return nil
	}, nil
}

func TestNewMedia(t *testing.T) {

	doc := GetArticleDocument(t)

	m := NewMedia("https://dst.com", NewLoader())

	p := &Payload{
		Doc:       doc,
		Selection: doc.DOM,
		URL:       &url.URL{},
		Data:      map[string]any{},
	}

	// Claims extracts all images urls from `src` attribute in the document.
	err := m.Claims(p)
	require.NoError(t, err)
	require.NotZero(t, len(p.Claims.All()))

	// check publicURL

}
