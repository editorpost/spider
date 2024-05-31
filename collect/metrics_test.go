package collect_test

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/editorpost/spider/collect"
	"github.com/gocolly/colly/v2"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"
)

func TestMetrics_OnRequest(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}

	m.OnRequest(req)

	startTime := req.Ctx.Get(collect.StartTimeCtx)
	assert.NotEmpty(t, startTime, "Start time should be set in context")
	assert.NotNil(t, m.Counter(collect.RequestEvent), "Request counter should not be nil")

	// Check retry logic
	req.Ctx.Put(collect.RetryCountCtx, "1")
	m.OnRequest(req)
	assert.NotNil(t, m.Counter(collect.RetryEvent), "Retry counter should not be nil")
}

func TestMetrics_OnRetry(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}

	req.Ctx.Put(collect.StartTimeCtx, strconv.Itoa(m.NowMilli()))

	m.OnRetry(req)
	assert.NotNil(t, m.Counter(collect.RetryEvent), "Retry counter should not be nil")
}

func TestMetrics_OnResponse(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}
	resp := &colly.Response{
		Request: req,
	}

	req.Ctx.Put(collect.StartTimeCtx, strconv.Itoa(m.NowMilli()))

	m.OnResponse(resp)
	assert.NotNil(t, m.Counter(collect.ResponseEvent), "Response counter should not be nil")
}

func TestMetrics_OnError(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}
	resp := &colly.Response{
		Request: req,
	}

	req.Ctx.Put(collect.StartTimeCtx, strconv.Itoa(m.NowMilli()))

	m.OnError(resp, nil)
	assert.NotNil(t, m.Counter(collect.ErrorEvent), "Error counter should not be nil")
}

func TestMetrics_OnScraped(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}
	resp := &colly.Response{
		Request: req,
	}

	req.Ctx.Put(collect.StartTimeCtx, strconv.Itoa(m.NowMilli()))

	m.OnScraped(resp)
	assert.NotNil(t, m.Counter(collect.ScrapedEvent), "Scraped counter should not be nil")
}

func TestMetrics_OnExtract(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}
	resp := &colly.Response{
		Request: req,
	}

	req.Ctx.Put(collect.StartTimeCtx, strconv.Itoa(m.NowMilli()))

	m.OnExtract(resp)
	assert.NotNil(t, m.Counter(collect.ExtractionEvent), "Extraction counter should not be nil")
}

func TestMetrics_SetLatency(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}

	startTime := strconv.Itoa(m.NowMilli())
	req.Ctx.Put(collect.StartTimeCtx, startTime)

	m.SetLatency(collect.RequestEvent, req)
	histogram := metrics.GetOrCreateHistogram("spider_request_lat{job=\"job1\", spider=\"spider1\"}")
	assert.NotNil(t, histogram, "Histogram should not be nil")
}

func TestMetrics_SetLatency_InvalidStartTime(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	req := &colly.Request{
		Ctx: colly.NewContext(),
	}

	req.Ctx.Put(collect.StartTimeCtx, "invalid")

	m.SetLatency(collect.RequestEvent, req)
	// Check that no panic occurred and invalid start time is handled gracefully
}

func TestMetrics_Counter(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	counter := m.Counter(collect.RequestEvent)
	assert.NotNil(t, counter, "Counter should not be nil")
}

func TestMetrics_Gauge(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	gauge := m.Gauge(collect.RequestEvent)
	assert.NotNil(t, gauge, "Gauge should not be nil")
}

func TestMetrics_Now(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	now := m.Now()
	assert.WithinDuration(t, time.Now().UTC(), now, time.Second, "Now should return current UTC time")
}

func TestMetrics_NowMilli(t *testing.T) {
	m := collect.NewMetrics("job1", "spider1")
	nowMilli := m.NowMilli()
	assert.True(t, nowMilli > 0, "NowMilli should return a positive millisecond timestamp")
}

func TestExampleMetrics(t *testing.T) {

	// Создаем экземпляр Metrics
	m := collect.NewMetrics("job1", "spider1")

	// Создаем новый запрос Colly
	req := &colly.Request{
		Ctx: colly.NewContext(),
		URL: &url.URL{
			Scheme: "http",
			Host:   "example.com",
			Path:   "/",
		},
	}

	// Устанавливаем время начала и вызываем обработчики событий
	req.Ctx.Put(collect.StartTimeCtx, strconv.Itoa(m.NowMilli()))

	m.OnRequest(req)

	// Эмулируем успешный ответ
	resp := &colly.Response{
		Request:    req,
		StatusCode: http.StatusOK,
	}
	m.OnResponse(resp)

	// Эмулируем ошибку
	m.OnError(resp, fmt.Errorf("example error"))

	// Эмулируем событие Scraped
	m.OnScraped(resp)

	// Эмулируем событие Extract
	m.OnExtract(resp)

	// Задержка для гарантированной отправки данных
	time.Sleep(2 * time.Second)

	// Отправка метрик в VictoriaMetrics
	http.HandleFunc("/metrics", func(w http.ResponseWriter, req *http.Request) {
		metrics.WritePrometheus(w, true)
	})
	go func() {
		err := http.ListenAndServe(":9090", nil) // Порт изменен на 9090
		if err != nil {
			fmt.Printf("Error starting HTTP server: %v\n", err)
		}
	}()
	time.Sleep(60 * time.Second)
}
