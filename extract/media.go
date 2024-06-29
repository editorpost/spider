package extract

import (
	"github.com/editorpost/spider/extract/media"
	"github.com/editorpost/spider/extract/payload"
	"log/slog"
)

type Media struct {
	// publicURL for public media storagePath
	publicURL string
	// storagePath for media storage
	storagePath string
	loader      *media.Loader
}

// NewMedia creates a new media extractors.
// Claims extracts and replaces image urls in the document. Must be called before extractors relying on document content.
// Uploads requested media to the destination. Must be called right before saving the payload. Adds upload result to the payload.
func NewMedia(url, path string, loader *media.Loader) *Media {
	return &Media{
		publicURL:   url,
		storagePath: path,
		loader:      loader,
	}
}

// Claims extracts all images urls from `src` attribute in the document.
// Creates a claim for each image with the source and desired destination publicURL.
// DOM source urls are replaced with the destination publicURL.
//
// Url replacement affected to every next extract.Fn in the pipeline.
// Any extract.Fn might mark media.Claim as required to upload.
// For example:
//
//	func(payload *Context) error {
//		uri := getImageUrlFromDocumentToUpload()
//		payload.Claims.Request(uri)
//		return nil
//	}
func (m *Media) Claims(payload *payload.Payload) error {

	payload.Claims = media.NewClaims(m.publicURL).ExtractAndReplace(payload.Doc.DOM)

	return nil
}

// Upload creates Fn to upload requested media from claims.
func (m *Media) Upload() (payload.Extractor, error) {

	return func(payload *payload.Payload) error {

		// skip if no requested claims
		requested := payload.Claims.Requested()
		if len(requested) == 0 {
			return nil
		}

		// download source and upload to destination
		for _, claim := range requested {
			if _, err := m.loader.Upload(claim.Src, m.storagePath); err != nil {
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
