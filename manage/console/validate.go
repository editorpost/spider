package console

import (
	"github.com/editorpost/spider/manage/setup"
)

// Validate configuration
func Validate(s *setup.Spider, deploy setup.Deploy) error {
	_, err := s.NewCrawler(deploy)
	if err != nil {
		return err
	}
	return nil
}
