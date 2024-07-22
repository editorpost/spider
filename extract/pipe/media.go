package pipe

import (
	"log/slog"
)

type (
	Media struct {
		// publicURL for public media storagePath
		publicURL string
		loader    Downloader
	}

	Downloader interface {
		Download(src, dst string) error
	}
)

// NewMedia creates a new media extractors.
// Claims extracts and replaces image urls in the document. Must be called before extractors got access to document content.
// Uploads requested media to the destination. Must be called right before saving the payload. Adds upload result to the payload.
func NewMedia(publicURL string, loader Downloader) *Media {
	return &Media{
		publicURL: publicURL,
		loader:    loader,
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
//		payload.download.Request(uri)
//		return nil
//	}
func (m *Media) Claims(payload *Payload) error {
	// set download claims implementation
	payload.claims = NewClaims(m.publicURL)
	return nil
}

// Upload creates Fn to upload requested media from claims.
func (m *Media) Upload(payload *Payload) error {

	// no media claims
	if payload.claims == nil {
		slog.Error("unexpected nil claims")
		return nil
	}

	if payload.claims.Empty() {
		return nil
	}

	// download source and upload to destination
	for _, claim := range payload.claims.All() {

		filename, err := Filename(claim.Src)
		if err != nil {
			slog.Error("failed to hash filename", slog.String("claim.Src", claim.Src), slog.String("err", err.Error()))
			continue
		}

		if err = m.loader.Download(claim.Src, filename); err != nil {
			slog.Error("failed to download media", slog.String("claim.Src", claim.Src), slog.String("err", err.Error()))
			continue
		}

		payload.claims.Done(claim.Dst)
	}

	// set uploaded media mapping to payload
	payload.Data["extract_media"] = payload.claims.Uploaded()

	return nil
}
