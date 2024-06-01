package vlog

import (
	"log/slog"
)

func Mapper(r slog.Record) map[string]any {

	m := map[string]any{
		"_time": r.Time.UTC(),
		"_msg":  r.Message,
		"level": r.Level.String(),
	}

	r.Attrs(func(a slog.Attr) bool {
		m[a.Key] = a.Value.Any()
		return true
	})

	return m
}
