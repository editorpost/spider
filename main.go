package main

import (
	"errors"
	"flag"
	"github.com/editorpost/spider/manage/provider/windmill"
	"github.com/editorpost/spider/manage/setup"
	_ "github.com/lib/pq"
	"log/slog"
	"os"
)

var (
	fCmd    = flag.String("cmd", "", "Available commands: start, trial")
	fSpider = flag.String("spider", "", "Spider arguments as JSON string")
	fDeploy = flag.String("deploy", "", "Deploy arguments as JSON string")
)

func main() {

	cmd, spider, err := Flags()

	if err != nil {
		slog.Error("flags", slog.String("error", err.Error()))
		os.Exit(1)
	}

	if err = windmill.Command(cmd, spider); err != nil {
		slog.Error("cmd:"+cmd, slog.String("error", err.Error()))
		return
	}
}

func Flags() (cmd string, spider *setup.Spider, err error) {

	// parse command and flags
	flag.Parse()

	cmd = FlagToString(fCmd)
	if cmd == "" {
		err = errors.New("cmd flag for spider binary is not set")
		return
	}

	spider, err = setup.SpiderFromJSON([]byte(FlagToString(fSpider)))
	if err != nil {
		err = errors.New("failed to parse spider JSON")
		return
	}

	deploy := FlagToString(fDeploy)
	if len(deploy) > 0 {
		spider.Deploy, err = setup.NewDeploy(deploy)
	}

	return
}

func FlagToString(flag *string) string {
	if flag == nil {
		return ""
	}
	return *flag
}
