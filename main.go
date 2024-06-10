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

	fCmd, fArgs, fNamedExtract, fBuildExtract := Flags()
	if fCmd == "" || fArgs == "" {
		slog.Error("cmd or args flag for spider is not set")
		return
	}

	if err := Run(fCmd, fArgs, fNamedExtract, fBuildExtract); err != nil {
		slog.Error("run", err)
		return
	}
}

func Flags() (cmd, args, named, build string) {

	cmdFlag := flag.String("cmd", "", "Available commands: start, trial")
	if cmdFlag == nil {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	// argsFlag string is the JSON string of spider arguments
	argsFlag := flag.String("args", "", "Spider arguments JSON")
	if argsFlag == nil {
		slog.Error("args flag for spider binary is not set")
		return
	}

	// fNamedExtract string is the list of extractors to apply, e.g. "html,article"
	fNamedExtract := flag.String("named-extract", "", "Extractor function name")
	if fNamedExtract == nil {
		slog.Info("extract flag is not set, use default html extractor")
	}

	// fBuildExtract is the JSON string of array of field extractor functions
	fBuildExtract := flag.String("build-extract", "", "Field extractor functions JSON")
	if fBuildExtract == nil {
		slog.Info("extract flag is not set, use default html extractor")
	}

	// parse command and flags
	flag.Parse()

	return *cmdFlag, *argsFlag, *fNamedExtract, *fBuildExtract
}

func Run(fCmd, fArgs, fNamedExtractors, fBuildExtractors string) (err error) {

	// parse json from args string to map[string]interface{}
	args, err := JSONStringToMap(fArgs)
	if err != nil {
		slog.Error("parse args", err)
		return
	}

	// parse extractors
	extractors := extract.ExtractorsByName(fNamedExtractors)
	// build extractors
	extractors = append(extractors, extract.ExtractorsByJsonString(fBuildExtractors)...)

	switch fCmd {
	case "start":
		if err = windmill.Start(args, extractors...); err != nil {
			slog.Error("start", err)
			return
		}
	case "trial":
		payloads := make([]*extract.Payload, 0)
		if payloads, err = windmill.Trial(args, extractors...); err != nil {
			slog.Error("trial", err)
			return
		}
		// write extracted data to `./result.json` as windmill expects
		if err = vars.WriteScriptResult(payloads, "./result.json"); err != nil {
			slog.Error("write payloads", err)
			return
		}
	}

	return nil
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
