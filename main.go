package main

import (
	"flag"
	"github.com/editorpost/spider/manage/provider/windmill"
	"github.com/editorpost/spider/manage/setup"
	_ "github.com/lib/pq"
	"log/slog"
)

var (
	fCmd    = flag.String("cmd", "", "Available commands: start, trial")
	fSpider = flag.String("spider", "", "Spider arguments as JSON string")
	fDeploy = flag.String("deploy", "", "Deploy arguments as JSON string")
)

func main() {
	cmd, spider, deploy := Flags()
	if err := windmill.Command(cmd, spider, deploy); err != nil {
		slog.Error("cmd:"+cmd, slog.String("error", err.Error()))
		return
	}
}

func Flags() (cmd string, spider *setup.Spider, deploy setup.Deploy) {

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

	// JSON string of setup.Spider
	deployJson := FlagToString(fDeploy)
	if deployJson != "" {
		deploy, err = setup.NewDeploy(deployJson)
	}

	return
}

func FlagToString(flag *string) string {
	if flag == nil {
		return ""
	}
	return *flag
}
