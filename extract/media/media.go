package media

import (
	"context"
	"github.com/editorpost/spider/extract/pipe"
	"log/slog"
	"strings"
)

var ClaimsCtxKey ClaimsCtx = "media.claims"

type (
	ClaimsCtx string

	Media struct {
		// publicURL for public media storagePath
		publicURL string
		loader    Downloader
	}

	Downloader interface {
		Download(src, dst string) error
	}
)

func NewMedia(publicURL string, loader Downloader) *Media {
	return &Media{
		publicURL: publicURL,
		loader:    loader,
	}
}

// Claims extracts all images urls from `src` attribute in the document.
// Creates a claim for each image with the source and desired destination publicURL.
func (m *Media) Claims(payload *pipe.Payload) error {
	// set claims to payload context
	payload.Ctx = context.WithValue(payload.Ctx, ClaimsCtxKey, NewClaims(m.publicURL))
	return nil
}

// Upload creates Fn to upload requested media from claims.
func (m *Media) Upload(payload *pipe.Payload) error {

	// no media claims
	claims := ClaimsFromContext(payload.Ctx)
	if claims == nil {
		return nil
	}

	if claims.Empty() {
		return nil
	}

	// download source and upload to destination
	for _, claim := range claims.All() {

		dst := strings.TrimPrefix(claim.Dst, m.publicURL)

		if err := m.loader.Download(claim.Src, dst); err != nil {
			slog.Error("failed to download media", slog.String("claim.Src", claim.Src), slog.String("err", err.Error()))
			continue
		}

		claims.Done(claim.Dst)
	}

	// set uploaded media mapping to payload
	payload.Data["extract_media"] = claims.Uploaded()

	return nil
}

func ClaimsFromContext(ctx context.Context) *Claims {
	claims, ok := ctx.Value(ClaimsCtxKey).(*Claims)
	if !ok {
		return nil
	}
	return claims
}
