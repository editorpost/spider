package vlog_test

import (
	"sync"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/editorpost/spider/pkg/vlog"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"log/slog"
)

type FakeSender struct {
	mock.Mock
}

func (f *FakeSender) Send(logs []slog.Record) error {
	args := f.Called(logs)
	return args.Error(0)
}

func TestNewPool(t *testing.T) {

	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper))
	assert.NotNil(t, pool)
	assert.Equal(t, 100, pool.BufSize)
	assert.Equal(t, 10*time.Second, pool.Timeout)
	assert.NotNil(t, pool.Sender)
	assert.Empty(t, pool.DumpPool())
}

func TestPool_Add(t *testing.T) {

	bufSize := 5
	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper)).SetBufferSize(bufSize)

	record := slog.Record{
		Message: gofakeit.Word(),
		Time:    time.Now(),
		Level:   slog.LevelInfo,
	}

	pool.Add(record)
	assert.Equal(t, 1, len(pool.DumpPool()))
	assert.Equal(t, record, pool.DumpPool()[0])

	// Test adding records beyond buffer size
	count := (bufSize * 2) + 1
	rest := (len(pool.DumpPool()) + count) / bufSize
	for i := 0; i < count; i++ {
		pool.Add(record)
	}

	assert.Equal(t, rest, len(pool.DumpPool())) // Should not exceed buffer size
}

func TestPool_Flush(t *testing.T) {

	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper))
	fakeSender := new(FakeSender)
	pool.SetSender(fakeSender.Send)
	pool.Add(slog.Record{
		Message: gofakeit.Word(),
		Time:    time.Now(),
		Level:   slog.LevelInfo,
	})

	done := make(chan struct{})

	fakeSender.On("Send", mock.Anything).Return(nil).Once().Run(func(args mock.Arguments) {
		close(done)
	})

	pool.Flush()

	select {
	case <-done:
		assert.Empty(t, pool.DumpPool())
		fakeSender.AssertExpectations(t)
	case <-time.After(5 * time.Second):
		t.Fatal("timeout")
	}

}

func TestPool_SetBufferSize(t *testing.T) {
	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper))
	newSize := 500

	pool.SetBufferSize(newSize)
	assert.Equal(t, newSize, pool.BufSize)
}

func TestPool_Pop(t *testing.T) {
	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper))
	record := slog.Record{
		Message: gofakeit.Word(),
		Time:    time.Now(),
		Level:   slog.LevelInfo,
	}
	pool.Add(record)

	popped := pool.Pop()
	assert.Equal(t, 1, len(popped))
	assert.Equal(t, record, popped[0])
	assert.Empty(t, pool.DumpPool())
}

func TestPool_Send(t *testing.T) {
	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper))
	fakeSender := new(FakeSender)
	pool.SetSender(fakeSender.Send)

	record := slog.Record{
		Message: gofakeit.Word(),
		Time:    time.Now(),
		Level:   slog.LevelInfo,
	}
	logs := []slog.Record{record}

	fakeSender.On("Send", logs).Return(nil).Once()
	pool.Send(logs)
	fakeSender.AssertExpectations(t)
}

func TestPool_DumpPool(t *testing.T) {
	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper))
	record := slog.Record{
		Message: gofakeit.Word(),
		Time:    time.Now(),
		Level:   slog.LevelInfo,
	}
	pool.Add(record)

	dump := pool.DumpPool()
	assert.Equal(t, 1, len(dump))
	assert.Equal(t, record, dump[0])
}

func TestPool_WithTicker(t *testing.T) {
	pool := vlog.NewPool(vlog.StdoutSender(vlog.Mapper))
	fakeSender := new(FakeSender)
	pool.SetSender(fakeSender.Send)

	var wg sync.WaitGroup
	wg.Add(1)

	go func() {
		defer wg.Done()
		pool.Ticker(1 * time.Millisecond)
	}()

	record := slog.Record{
		Message: gofakeit.Word(),
		Time:    time.Now(),
		Level:   slog.LevelInfo,
	}
	pool.Add(record)

	fakeSender.On("Send", mock.Anything).Return(nil).Once()

	time.Sleep(100 * time.Millisecond)

	assert.Empty(t, pool.DumpPool())
	fakeSender.AssertExpectations(t)
}
