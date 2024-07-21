package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
	"time"
)

//
//ExtractIndex struct {
//PayloadID   string `json:"PayloadID"`
//SpiderID    string `json:"SpiderID"`
//ExtractedAt string `json:"ExtractedAt"`
//}

// ExtractIndex holds the schema definition for the ExtractIndex entity.
type ExtractIndex struct {
	ent.Schema
}

// Fields of the ExtractIndex.
func (ExtractIndex) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.UUID{}).Default(uuid.New),
		field.UUID("payload_id", uuid.UUID{}),
		field.UUID("spider_id", uuid.UUID{}),
		field.String("title"),
		field.Time("extracted_at").Default(time.Now),
		field.Uint8("status").Default(1),
	}
}

// Edges of the ExtractIndex.
func (ExtractIndex) Edges() []ent.Edge {
	return nil
}
