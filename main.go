package main

import (
	"encoding/json"
	"flag"
	"github.com/editorpost/spider/collect/config"
	"github.com/editorpost/spider/extract"
	"github.com/editorpost/spider/extract/fields"
	"github.com/editorpost/spider/manage/provider/windmill"
	"log/slog"
)

var (
	fCmd      = flag.String("cmd", "", "Available commands: start, trial")
	fArgs     = flag.String("args", "", "Spider arguments JSON")
	fFields   = flag.String("fields", "", "Field extractor functions JSON")
	fEntities = flag.String("entities", "", "Comma separated list of named extractors")
)

func main() {
	cmd, args, entities, ff := Flags()
	if err := Run(cmd, args, entities, ff); err != nil {
		slog.Error("run", err)
		return
	}
}

func Run(cmd string, args *config.Args, entities string, fields []*fields.Field) (err error) {

	extractors, err := extract.Extractors(fields, entities)
	if err != nil {
		return
	}

	switch cmd {
	case "start":
		return windmill.Start(args, extractors...)
	case "trial":
		return windmill.Trial(args, extractors...)
	}

	return nil
}

func Flags() (cmd string, args *config.Args, entities string, ff []*fields.Field) {

	// parse command and flags
	flag.Parse()

	cmd = FlagToString(fCmd)
	if cmd == "" {
		slog.Error("cmd flag for spider binary is not set")
		return
	}

	// argsFlag string is the JSON string of spider arguments
	argsJson := FlagToString(fArgs)
	if argsJson == "" {
		slog.Error("args flag for spider binary is not set")
		return
	}

	// entities string is the list of extractors to apply, e.g. "html,article"
	entities = FlagToString(fEntities)
	if entities == "" {
		slog.Info("extract flag is not set, use default html extractor")
	}

	// fields is the JSON string of array of field extractor functions
	fieldsJson := FlagToString(fFields)

	args = &config.Args{}
	if err := JsonToType(argsJson, args); err != nil {
		slog.Error("parse args", slog.String("args", argsJson), err)
		return
	}

	ff = make([]*fields.Field, 0)
	if err := JsonToType(fieldsJson, &ff); err != nil {
		slog.Error("parse fields", slog.String("fields", fieldsJson), err)
		return
	}

	return cmd, args, entities, ff
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
