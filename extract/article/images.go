package article

import (
	dto "github.com/editorpost/article"
	"log/slog"
	"regexp"
	"strings"
)

// Compile regex once at package initialization
var markdownImgTag = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)

// MediaClaims is the interface for downloading images
type MediaClaims interface {
	Add(srcAbsoluteUrl string) (dst string, err error)
}

func Images(a *dto.Article, d MediaClaims) {

	matches := MarkdownSourceUrls(a.Markup)
	if matches == nil {
		return
	}

	srcDst := ImageClaims(matches, d)
	if len(srcDst) == 0 {
		return
	}

	a.Markup = MarkdownReplaceUrls(a.Markup, srcDst)

	images := dto.NewImages()
	for _, dst := range srcDst {
		image := dto.NewImage(dst)
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

func MarkdownReplaceUrls(md string, srcDst map[string]string) string {

	for src, dst := range srcDst {
		md = strings.ReplaceAll(md, src, dst)
	}

	return md
}

func ImageClaims(srcUrls []string, d MediaClaims) map[string]string {

	m := map[string]string{}

	for _, src := range srcUrls {
		dst, err := d.Add(src)
		if err != nil {
			slog.Warn("failed to add download claim", slog.String("src", src), slog.String("err", err.Error()))
			continue
		}
		m[src] = dst
	}

	return m
}
