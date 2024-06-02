package setup

import (
	"github.com/editorpost/donq/pkg/vars"
	"github.com/editorpost/donq/pkg/vlog"
	"log/slog"
)

// VictoriaLogs sets up slog ingester to VictoriaLogs server.
// All slog messages will be sent to VictoriaLogs server.
func VictoriaLogs(uri, lvl, spider string) {

	// set windmill attributes to the logger
	vlog.VictoriaLogger(uri, LevelParse(lvl), vars.LoggerAttr(spider)...)

	// log arguments on start
	slog.Debug("start logging", slog.Any("vars", vars.FromEnv()))
}

func LevelParse(label string) slog.Level {

	var got slog.Level
	if err := got.UnmarshalText([]byte(label)); err != nil {
		slog.Error("failed to parse log level", slog.String("error", err.Error()))
		got = slog.LevelDebug
	}

	slog.Debug("log level", slog.String("level", got.String()))

	return got
}
