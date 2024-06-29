package main

import (
	"encoding/json"
	"flag"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/payload"
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
	if err := Run(cmd, spider); err != nil {
		slog.Error("cmd:"+cmd, err)
		return
	}
}

func Run(cmd string, s *setup.Spider) (err error) {

	extractors, err := extract.Extractors(s.ExtractFields, s.ExtractEntities...)
	if err != nil {
		return
	}

	pipe := payload.NewPipeline(extractors...)

	switch cmd {
	case "start":
		return windmill.Start(s.Args, pipe)
	case "trial":
		return windmill.Trial(s.Args, pipe)
	}

	return nil
}

func Flags() (cmd string, spider *setup.Spider) {

	// parse command and flags
	flag.Parse()

	cmd = FlagToString(fCmd)
	if cmd == "" {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	// argsFlag string is the JSON string of spider arguments
	spiderJson := FlagToString(fSpider)
	if spiderJson == "" {
		slog.Error("args flag for spider binary is not set")
		return
	}

	spider = &setup.Spider{}
	if err := JsonToType(spiderJson, spider); err != nil {
		slog.Error("parse spider JSON", slog.String("arg", spiderJson), err)
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

func JsonToType[T any](str string, typ T) error {

	if str == "" {
		return nil
	}

	err := json.Unmarshal([]byte(str), typ)
	if err != nil {
		return err
	}

	return nil
}
