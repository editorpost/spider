package console

import (
	"github.com/editorpost/spider/manage/setup"
)

// Validate configuration
func Validate(s *setup.Spider) error {
	_, err := s.NewCrawler()
	if err != nil {
		return err
	}
	return nil
}
