package main

import (
	"encoding/json"
	"flag"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/provider/windmill"
	"log/slog"
	"os"
)

func main() {

	// parse command and flags
	flag.Parse()

	cmd := flag.String("cmd", "", "Available commands: start, trial")
	if cmd == nil {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	args := flag.String("args", "", "Spider arguments JSON")
	if args == nil {
		slog.Error("args flag for spider binary is not set")
		return
	}

	ext := flag.String("extract", "", "Extractor function name")
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
		// write extracted data to `./result.json` as windmill expects
		if err = vars.WriteScriptResult(payloads, "./result.json"); err != nil {
			slog.Error("write payloads", err)
		}
	}
}

func MarshalToFile(payloads []*extract.Payload, path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	enc := json.NewEncoder(file)
	enc.SetIndent("", "  ")

	return enc.Encode(payloads)
}
