package windmill

import (
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/manage/setup"
)

func Command(cmd string, s *setup.Spider, d setup.Deploy) (err error) {

	switch cmd {

	case "start":
		return console.Start(s, d)
	case "validate":
		return console.Validate(s, d)
	case "trial":
		return Trial(s)
	}

	return nil
}
