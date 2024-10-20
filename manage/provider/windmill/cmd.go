package windmill

import (
	"github.com/editorpost/spider/manage/console"
	"github.com/editorpost/spider/manage/setup"
)

func Command(cmd string, s *setup.Spider) (err error) {

	switch cmd {

	case "start":
		return console.Start(s)
	case "validate":
		return console.Validate(s)
	case "trial":
		return Check(s)
	}

	return nil
}
