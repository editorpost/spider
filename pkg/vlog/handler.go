package vlog

import (
	"context"
	slogcommon "github.com/samber/slog-common"
	"log/slog"
)

type BaseHandler struct {
	Level  slog.Leveler
	Attrs  []slog.Attr
	Groups []string
	pool   *Pool
}

func NewBaseHandler(level slog.Leveler, pool *Pool, attrs []slog.Attr) *BaseHandler {
	return &BaseHandler{
		Level: level,
		pool:  pool,
		Attrs: attrs,
	}
}

// Enabled returns true if the record should be logged
func (h *BaseHandler) Enabled(_ context.Context, level slog.Level) bool {
	return level >= h.Level.Level()
}

// Handle logs the record and put to buffer
func (h *BaseHandler) Handle(ctx context.Context, record slog.Record) error {

	// merge global attrs
	record.AddAttrs(h.Attrs...)
	// add into pull
	h.pool.Add(record)
	return nil
}

func (h *BaseHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &BaseHandler{
		Level:  h.Level,
		Attrs:  slogcommon.AppendAttrsToGroup(h.Groups, h.Attrs, attrs...),
		Groups: h.Groups,
		pool:   h.pool,
	}
}

func (h *BaseHandler) WithGroup(name string) slog.Handler {
	return &BaseHandler{
		Level:  h.Level,
		Attrs:  h.Attrs,
		Groups: append(h.Groups, name),
		pool:   h.pool,
	}
}
