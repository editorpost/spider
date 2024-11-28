package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/index"
	"github.com/google/uuid"
	"time"
)

// SpiderPayload holds the schema definition for the SpiderPayload entity.
type SpiderPayload struct {
	ent.Schema
}

// Fields of the ExtractIndex.
func (SpiderPayload) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("spider_id", uuid.UUID{}),
		field.String("payload_id").NotEmpty(),
		field.Time("extracted_at").Default(time.Now),
		field.Uint8("status").Default(1),
		field.String("title"),
	}
}

// Edges of the ExtractIndex.
func (SpiderPayload) Edges() []ent.Edge {
	return nil
}

func (SpiderPayload) Indexes() []ent.Index {
	return []ent.Index{
		// unique constraint on (identity_id, user_agent)
		index.Fields("spider_id", "extracted_at").Unique(),
	}
}
