package windmill

import (
	"github.com/editorpost/spider/manage/setup"
)

// Validate configuration
func Validate(s *setup.Spider) error {

	deploy := setup.Deploy{}
	if err := LoadDeployResource(&deploy); err != nil {
		return err
	}

	_, err := s.NewCrawler(deploy)
	if err != nil {
		return err
	}
	return nil
}
