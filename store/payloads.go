package store

import (
	"context"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/store/ent"
	"github.com/editorpost/spider/store/ent/spiderpayload"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"strings"
)

const (
	ExtractIndexStatusPending = 1
)

type (
	SpiderPayloads struct {
		db       *ent.Client
		paths    PayloadPaths
		spiderID string
	}

	PayloadPaths interface {
		PayloadFile(spiderID, payloadID string) string
	}
)

func NewSpiderPayloads(spiderID, dsn string, paths PayloadPaths) (_ *SpiderPayloads, err error) {

	db, err := NewEntClient(dsn)
	if err != nil {
		return nil, err
	}

	// migrate
	if err = db.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return &SpiderPayloads{
		db:       db,
		paths:    paths,
		spiderID: spiderID,
	}, nil
}

func NewEntClient(dsn string) (c *ent.Client, err error) {

	driver := "sqlite3"

	if strings.HasPrefix(dsn, "postgres://") {
		driver = "postgres"
		dsn, err = pq.ParseURL(dsn)
		if err != nil {
			return nil, err
		}
	}

	if strings.HasPrefix(dsn, "sqlite3://") {
		driver = "sqlite3"
		dsn = strings.TrimPrefix(dsn, "sqlite3://")
	}

	_ = dsn

	return ent.Open(driver, dsn)
}

func (e *SpiderPayloads) Save(p *pipe.Payload) error {

	spiderID, err := uuid.Parse(e.spiderID)
	if err != nil {
		return err
	}

	_, err = e.db.SpiderPayload.Create().
		SetPayloadID(p.ID).
		SetSpiderID(spiderID).
		SetTitle(p.Doc.DOM.Find("title").Text()).
		SetURL(p.URL.String()).
		SetPath(e.paths.PayloadFile(e.spiderID, p.ID)).
		SetStatus(ExtractIndexStatusPending).
		Save(context.Background())

	return err
}

func (e *SpiderPayloads) ByPayloadID(payloadID string) (*ent.SpiderPayload, error) {
	return e.db.SpiderPayload.Query().
		Where(spiderpayload.PayloadID(payloadID)).
		Only(context.Background())
}

func (e *SpiderPayloads) Client() *ent.Client {
	return e.db
}

func (e *SpiderPayloads) Close() error {
	return e.db.Close()
}
