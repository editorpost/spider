package vlog_test

import (
	"github.com/editorpost/spider/pkg/vlog"
	"log/slog"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	vlog.StdoutLogger(slog.String("cid", "test"), slog.Time("upd", time.Now()))
	slog.Info("test message", slog.String("key", "value"))
	time.Sleep(50 * time.Millisecond)
}

func TestNewElastic(t *testing.T) {
	vlog.ElasticLogger("", slog.String("cid", "test"))
	slog.Info("test message", slog.String("key", "value"))
	time.Sleep(50 * time.Millisecond)
}
