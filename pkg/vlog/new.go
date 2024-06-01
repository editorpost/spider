package vlog

import "log/slog"

// StdoutLogger sets the vlog as default slog logger
func StdoutLogger(attrs ...slog.Attr) {
	slog.SetDefault(New(attrs...))
}

// ElasticLogger sets the vlog as default slog logger with ElasticSearch sender
func ElasticLogger(uri string, attrs ...slog.Attr) {
	slog.SetDefault(NewElastic(uri, attrs...))
}

// New creates a new vlog with Stdout sender.
func New(attrs ...slog.Attr) *slog.Logger {

	// buffering/sending logs
	pool := NewPool(StdoutSender(Mapper))
	go pool.Ticker(1)

	// catching logs
	handler := NewBaseHandler(slog.LevelInfo, pool, attrs)
	return slog.New(handler)
}

// NewElastic creates a new vlog with ElasticSearch sender.
func NewElastic(uri string, attrs ...slog.Attr) *slog.Logger {

	// buffering/sending logs
	pool := NewPool(NewElasticIngest(uri, Mapper).Sender())
	go pool.Ticker(5)

	// catching logs
	handler := NewBaseHandler(slog.LevelInfo, pool, attrs)
	return slog.New(handler)
}
