package collect

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"strconv"
	"time"
)

const (
	RequestEvent    = "request"
	RetryEvent      = "retry"
	ErrorEvent      = "error"
	ScrapedEvent    = "scraped"
	ExtractionEvent = "extracted"
	ResponseEvent   = "response"

	StartTimeCtx = "metrics-request-start-time"
)

type (
	// Metrics is spider event dispatcher and VictoriaMetrics
	Metrics struct {
		jobID    string
		spiderID string
	}
)

func NewMetrics(job, spider string) *Metrics {
	return &Metrics{
		jobID:    job,
		spiderID: spider,
	}
}

func (m *Metrics) OnRequest(req *colly.Request) {

	// start time
	req.Ctx.Put(StartTimeCtx, strconv.Itoa(m.NowMilli()))

	// request count
	m.Counter(RequestEvent).Inc()

	// retry count
	if req.Ctx.Get(RetryCountCtx) != "" {
		m.OnRetry(req)
	}
}

func (m *Metrics) OnRetry(req *colly.Request) {

	m.Counter(RetryEvent).Inc()
	m.SetLatency(RetryEvent, req)
}

func (m *Metrics) OnResponse(resp *colly.Response) {
	m.Counter(ResponseEvent).Inc()
	m.SetLatency(ResponseEvent, resp.Request)
}

func (m *Metrics) OnError(resp *colly.Response, err error) {
	m.Counter(ErrorEvent).Inc()
	m.SetLatency(ErrorEvent, resp.Request)
}

func (m *Metrics) OnScraped(resp *colly.Response) {
	m.Counter(ScrapedEvent).Inc()
	m.SetLatency(ScrapedEvent, resp.Request)
}

func (m *Metrics) OnExtract(resp *colly.Response) {
	m.Counter(ExtractionEvent).Inc()
	m.SetLatency(ExtractionEvent, resp.Request)
}

func (m *Metrics) SetLatency(event string, req *colly.Request) {

	startTime := req.Ctx.Get(StartTimeCtx)
	if startTime == "" {
		return
	}

	start, err := strconv.Atoi(startTime)
	if err != nil {
		slog.Error("failed to parse start time", slog.String("err", err.Error()))
		return
	}

	// latency, seconds
	latency := float64(time.Now().UnixMilli()-int64(start)) / 1000

	// set metric name
	labels := fmt.Sprintf(`job="%s", spider="%s"`, m.jobID, m.spiderID)
	metrics.GetOrCreateHistogram(fmt.Sprintf(`spider_%s_lat{%s}`, event, labels)).Update(latency)
}

func (m *Metrics) Counter(event string) *metrics.Counter {
	format := `spider_%s_count{job="%s", spider="%s"}`
	return metrics.GetOrCreateCounter(fmt.Sprintf(format, event, m.jobID, m.spiderID))
}

func (m *Metrics) Gauge(event string) *metrics.Gauge {
	format := `spider_%s_gauge{job="%s", spider="%s"}`
	return metrics.GetOrCreateGauge(fmt.Sprintf(format, event, m.jobID, m.spiderID), nil)
}

func (m *Metrics) Now() time.Time {
	return time.Now().UTC()
}

func (m *Metrics) NowMilli() int {
	return int(m.Now().UnixMilli())
}
