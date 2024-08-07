// Code generated by ent, DO NOT EDIT.

package ent

import (
	"time"

	"github.com/editorpost/spider/store/ent/extractindex"
	"github.com/editorpost/spider/store/ent/schema"
	"github.com/google/uuid"
)

// The init function reads all schema descriptors with runtime code
// (default values, validators, hooks and policies) and stitches it
// to their package variables.
func init() {
	extractindexFields := schema.ExtractIndex{}.Fields()
	_ = extractindexFields
	// extractindexDescPayloadID is the schema descriptor for payload_id field.
	extractindexDescPayloadID := extractindexFields[2].Descriptor()
	// extractindex.PayloadIDValidator is a validator for the "payload_id" field. It is called by the builders before save.
	extractindex.PayloadIDValidator = extractindexDescPayloadID.Validators[0].(func(string) error)
	// extractindexDescExtractedAt is the schema descriptor for extracted_at field.
	extractindexDescExtractedAt := extractindexFields[3].Descriptor()
	// extractindex.DefaultExtractedAt holds the default value on creation for the extracted_at field.
	extractindex.DefaultExtractedAt = extractindexDescExtractedAt.Default.(func() time.Time)
	// extractindexDescStatus is the schema descriptor for status field.
	extractindexDescStatus := extractindexFields[4].Descriptor()
	// extractindex.DefaultStatus holds the default value on creation for the status field.
	extractindex.DefaultStatus = extractindexDescStatus.Default.(uint8)
	// extractindexDescID is the schema descriptor for id field.
	extractindexDescID := extractindexFields[0].Descriptor()
	// extractindex.DefaultID holds the default value on creation for the id field.
	extractindex.DefaultID = extractindexDescID.Default.(func() uuid.UUID)
}
