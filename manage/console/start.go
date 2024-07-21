package console

import (
	"github.com/editorpost/spider/manage/setup"
)

// Start is a code for running spider
// as Windmill Script with extract.Article
func Start(s *setup.Spider, deploy *setup.Deploy) error {

	crawler, err := s.NewCrawler(deploy)
	if err != nil {
		return err
	}

	defer s.Shutdown()
	return crawler.Run()
}
