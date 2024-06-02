package collect

import (
	"github.com/gocolly/colly/v2"
)

type Metrics interface {
	OnRequest(req *colly.Request)
	OnRetry(req *colly.Request)
	OnResponse(resp *colly.Response)
	OnError(resp *colly.Response, err error)
	OnScraped(resp *colly.Response)
	OnExtract(resp *colly.Response)
}

type MetricsFallback struct{}

func (m *MetricsFallback) OnRequest(req *colly.Request) {}

func (m *MetricsFallback) OnRetry(req *colly.Request) {}

func (m *MetricsFallback) OnResponse(resp *colly.Response) {}

func (m *MetricsFallback) OnError(resp *colly.Response, err error) {}

func (m *MetricsFallback) OnScraped(resp *colly.Response) {}

func (m *MetricsFallback) OnExtract(resp *colly.Response) {}
