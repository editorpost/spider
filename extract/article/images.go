package article

import (
	dto "github.com/editorpost/article"
	"github.com/editorpost/spider/extract/media"
	"log/slog"
	"regexp"
	"strings"
)

// Compile regex once at package initialization
var markdownImgTag = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)

// MediaClaims is the interface for downloading images
type MediaClaims interface {
	Add(payloadID string, srcAbsoluteUrl string) (media.Claim, error)
}

// Images extracts images from the article and sets the images field
func Images(payloadID string, a *dto.Article, d MediaClaims) {

	// extract image urls from markdown
	matches := MarkdownSourceUrls(a.Markup)
	if matches == nil {
		return
	}

	claims := ImageClaims(payloadID, matches, d)
	if len(claims) == 0 {
		return
	}

	a.Markup = MarkdownReplaceUrls(a.Markup, claims)

	images := dto.NewImages()
	for _, dst := range claims {
		image := dto.NewImage(dst.DstPath)
		images.Add(image)
	}

	a.Images = images
}

func MarkdownSourceUrls(md string) []string {

	matches := markdownImgTag.FindAllStringSubmatch(md, -1)
	if matches == nil {
		return nil
	}

	var urls []string
	for _, match := range matches {
		urls = append(urls, match[1])
	}

	return urls
}

func MarkdownReplaceUrls(md string, claims []media.Claim) string {

	for _, claim := range claims {
		md = strings.ReplaceAll(md, claim.Src, claim.Dst)
	}

	return md
}

func ImageClaims(payloadID string, srcUrls []string, d MediaClaims) []media.Claim {

	claims := []media.Claim{}

	for _, src := range srcUrls {
		claim, err := d.Add(payloadID, src)
		if err != nil {
			slog.Warn("failed to add download claim", slog.String("src", src), slog.String("err", err.Error()))
			continue
		}
		claims = append(claims, claim)
	}

	return claims
}
