package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/manage/provider/windmill"
	"log/slog"
)

func main() {

	cmdFlag := flag.String("cmd", "", "Available commands: start, trial")
	if cmdFlag == nil {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	argsFlag := flag.String("args", "", "Spider arguments JSON")
	if argsFlag == nil {
		slog.Error("args flag for spider binary is not set")
		return
	}

	extractFlag := flag.String("extract", "", "Extractor function name")
	if extractFlag == nil {
		slog.Info("extract flag is not set, use default html extractor")
	}

	// parse command and flags
	flag.Parse()

	// parse json from args string to map[string]interface{}
	args, err := JSONStringToMap(*argsFlag)
	if err != nil {
		slog.Error("parse args", err)
		return
	}

	switch *cmdFlag {
	case "start":
		if err = windmill.Start(args); err != nil {
			slog.Error("start", err)
			return
		}
	case "trial":
		payloads := make([]*extract.Payload, 0)
		if payloads, err = windmill.Trial(args); err != nil {
			slog.Error("trial", err)
			return
		}
		// write extracted data to `./result.json` as windmill expects
		if err = vars.WriteScriptResult(payloads, "./result.json"); err != nil {
			slog.Error("write payloads", err)
			return
		}
	}
}

func JSONStringToMap(jsonStr string) (map[string]interface{}, error) {
	var argsMap map[string]interface{}
	fmt.Println("jsonStr: ", jsonStr)
	if err := json.Unmarshal([]byte(jsonStr), &argsMap); err != nil {
		return nil, err
	}
	fmt.Println("argsMap: ", argsMap)
	return argsMap, nil
}
