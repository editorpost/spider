package config

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

func (m *MetricsFallback) OnRequest(_ *colly.Request) {}

func (m *MetricsFallback) OnRetry(_ *colly.Request) {}

func (m *MetricsFallback) OnResponse(_ *colly.Response) {}

func (m *MetricsFallback) OnError(_ *colly.Response, _ error) {}

func (m *MetricsFallback) OnScraped(_ *colly.Response) {}

func (m *MetricsFallback) OnExtract(_ *colly.Response) {}
