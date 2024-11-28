// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"sync"
	"time"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"github.com/editorpost/spider/store/ent/predicate"
	"github.com/editorpost/spider/store/ent/spiderpayload"
	"github.com/google/uuid"
)

const (
	// Operation types.
	OpCreate    = ent.OpCreate
	OpDelete    = ent.OpDelete
	OpDeleteOne = ent.OpDeleteOne
	OpUpdate    = ent.OpUpdate
	OpUpdateOne = ent.OpUpdateOne

	// Node types.
	TypeSpiderPayload = "SpiderPayload"
)

// SpiderPayloadMutation represents an operation that mutates the SpiderPayload nodes in the graph.
type SpiderPayloadMutation struct {
	config
	op            Op
	typ           string
	id            *uuid.UUID
	spider_id     *uuid.UUID
	payload_id    *string
	extracted_at  *time.Time
	url           *string
	_path         *string
	status        *uint8
	addstatus     *int8
	title         *string
	clearedFields map[string]struct{}
	done          bool
	oldValue      func(context.Context) (*SpiderPayload, error)
	predicates    []predicate.SpiderPayload
}

var _ ent.Mutation = (*SpiderPayloadMutation)(nil)

// spiderpayloadOption allows management of the mutation configuration using functional options.
type spiderpayloadOption func(*SpiderPayloadMutation)

// newSpiderPayloadMutation creates new mutation for the SpiderPayload entity.
func newSpiderPayloadMutation(c config, op Op, opts ...spiderpayloadOption) *SpiderPayloadMutation {
	m := &SpiderPayloadMutation{
		config:        c,
		op:            op,
		typ:           TypeSpiderPayload,
		clearedFields: make(map[string]struct{}),
	}
	for _, opt := range opts {
		opt(m)
	}
	return m
}

// withSpiderPayloadID sets the ID field of the mutation.
func withSpiderPayloadID(id uuid.UUID) spiderpayloadOption {
	return func(m *SpiderPayloadMutation) {
		var (
			err   error
			once  sync.Once
			value *SpiderPayload
		)
		m.oldValue = func(ctx context.Context) (*SpiderPayload, error) {
			once.Do(func() {
				if m.done {
					err = errors.New("querying old values post mutation is not allowed")
				} else {
					value, err = m.Client().SpiderPayload.Get(ctx, id)
				}
			})
			return value, err
		}
		m.id = &id
	}
}

// withSpiderPayload sets the old SpiderPayload of the mutation.
func withSpiderPayload(node *SpiderPayload) spiderpayloadOption {
	return func(m *SpiderPayloadMutation) {
		m.oldValue = func(context.Context) (*SpiderPayload, error) {
			return node, nil
		}
		m.id = &node.ID
	}
}

// Client returns a new `ent.Client` from the mutation. If the mutation was
// executed in a transaction (ent.Tx), a transactional client is returned.
func (m SpiderPayloadMutation) Client() *Client {
	client := &Client{config: m.config}
	client.init()
	return client
}

// Tx returns an `ent.Tx` for mutations that were executed in transactions;
// it returns an error otherwise.
func (m SpiderPayloadMutation) Tx() (*Tx, error) {
	if _, ok := m.driver.(*txDriver); !ok {
		return nil, errors.New("ent: mutation is not running in a transaction")
	}
	tx := &Tx{config: m.config}
	tx.init()
	return tx, nil
}

// SetID sets the value of the id field. Note that this
// operation is only accepted on creation of SpiderPayload entities.
func (m *SpiderPayloadMutation) SetID(id uuid.UUID) {
	m.id = &id
}

// ID returns the ID value in the mutation. Note that the ID is only available
// if it was provided to the builder or after it was returned from the database.
func (m *SpiderPayloadMutation) ID() (id uuid.UUID, exists bool) {
	if m.id == nil {
		return
	}
	return *m.id, true
}

// IDs queries the database and returns the entity ids that match the mutation's predicate.
// That means, if the mutation is applied within a transaction with an isolation level such
// as sql.LevelSerializable, the returned ids match the ids of the rows that will be updated
// or updated by the mutation.
func (m *SpiderPayloadMutation) IDs(ctx context.Context) ([]uuid.UUID, error) {
	switch {
	case m.op.Is(OpUpdateOne | OpDeleteOne):
		id, exists := m.ID()
		if exists {
			return []uuid.UUID{id}, nil
		}
		fallthrough
	case m.op.Is(OpUpdate | OpDelete):
		return m.Client().SpiderPayload.Query().Where(m.predicates...).IDs(ctx)
	default:
		return nil, fmt.Errorf("IDs is not allowed on %s operations", m.op)
	}
}

// SetSpiderID sets the "spider_id" field.
func (m *SpiderPayloadMutation) SetSpiderID(u uuid.UUID) {
	m.spider_id = &u
}

// SpiderID returns the value of the "spider_id" field in the mutation.
func (m *SpiderPayloadMutation) SpiderID() (r uuid.UUID, exists bool) {
	v := m.spider_id
	if v == nil {
		return
	}
	return *v, true
}

// OldSpiderID returns the old "spider_id" field's value of the SpiderPayload entity.
// If the SpiderPayload object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SpiderPayloadMutation) OldSpiderID(ctx context.Context) (v uuid.UUID, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldSpiderID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldSpiderID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldSpiderID: %w", err)
	}
	return oldValue.SpiderID, nil
}

// ResetSpiderID resets all changes to the "spider_id" field.
func (m *SpiderPayloadMutation) ResetSpiderID() {
	m.spider_id = nil
}

// SetPayloadID sets the "payload_id" field.
func (m *SpiderPayloadMutation) SetPayloadID(s string) {
	m.payload_id = &s
}

// PayloadID returns the value of the "payload_id" field in the mutation.
func (m *SpiderPayloadMutation) PayloadID() (r string, exists bool) {
	v := m.payload_id
	if v == nil {
		return
	}
	return *v, true
}

// OldPayloadID returns the old "payload_id" field's value of the SpiderPayload entity.
// If the SpiderPayload object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SpiderPayloadMutation) OldPayloadID(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPayloadID is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPayloadID requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPayloadID: %w", err)
	}
	return oldValue.PayloadID, nil
}

// ResetPayloadID resets all changes to the "payload_id" field.
func (m *SpiderPayloadMutation) ResetPayloadID() {
	m.payload_id = nil
}

// SetExtractedAt sets the "extracted_at" field.
func (m *SpiderPayloadMutation) SetExtractedAt(t time.Time) {
	m.extracted_at = &t
}

// ExtractedAt returns the value of the "extracted_at" field in the mutation.
func (m *SpiderPayloadMutation) ExtractedAt() (r time.Time, exists bool) {
	v := m.extracted_at
	if v == nil {
		return
	}
	return *v, true
}

// OldExtractedAt returns the old "extracted_at" field's value of the SpiderPayload entity.
// If the SpiderPayload object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SpiderPayloadMutation) OldExtractedAt(ctx context.Context) (v time.Time, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldExtractedAt is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldExtractedAt requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldExtractedAt: %w", err)
	}
	return oldValue.ExtractedAt, nil
}

// ResetExtractedAt resets all changes to the "extracted_at" field.
func (m *SpiderPayloadMutation) ResetExtractedAt() {
	m.extracted_at = nil
}

// SetURL sets the "url" field.
func (m *SpiderPayloadMutation) SetURL(s string) {
	m.url = &s
}

// URL returns the value of the "url" field in the mutation.
func (m *SpiderPayloadMutation) URL() (r string, exists bool) {
	v := m.url
	if v == nil {
		return
	}
	return *v, true
}

// OldURL returns the old "url" field's value of the SpiderPayload entity.
// If the SpiderPayload object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SpiderPayloadMutation) OldURL(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldURL is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldURL requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldURL: %w", err)
	}
	return oldValue.URL, nil
}

// ClearURL clears the value of the "url" field.
func (m *SpiderPayloadMutation) ClearURL() {
	m.url = nil
	m.clearedFields[spiderpayload.FieldURL] = struct{}{}
}

// URLCleared returns if the "url" field was cleared in this mutation.
func (m *SpiderPayloadMutation) URLCleared() bool {
	_, ok := m.clearedFields[spiderpayload.FieldURL]
	return ok
}

// ResetURL resets all changes to the "url" field.
func (m *SpiderPayloadMutation) ResetURL() {
	m.url = nil
	delete(m.clearedFields, spiderpayload.FieldURL)
}

// SetPath sets the "path" field.
func (m *SpiderPayloadMutation) SetPath(s string) {
	m._path = &s
}

// Path returns the value of the "path" field in the mutation.
func (m *SpiderPayloadMutation) Path() (r string, exists bool) {
	v := m._path
	if v == nil {
		return
	}
	return *v, true
}

// OldPath returns the old "path" field's value of the SpiderPayload entity.
// If the SpiderPayload object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SpiderPayloadMutation) OldPath(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldPath is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldPath requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldPath: %w", err)
	}
	return oldValue.Path, nil
}

// ClearPath clears the value of the "path" field.
func (m *SpiderPayloadMutation) ClearPath() {
	m._path = nil
	m.clearedFields[spiderpayload.FieldPath] = struct{}{}
}

// PathCleared returns if the "path" field was cleared in this mutation.
func (m *SpiderPayloadMutation) PathCleared() bool {
	_, ok := m.clearedFields[spiderpayload.FieldPath]
	return ok
}

// ResetPath resets all changes to the "path" field.
func (m *SpiderPayloadMutation) ResetPath() {
	m._path = nil
	delete(m.clearedFields, spiderpayload.FieldPath)
}

// SetStatus sets the "status" field.
func (m *SpiderPayloadMutation) SetStatus(u uint8) {
	m.status = &u
	m.addstatus = nil
}

// Status returns the value of the "status" field in the mutation.
func (m *SpiderPayloadMutation) Status() (r uint8, exists bool) {
	v := m.status
	if v == nil {
		return
	}
	return *v, true
}

// OldStatus returns the old "status" field's value of the SpiderPayload entity.
// If the SpiderPayload object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SpiderPayloadMutation) OldStatus(ctx context.Context) (v uint8, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldStatus is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldStatus requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldStatus: %w", err)
	}
	return oldValue.Status, nil
}

// AddStatus adds u to the "status" field.
func (m *SpiderPayloadMutation) AddStatus(u int8) {
	if m.addstatus != nil {
		*m.addstatus += u
	} else {
		m.addstatus = &u
	}
}

// AddedStatus returns the value that was added to the "status" field in this mutation.
func (m *SpiderPayloadMutation) AddedStatus() (r int8, exists bool) {
	v := m.addstatus
	if v == nil {
		return
	}
	return *v, true
}

// ResetStatus resets all changes to the "status" field.
func (m *SpiderPayloadMutation) ResetStatus() {
	m.status = nil
	m.addstatus = nil
}

// SetTitle sets the "title" field.
func (m *SpiderPayloadMutation) SetTitle(s string) {
	m.title = &s
}

// Title returns the value of the "title" field in the mutation.
func (m *SpiderPayloadMutation) Title() (r string, exists bool) {
	v := m.title
	if v == nil {
		return
	}
	return *v, true
}

// OldTitle returns the old "title" field's value of the SpiderPayload entity.
// If the SpiderPayload object wasn't provided to the builder, the object is fetched from the database.
// An error is returned if the mutation operation is not UpdateOne, or the database query fails.
func (m *SpiderPayloadMutation) OldTitle(ctx context.Context) (v string, err error) {
	if !m.op.Is(OpUpdateOne) {
		return v, errors.New("OldTitle is only allowed on UpdateOne operations")
	}
	if m.id == nil || m.oldValue == nil {
		return v, errors.New("OldTitle requires an ID field in the mutation")
	}
	oldValue, err := m.oldValue(ctx)
	if err != nil {
		return v, fmt.Errorf("querying old value for OldTitle: %w", err)
	}
	return oldValue.Title, nil
}

// ResetTitle resets all changes to the "title" field.
func (m *SpiderPayloadMutation) ResetTitle() {
	m.title = nil
}

// Where appends a list predicates to the SpiderPayloadMutation builder.
func (m *SpiderPayloadMutation) Where(ps ...predicate.SpiderPayload) {
	m.predicates = append(m.predicates, ps...)
}

// WhereP appends storage-level predicates to the SpiderPayloadMutation builder. Using this method,
// users can use type-assertion to append predicates that do not depend on any generated package.
func (m *SpiderPayloadMutation) WhereP(ps ...func(*sql.Selector)) {
	p := make([]predicate.SpiderPayload, len(ps))
	for i := range ps {
		p[i] = ps[i]
	}
	m.Where(p...)
}

// Op returns the operation name.
func (m *SpiderPayloadMutation) Op() Op {
	return m.op
}

// SetOp allows setting the mutation operation.
func (m *SpiderPayloadMutation) SetOp(op Op) {
	m.op = op
}

// Type returns the node type of this mutation (SpiderPayload).
func (m *SpiderPayloadMutation) Type() string {
	return m.typ
}

// Fields returns all fields that were changed during this mutation. Note that in
// order to get all numeric fields that were incremented/decremented, call
// AddedFields().
func (m *SpiderPayloadMutation) Fields() []string {
	fields := make([]string, 0, 7)
	if m.spider_id != nil {
		fields = append(fields, spiderpayload.FieldSpiderID)
	}
	if m.payload_id != nil {
		fields = append(fields, spiderpayload.FieldPayloadID)
	}
	if m.extracted_at != nil {
		fields = append(fields, spiderpayload.FieldExtractedAt)
	}
	if m.url != nil {
		fields = append(fields, spiderpayload.FieldURL)
	}
	if m._path != nil {
		fields = append(fields, spiderpayload.FieldPath)
	}
	if m.status != nil {
		fields = append(fields, spiderpayload.FieldStatus)
	}
	if m.title != nil {
		fields = append(fields, spiderpayload.FieldTitle)
	}
	return fields
}

// Field returns the value of a field with the given name. The second boolean
// return value indicates that this field was not set, or was not defined in the
// schema.
func (m *SpiderPayloadMutation) Field(name string) (ent.Value, bool) {
	switch name {
	case spiderpayload.FieldSpiderID:
		return m.SpiderID()
	case spiderpayload.FieldPayloadID:
		return m.PayloadID()
	case spiderpayload.FieldExtractedAt:
		return m.ExtractedAt()
	case spiderpayload.FieldURL:
		return m.URL()
	case spiderpayload.FieldPath:
		return m.Path()
	case spiderpayload.FieldStatus:
		return m.Status()
	case spiderpayload.FieldTitle:
		return m.Title()
	}
	return nil, false
}

// OldField returns the old value of the field from the database. An error is
// returned if the mutation operation is not UpdateOne, or the query to the
// database failed.
func (m *SpiderPayloadMutation) OldField(ctx context.Context, name string) (ent.Value, error) {
	switch name {
	case spiderpayload.FieldSpiderID:
		return m.OldSpiderID(ctx)
	case spiderpayload.FieldPayloadID:
		return m.OldPayloadID(ctx)
	case spiderpayload.FieldExtractedAt:
		return m.OldExtractedAt(ctx)
	case spiderpayload.FieldURL:
		return m.OldURL(ctx)
	case spiderpayload.FieldPath:
		return m.OldPath(ctx)
	case spiderpayload.FieldStatus:
		return m.OldStatus(ctx)
	case spiderpayload.FieldTitle:
		return m.OldTitle(ctx)
	}
	return nil, fmt.Errorf("unknown SpiderPayload field %s", name)
}

// SetField sets the value of a field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SpiderPayloadMutation) SetField(name string, value ent.Value) error {
	switch name {
	case spiderpayload.FieldSpiderID:
		v, ok := value.(uuid.UUID)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetSpiderID(v)
		return nil
	case spiderpayload.FieldPayloadID:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPayloadID(v)
		return nil
	case spiderpayload.FieldExtractedAt:
		v, ok := value.(time.Time)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetExtractedAt(v)
		return nil
	case spiderpayload.FieldURL:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetURL(v)
		return nil
	case spiderpayload.FieldPath:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetPath(v)
		return nil
	case spiderpayload.FieldStatus:
		v, ok := value.(uint8)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetStatus(v)
		return nil
	case spiderpayload.FieldTitle:
		v, ok := value.(string)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.SetTitle(v)
		return nil
	}
	return fmt.Errorf("unknown SpiderPayload field %s", name)
}

// AddedFields returns all numeric fields that were incremented/decremented during
// this mutation.
func (m *SpiderPayloadMutation) AddedFields() []string {
	var fields []string
	if m.addstatus != nil {
		fields = append(fields, spiderpayload.FieldStatus)
	}
	return fields
}

// AddedField returns the numeric value that was incremented/decremented on a field
// with the given name. The second boolean return value indicates that this field
// was not set, or was not defined in the schema.
func (m *SpiderPayloadMutation) AddedField(name string) (ent.Value, bool) {
	switch name {
	case spiderpayload.FieldStatus:
		return m.AddedStatus()
	}
	return nil, false
}

// AddField adds the value to the field with the given name. It returns an error if
// the field is not defined in the schema, or if the type mismatched the field
// type.
func (m *SpiderPayloadMutation) AddField(name string, value ent.Value) error {
	switch name {
	case spiderpayload.FieldStatus:
		v, ok := value.(int8)
		if !ok {
			return fmt.Errorf("unexpected type %T for field %s", value, name)
		}
		m.AddStatus(v)
		return nil
	}
	return fmt.Errorf("unknown SpiderPayload numeric field %s", name)
}

// ClearedFields returns all nullable fields that were cleared during this
// mutation.
func (m *SpiderPayloadMutation) ClearedFields() []string {
	var fields []string
	if m.FieldCleared(spiderpayload.FieldURL) {
		fields = append(fields, spiderpayload.FieldURL)
	}
	if m.FieldCleared(spiderpayload.FieldPath) {
		fields = append(fields, spiderpayload.FieldPath)
	}
	return fields
}

// FieldCleared returns a boolean indicating if a field with the given name was
// cleared in this mutation.
func (m *SpiderPayloadMutation) FieldCleared(name string) bool {
	_, ok := m.clearedFields[name]
	return ok
}

// ClearField clears the value of the field with the given name. It returns an
// error if the field is not defined in the schema.
func (m *SpiderPayloadMutation) ClearField(name string) error {
	switch name {
	case spiderpayload.FieldURL:
		m.ClearURL()
		return nil
	case spiderpayload.FieldPath:
		m.ClearPath()
		return nil
	}
	return fmt.Errorf("unknown SpiderPayload nullable field %s", name)
}

// ResetField resets all changes in the mutation for the field with the given name.
// It returns an error if the field is not defined in the schema.
func (m *SpiderPayloadMutation) ResetField(name string) error {
	switch name {
	case spiderpayload.FieldSpiderID:
		m.ResetSpiderID()
		return nil
	case spiderpayload.FieldPayloadID:
		m.ResetPayloadID()
		return nil
	case spiderpayload.FieldExtractedAt:
		m.ResetExtractedAt()
		return nil
	case spiderpayload.FieldURL:
		m.ResetURL()
		return nil
	case spiderpayload.FieldPath:
		m.ResetPath()
		return nil
	case spiderpayload.FieldStatus:
		m.ResetStatus()
		return nil
	case spiderpayload.FieldTitle:
		m.ResetTitle()
		return nil
	}
	return fmt.Errorf("unknown SpiderPayload field %s", name)
}

// AddedEdges returns all edge names that were set/added in this mutation.
func (m *SpiderPayloadMutation) AddedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// AddedIDs returns all IDs (to other nodes) that were added for the given edge
// name in this mutation.
func (m *SpiderPayloadMutation) AddedIDs(name string) []ent.Value {
	return nil
}

// RemovedEdges returns all edge names that were removed in this mutation.
func (m *SpiderPayloadMutation) RemovedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// RemovedIDs returns all IDs (to other nodes) that were removed for the edge with
// the given name in this mutation.
func (m *SpiderPayloadMutation) RemovedIDs(name string) []ent.Value {
	return nil
}

// ClearedEdges returns all edge names that were cleared in this mutation.
func (m *SpiderPayloadMutation) ClearedEdges() []string {
	edges := make([]string, 0, 0)
	return edges
}

// EdgeCleared returns a boolean which indicates if the edge with the given name
// was cleared in this mutation.
func (m *SpiderPayloadMutation) EdgeCleared(name string) bool {
	return false
}

// ClearEdge clears the value of the edge with the given name. It returns an error
// if that edge is not defined in the schema.
func (m *SpiderPayloadMutation) ClearEdge(name string) error {
	return fmt.Errorf("unknown SpiderPayload unique edge %s", name)
}

// ResetEdge resets all changes to the edge with the given name in this mutation.
// It returns an error if the edge is not defined in the schema.
func (m *SpiderPayloadMutation) ResetEdge(name string) error {
	return fmt.Errorf("unknown SpiderPayload edge %s", name)
}
