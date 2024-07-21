package main

import (
	"flag"
	"github.com/editorpost/spider/manage/provider/windmill"
	"github.com/editorpost/spider/manage/setup"
	"log/slog"
)

var (
	fCmd    = flag.String("cmd", "", "Available commands: start, trial")
	fSpider = flag.String("spider", "", "Spider arguments as JSON string")
)

func main() {
	cmd, spider := Flags()
	if err := windmill.Command(cmd, spider); err != nil {
		slog.Error("cmd:"+cmd, slog.String("error", err.Error()))
		return
	}
}

func Flags() (cmd string, spider *setup.Spider) {

	// parse command and flags
	flag.Parse()

	cmd = FlagToString(fCmd)
	if cmd == "" {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	// JSON string of setup.Spider
	spiderJson := FlagToString(fSpider)
	if spiderJson == "" {
		slog.Error("args flag for spider binary is not set")
		return
	}

	spider, err := setup.NewSpiderFromJSON([]byte(spiderJson))
	if err != nil {
		slog.Error("parse spider JSON", slog.String("arg", spiderJson), slog.String("error", err.Error()))
		return
	}

	return cmd, spider
}

func FlagToString(flag *string) string {
	if flag == nil {
		return ""
	}
	return *flag
}
