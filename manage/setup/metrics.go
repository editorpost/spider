package setup

import (
	"fmt"
	"github.com/VictoriaMetrics/metrics"
	"github.com/avast/retry-go"
	"github.com/editorpost/spider/collect"
	"github.com/gocolly/colly/v2"
	"log/slog"
	"net"
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

// Metrics is spider event dispatcher and VictoriaMetrics
type VictoriaMetrics struct {
	jobID    string
	spiderID string
}

// NewMetrics creates a new Metrics instance
func NewMetrics(job, spider, promUrl string) (collect.Metrics, error) {

	// pusher
	err := retry.Do(
		func() error {
			slog.Info("init metrics", slog.String("url", promUrl))
			return metrics.InitPush(
				promUrl, // "http://localhost:35021/api/v1/import/prometheus",
				5*time.Second,
				fmt.Sprintf(`job="%s", spider="%s"`, job, spider),
				false,
			)
		},
		retry.Attempts(3), retry.Delay(5*time.Second),
	)

	if err != nil {
		slog.Error("failed to init metrics", slog.String("err", err.Error()))
		return nil, err
	}

	return &VictoriaMetrics{
		jobID:    job,
		spiderID: spider,
	}, nil
}

func (m *VictoriaMetrics) Init() *VictoriaMetrics {

	return m
}

func (m *VictoriaMetrics) OnRequest(req *colly.Request) {

	// start time
	req.Ctx.Put(StartTimeCtx, strconv.Itoa(m.NowMilli()))

	// request count
	m.Counter(RequestEvent).Inc()

	// retry count
	if req.Ctx.Get(collect.RetryCountCtx) != "" {
		m.OnRetry(req)
	}
}

func (m *VictoriaMetrics) OnRetry(req *colly.Request) {

	m.Counter(RetryEvent).Inc()
	m.SetLatency(RetryEvent, req)
}

func (m *VictoriaMetrics) OnResponse(resp *colly.Response) {
	m.Counter(ResponseEvent).Inc()
	m.SetLatency(ResponseEvent, resp.Request)
}

func (m *VictoriaMetrics) OnError(resp *colly.Response, err error) {
	m.Counter(ErrorEvent).Inc()
	m.SetLatency(ErrorEvent, resp.Request)
}

func (m *VictoriaMetrics) OnScraped(resp *colly.Response) {
	m.Counter(ScrapedEvent).Inc()
	m.SetLatency(ScrapedEvent, resp.Request)
}

func (m *VictoriaMetrics) OnExtract(resp *colly.Response) {
	m.Counter(ExtractionEvent).Inc()
	m.SetLatency(ExtractionEvent, resp.Request)
}

func (m *VictoriaMetrics) SetLatency(event string, req *colly.Request) {

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

	// get host without port
	host, _, _ := net.SplitHostPort(req.URL.Host)

	// set metric name
	labels := fmt.Sprintf(`host="%s", path="%s"`, host, req.URL.Path)
	metrics.GetOrCreateHistogram(fmt.Sprintf(`spider_%s_lat{%s}`, event, labels)).Update(latency)
}

func (m *VictoriaMetrics) Counter(event string) *metrics.Counter {
	format := `spider_%s_count{job="%s", spider="%s"}`
	return metrics.GetOrCreateCounter(fmt.Sprintf(format, event, m.jobID, m.spiderID))
}

func (m *VictoriaMetrics) CounterUrl(event, url string) *metrics.Counter {
	format := `spider_%s_count{url="%s"}`
	return metrics.GetOrCreateCounter(fmt.Sprintf(format, event, url))
}

func (m *VictoriaMetrics) Gauge(event string) *metrics.Gauge {
	format := `spider_%s_gauge{job="%s", spider="%s"}`
	return metrics.GetOrCreateGauge(fmt.Sprintf(format, event, m.jobID, m.spiderID), nil)
}

func (m *VictoriaMetrics) GaugeUrl(event, url string) *metrics.Gauge {
	format := `spider_%s_gauge{url="%s"}`
	return metrics.GetOrCreateGauge(fmt.Sprintf(format, event, url), nil)

}

func (m *VictoriaMetrics) Now() time.Time {
	return time.Now().UTC()
}

func (m *VictoriaMetrics) NowMilli() int {
	return int(m.Now().UnixMilli())
}
