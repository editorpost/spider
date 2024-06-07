package main

import (
	"flag"
	"fmt"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/provider/windmill"
	"log/slog"
)

func main() {

	// parse command and flags
	flag.Parse()

	cmd := flag.String("cmd", "", "Available commands: start, trial")
	if cmd == nil {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	args := flag.String("args", "", "argsJSON")
	if args == nil {
		slog.Error("args flag for spider binary is not set")
		return
	}

	ext := flag.String("extract", "", "argsJSON")
	if ext == nil {
		slog.Info("extract flag is not set, use default html extractor")
	}

	switch *cmd {
	case "start":
		if err := windmill.Start(*args); err != nil {
			slog.Error("start", err)
		}
	case "trial":
		var err error
		var payloads []*extract.Payload
		if payloads, err = windmill.Trial(*args); err != nil {
			slog.Error("trial", err)
		}
		fmt.Println(payloads)
	}
}
