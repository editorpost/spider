package setup

import (
	"log/slog"
	"sync"
)

func (s *Spider) onShutdown(fn func() error) {
	s.shutdown = append(s.shutdown, fn)
}

func (s *Spider) Shutdown(fns ...func() error) {

	// wrap all shutdown
	wg := sync.WaitGroup{}

	for _, fn := range fns {
		wg.Add(1)
		go func(fn func() error) {
			defer wg.Done()
			if err := fn(); err != nil {
				slog.Error("close error", slog.String("err", err.Error()))
			}
		}(fn)
	}

	wg.Wait()
}
