package store

import (
	"context"
	"github.com/editorpost/spider/extract/pipe"
	"github.com/editorpost/spider/store/ent"
	"github.com/editorpost/spider/store/ent/extractindex"
	"github.com/google/uuid"
	"github.com/lib/pq"
	_ "github.com/lib/pq"
	"strings"
)

const (
	ExtractIndexStatusPending = 1
)

type (
	ExtractIndex struct {
		db       *ent.Client
		spiderID string
	}
)

func NewExtractIndex(spiderID, dsn string) (_ *ExtractIndex, err error) {

	db, err := NewEntClient(dsn)
	if err != nil {
		return nil, err
	}

	// migrate
	if err = db.Schema.Create(context.Background()); err != nil {
		return nil, err
	}

	return &ExtractIndex{
		db:       db,
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

func (e *ExtractIndex) Save(p *pipe.Payload) error {

	spiderID, err := uuid.Parse(e.spiderID)
	if err != nil {
		return err
	}

	_, err = e.db.ExtractIndex.Create().
		SetPayloadID(p.ID).
		SetSpiderID(spiderID).
		SetTitle(p.Doc.DOM.Find("title").Text()).
		SetStatus(ExtractIndexStatusPending).
		Save(context.Background())

	return err
}

func (e *ExtractIndex) ByPayloadID(payloadID string) (*ent.ExtractIndex, error) {
	return e.db.ExtractIndex.Query().
		Where(extractindex.PayloadID(payloadID)).
		Only(context.Background())
}

func (e *ExtractIndex) Client() *ent.Client {
	return e.db
}

func (e *ExtractIndex) Close() error {
	return e.db.Close()
}
