package store_test

import (
	"context"
	"github.com/editorpost/spider/store"
	"github.com/editorpost/spider/store/ent"
	"github.com/editorpost/spider/tester"
	"github.com/google/uuid"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/require"
	"log"
	"testing"
)

func TestMain(m *testing.M) {

	client, err := ent.Open("sqlite3", "file:ent?mode=memory&cache=shared&_fk=1")
	if err != nil {
		log.Fatalf("failed opening connection to sqlite: %v", err)
	}
	defer func(client *ent.Client) {
		_ = client.Close()
	}(client)

	// Run the auto migration tool.
	if err = client.Schema.Create(context.Background()); err != nil {
		log.Fatalf("failed creating schema resources: %v", err)
	}

	m.Run()
}

func TestSpiderPayloads_Save(t *testing.T) {

	spiderID := uuid.New().String()
	payload := tester.TestPayload(t, "../tester/fixtures/news/article-1.html")
	deploy := tester.TestDeploy(t)

	idx, err := store.NewSpiderPayloads(spiderID, "sqlite3://file:ent?mode=memory&cache=shared&_fk=1", deploy.Paths)
	require.NoError(t, err)
	defer func(idx *store.SpiderPayloads) {
		_ = idx.Close()
	}(idx)

	require.NoError(t, idx.Save(payload))

	// load and check
	row, err := idx.ByPayloadID(payload.ID)
	require.NoError(t, err)
	require.Equal(t, payload.ID, row.PayloadID)
}
