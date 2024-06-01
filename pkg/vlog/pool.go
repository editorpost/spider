package vlog

import (
	"fmt"
	"log/slog"
	"slices"
	"sync"
	"time"
)

type (
	MapperFn func(r slog.Record) map[string]any
	SenderFn func([]slog.Record) error
)

// Pool is a thread-safe pool of slog.Logger
type Pool struct {
	Timeout time.Duration // default: 10s
	Sender  SenderFn
	BufSize int // default: 1000
	pool    []slog.Record
	mute    sync.RWMutex
}

// NewPool creates a new Pool
func NewPool(sender SenderFn) *Pool {

	bufSize := 100

	if sender == nil {
		sender = StdoutSender(Mapper)
	}

	return &Pool{
		BufSize: bufSize,
		Timeout: 10 * time.Second,
		Sender:  sender,
		pool:    make([]slog.Record, 0, bufSize),
	}
}

// Ticker runs Flush every interval
func (p *Pool) Ticker(interval time.Duration) {

	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			p.Flush()
		}
	}
}

// SetBufferSize sets the limit in count of records
func (p *Pool) SetBufferSize(size int) *Pool {
	p.mute.Lock()
	defer p.mute.Unlock()

	p.BufSize = size

	return p
}

// SetSender sets the sender function
func (p *Pool) SetSender(fn SenderFn) *Pool {
	p.mute.Lock()
	defer p.mute.Unlock()

	p.Sender = fn

	return p
}

// Add adds a record to the pool
func (p *Pool) Add(logs ...slog.Record) {
	p.mute.RLock()
	defer p.mute.RUnlock()

	// skip if pool is full
	// prevent pool from growing too large
	if len(p.pool) >= p.BufSize {
		return
	}

	p.pool = append(p.pool, logs...)

	// flush the pool if it's full
	if len(p.pool) >= p.BufSize {
		p.flushUnsafe()
	}
}

// Flush flushes the pool
func (p *Pool) Flush() {

	p.mute.Lock()
	defer p.mute.Unlock()

	p.flushUnsafe()
}

// Send sends the records to the endpoint
func (p *Pool) Send(logs []slog.Record) {

	// try to send the logs
	err := p.Sender(logs)

	if err != nil {
		fmt.Printf("error sending logs: %v", err)
	}
}

// Pop pops the records from the pool
func (p *Pool) Pop() []slog.Record {

	p.mute.Lock()
	defer p.mute.Unlock()

	return p.popUnsafe()
}

// Flush flushes the pool
func (p *Pool) flushUnsafe() {

	// skip if pool is empty
	if len(p.pool) == 0 {
		return
	}

	logs := p.popUnsafe()

	go p.Send(logs)
}

// popUnsafe ret
func (p *Pool) popUnsafe() []slog.Record {

	// skip if pool is empty
	if len(p.pool) == 0 {
		return nil
	}

	// todo: avoid allocations
	buf := slices.Clone(p.pool)
	p.pool = p.pool[:0]

	return buf
}

// DumpPool the copy of the pool
func (p *Pool) DumpPool() []slog.Record {
	p.mute.RLock()
	defer p.mute.RUnlock()

	return slices.Clone(p.pool)
}

// StdoutSender sender for logs that failed to be sent to the primary endpoint
func (e ElasticIngester) StdoutSender(logs []slog.Record) error {

	for _, record := range logs {

		attrs := make([]interface{}, 0)

		record.Attrs(func(attr slog.Attr) bool {
			attrs = append(attrs, attr.Key, attr.Value)
			return true
		})

		fmt.Printf("%s %s %s %s\n", record.Time.Format(time.RFC3339), record.Level, record.Message, attrs)
	}

	return nil
}
