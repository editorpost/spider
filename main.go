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

func Flags() (cmd, args, entities, fields string) {

	cmd = FlagToString(flag.String("cmd", "", "Available commands: start, trial"))
	if cmd == "" {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	// argsFlag string is the JSON string of spider arguments
	args = FlagToString(flag.String("args", "", "Spider arguments JSON"))
	if args == "" {
		slog.Error("args flag for spider binary is not set")
		return
	}

	// entities string is the list of extractors to apply, e.g. "html,article"
	entities = FlagToString(flag.String("entities", "", "Comma separated list of named extractors"))
	if entities == "" {
		slog.Info("extract flag is not set, use default html extractor")
	}

	// fields is the JSON string of array of field extractor functions
	fields = FlagToString(flag.String("fields", "", "Field extractor functions JSON"))
	if fields == "" {
		slog.Info("extract flag is not set, use default html extractor")
	}

	// parse command and flags
	flag.Parse()

	return cmd, args, entities, fields
}

func FlagToString(flag *string) string {
	if flag == nil {
		return ""
	}
	return *flag
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
