package extract

import (
	"errors"
	"github.com/PuerkitoBio/goquery"
	"strconv"
	"strings"
)

// FilterFn это функция для фильтрации изображений
type FilterFn func(src string, selection *goquery.Selection) bool

// ImagesURLs проходит по содержимому статьи и извлекает URL всех изображений, применяя фильтры.
func ImagesURLs(p *Payload, filters ...FilterFn) error {
	content, ok := p.Data["entity__content"].(string)
	if !ok {
		return errors.New("content not found or invalid")
	}

	doc, err := goquery.NewDocumentFromReader(strings.NewReader(content))
	if err != nil {
		return err
	}

	var images []string
	doc.Find("img").Each(func(i int, s *goquery.Selection) {
		src, exists := s.Attr("src")
		if exists {
			pass := true
			for _, filter := range filters {
				if !filter(src, s) {
					pass = false
					break
				}
			}
			if pass {
				images = append(images, src)
			}
		}
	})

	p.Data["entity__images"] = images
	return nil
}

// MinSizeFilter фильтрует изображения по минимальному размеру
func MinSizeFilter(minWidth, minHeight int) FilterFn {
	return func(src string, selection *goquery.Selection) bool {
		widthStr, _ := selection.Attr("width")
		heightStr, _ := selection.Attr("height")

		width, _ := strconv.Atoi(widthStr)
		height, _ := strconv.Atoi(heightStr)

		return (width == 0 || width >= minWidth) && (height == 0 || height >= minHeight)
	}
}
