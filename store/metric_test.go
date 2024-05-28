//go:build e2e

package store_test

import (
	"github.com/editorpost/donq/mongodb"
	"github.com/editorpost/spider/store"
	"log"
	"os"
	"testing"
	"time"
)

const (
	testDbName = "test_test"
)

var (
	mongoCfg *mongodb.Config
)

func TestMain(m *testing.M) {

	res := map[string]any{
		"db": "spider_meta",
		"servers": []any{
			map[string]any{
				"host": "localhost",
				"port": 27018,
			},
		},
		"credential": map[string]any{
			"username": "root",
			"password": "nopass",
		},
	}

	var err error
	if mongoCfg, err = mongodb.ConfigFromResource(res); err != nil {
		log.Fatal(err)
	}

	os.Exit(m.Run())
}

func TestTelemetryStore(t *testing.T) {

	ts, err := store.NewMetricStore(mongoCfg)
	if err != nil {
		t.Fatalf("NewMetricStore() error = %v", err)
	}

	// Define test cases
	tests := []struct {
		name         string
		events       []store.MetricRow
		wantErr      bool
		expectedBulk int
	}{
		{
			name: "Single Event",
			events: []store.MetricRow{
				{Job: "job1", URL: "http://example.com", EventType: "visited", Timestamp: time.Now()},
			},
			wantErr:      false,
			expectedBulk: 1,
		},
		{
			name: "Multiple Events",
			events: []store.MetricRow{
				{Job: "job1", URL: "http://example.com", EventType: "visited", Timestamp: time.Now()},
				{Job: "job2", URL: "http://example.net", EventType: "error", ErrorInfo: pointerToString("Page not found"), Timestamp: time.Now()},
			},
			wantErr:      false,
			expectedBulk: 2,
		},
		{
			name: "BulkWrite Error",
			events: []store.MetricRow{
				{Job: "job1", URL: "http://example.com", EventType: "visited", Timestamp: time.Now()},
			},
			wantErr:      true,
			expectedBulk: 1,
		},
	}

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			for _, event := range tt.events {
				switch event.EventType {
				case "visited":
					ts.Visited(event.Job, event.URL)
				case "extracted":
					ts.Extracted(event.Job, event.URL)
				case "error":
					ts.Error(event.Job, event.URL, *event.ErrorInfo)
				}
			}
		})
	}

	// Flush data
	ts.Close()
}

func pointerToString(s string) *string {
	return &s
}
