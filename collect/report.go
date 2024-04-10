package collect

import (
	"log/slog"
	"sync"
	"time"
)

type Report struct {
	collectSuccess int
	extractSuccess int
	extractFailed  int
	mux            *sync.RWMutex
}

func NewReport() *Report {

	report := &Report{
		mux: &sync.RWMutex{},
	}

	go report.PrintEvery(60 * time.Second)

	return report
}

func (r *Report) Visited() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.collectSuccess++
}

func (r *Report) Extracted() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.extractSuccess++
}

func (r *Report) ExtractFailed() {
	r.mux.Lock()
	defer r.mux.Unlock()
	r.extractFailed++
}

func (r *Report) Print() {
	r.mux.RLock()
	defer r.mux.RUnlock()
	slog.Info("report",
		slog.Int("collectSuccess", r.collectSuccess),
		slog.Int("extractSuccess", r.extractSuccess),
		slog.Int("extractFailed", r.extractFailed),
	)
}

func (r *Report) PrintEvery(d time.Duration) {

	ticker := time.NewTicker(d)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			r.Print()
		}
	}
}
