package article

import (
	dto "github.com/editorpost/article"
	"regexp"
	"strings"
)

// Compile regex once at package initialization
var markdownImgTag = regexp.MustCompile(`!\[.*?\]\((.*?)\)`)

// Downloader is the interface for downloading images
type Downloader interface {
	Download(srcAbsoluteUrl string) (dst string, err error)
}

func ArticleImages(a *dto.Article, d Downloader) error {

	md, err := MarkdownImages(a.Markup, d)
	if err != nil {
		return err
	}

	a.Markup = md
	return nil
}

// MarkdownImages replaces images in markdown with claims
func MarkdownImages(md string, d Downloader) (string, error) {

	// Find all matches
	matches := markdownImgTag.FindAllStringSubmatch(md, -1)
	if matches == nil {
		return md, nil
	}

	replacements := make(map[string]string)
	for _, match := range matches {
		src := match[1]
		if _, exists := replacements[src]; !exists {
			dst, err := d.Download(src)
			if err != nil {
				return "", err
			}
			replacements[src] = dst
		}
	}

	// Replace all occurrences of the image URLs in the markdown
	for src, dst := range replacements {
		md = strings.ReplaceAll(md, src, dst)
	}

	return md, nil
}
