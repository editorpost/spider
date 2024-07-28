package windmill

import (
	"github.com/editorpost/spider/manage/setup"
)

func Command(cmd string, s *setup.Spider) (err error) {

	switch cmd {

	case "start":
		return Start(s)
	case "trial":
		return Trial(s)
	case "validate":
		return Validate(s)
	}

	return nil
}
