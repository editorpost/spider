package media

import (
	"context"
	"github.com/editorpost/spider/extract/payload"
	"log/slog"
)

type ClaimsCtx string

const (
	// ClaimsCtxKey is a key for media claims in the payload context.
	ClaimsCtxKey ClaimsCtx = "extract.media.claims"
)

type (
	Media struct {
		// publicURL for public media storagePath
		publicURL string
		// storagePath for media storage
		storagePath string
		loader      Uploader
	}

	Uploader interface {
		Upload(src, dst string) (string, error)
	}
)

// NewMedia creates a new media extractors.
// Claims extracts and replaces image urls in the document. Must be called before extractors relying on document content.
// Uploads requested media to the destination. Must be called right before saving the payload. Adds upload result to the payload.
func NewMedia(publicURL, storagePath string, loader Uploader) *Media {
	return &Media{
		publicURL:   publicURL,
		storagePath: storagePath,
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

	claims := NewClaims(m.publicURL).ExtractAndReplace(payload.Doc.DOM)
	payload.Ctx = context.WithValue(payload.Ctx, ClaimsCtxKey, claims)

	return nil
}

// Upload creates Fn to upload requested media from claims.
func (m *Media) Upload() (payload.Extractor, error) {

	return func(payload *payload.Payload) error {

		claims, ok := payload.Ctx.Value(ClaimsCtxKey).(*Claims)
		if !ok {
			slog.Error("claims not found in payload context")
			return nil
		}

		// skip if no requested claims
		requested := claims.Requested()
		if len(requested) == 0 {
			return nil
		}

		// download source and upload to destination
		for _, claim := range requested {
			if _, err := m.loader.Upload(claim.Src, m.storagePath); err != nil {
				slog.Error("failed to download media", slog.String("publicURL", claim.Src), slog.String("err", err.Error()))
				continue
			}
			claims.Done(claim.Dst)
		}

		// set uploaded media mapping to payload
		payload.Data["extract_media"] = claims.Uploaded()

		return nil
	}, nil
}
