package windmill

import (
	"github.com/editorpost/spider/manage/setup"
)

// Start is an example code for running spider
// as Windmill Script with extract.Article
//
//goland:noinspection GoUnusedExportedFunction
func Start(s *setup.Spider) (err error) {

	deploy := &setup.Deploy{}

	if err = DeployResource(deploy); err != nil {
		return err
	}

	crawler, err := s.NewCrawler(deploy)
	if err != nil {
		return err
	}

	return crawler.Run()
}
