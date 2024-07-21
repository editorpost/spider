package windmill

import "github.com/editorpost/spider/manage/setup"

func Command(cmd string, s *setup.Spider) (err error) {

	switch cmd {

	case "start":
		return Start(s)
	case "trial":
		return Trial(s)
	}

	return nil
}
