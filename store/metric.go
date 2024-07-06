package store

import (
	"context"
	"fmt"
	"github.com/editorpost/donq/mongodb"
	"log/slog"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	TelemetryDatabase   = "spider"
	TelemetryCollection = "telemetry"
)

type MetricRow struct {
	Job       string
	URL       string
	EventType string
	ErrorInfo *string
	Timestamp time.Time
}

type MetricStore struct {
	collection  *mongo.Collection
	accumulator map[string]MetricRow
	mu          sync.Mutex
	flushTicker *time.Ticker
	closeChan   chan struct{}
}

// NewMetricStore expects mongodb.Config for MongoDB connection:
//
//		mongodb.Config struct {
//			Db   string `json:"db"`
//			Host string `json:"host"`
//			Port int    `json:"port"`
//			DSN  string `json:"dsn"`
//			User string `json:"password"`
//			Pass string `json:"username"`
//			TLS  bool   `json:"tls"`
//	}
func NewMetricStore(cfg *mongodb.Config) (*MetricStore, error) {

	uri := options.Client().ApplyURI(cfg.DSN)
	client, err := mongo.Connect(context.Background(), uri)
	if err != nil {
		return nil, err
	}

	ts := &MetricStore{
		collection:  client.Database(TelemetryDatabase).Collection(TelemetryCollection),
		accumulator: make(map[string]MetricRow),
		flushTicker: time.NewTicker(1 * time.Minute),
		closeChan:   make(chan struct{}),
	}

	// Run the background flushing process
	go ts.backgroundFlush()

	return ts, nil
}

// Close flushes any remaining data and stops the background flush process
func (ts *MetricStore) Close() {

	ts.flushTicker.Stop()

	if err := ts.flush(); err != nil {
		slog.Error("Error flushing data on close", slog.String("err", err.Error()))
	}

	close(ts.closeChan)
}

func (ts *MetricStore) Visited(job, url string) {
	ts.cache(job, url, "visited", nil)
}

func (ts *MetricStore) Extracted(job, url string) {
	ts.cache(job, url, "extracted", nil)
}

func (ts *MetricStore) Error(job, url string, errorMsg string) {
	ts.cache(job, url, "error", &errorMsg)
}

func (ts *MetricStore) backgroundFlush() {
	for {
		select {
		case <-ts.flushTicker.C:
			if err := ts.flush(); err != nil {
				fmt.Println("Error flushing data:", err)
			}
		case <-ts.closeChan:
			// Perform one last flush before closing
			if err := ts.flush(); err != nil {
				fmt.Println("Error flushing data on close:", err)
			}
			return
		}
	}
}

func (ts *MetricStore) cache(job, url, eventType string, errorInfo *string) {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	key := fmt.Sprintf("%s_%s", job, url) // Unique key for job and Endpoint
	event := MetricRow{
		Job:       job,
		URL:       url,
		EventType: eventType,
		ErrorInfo: errorInfo,
		Timestamp: time.Now(),
	}
	ts.accumulator[key] = event
}

//goland:noinspection GoLinter
func (ts *MetricStore) flush() error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	models := make([]mongo.WriteModel, 0, len(ts.accumulator))
	for _, event := range ts.accumulator {
		filter := bson.M{"job": event.Job, "url": event.URL}
		update := bson.M{
			"$set": bson.M{
				event.EventType: event.Timestamp,
			},
		}
		if event.ErrorInfo != nil {

			update["$set"].(bson.M)["error"] = *event.ErrorInfo
		}
		models = append(models, mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true))
	}

	if len(models) > 0 {
		_, err := ts.collection.BulkWrite(context.Background(), models)
		if err != nil {
			return err
		}
		// Clear the accumulator after flushing
		ts.accumulator = make(map[string]MetricRow)
	}
	return nil
}
