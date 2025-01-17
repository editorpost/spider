// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"errors"
	"fmt"
	"time"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/editorpost/spider/store/ent/predicate"
	"github.com/editorpost/spider/store/ent/spiderpayload"
	"github.com/google/uuid"
)

// SpiderPayloadUpdate is the builder for updating SpiderPayload entities.
type SpiderPayloadUpdate struct {
	config
	hooks    []Hook
	mutation *SpiderPayloadMutation
}

// Where appends a list predicates to the SpiderPayloadUpdate builder.
func (spu *SpiderPayloadUpdate) Where(ps ...predicate.SpiderPayload) *SpiderPayloadUpdate {
	spu.mutation.Where(ps...)
	return spu
}

// SetSpiderID sets the "spider_id" field.
func (spu *SpiderPayloadUpdate) SetSpiderID(u uuid.UUID) *SpiderPayloadUpdate {
	spu.mutation.SetSpiderID(u)
	return spu
}

// SetNillableSpiderID sets the "spider_id" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillableSpiderID(u *uuid.UUID) *SpiderPayloadUpdate {
	if u != nil {
		spu.SetSpiderID(*u)
	}
	return spu
}

// SetExtractedAt sets the "extracted_at" field.
func (spu *SpiderPayloadUpdate) SetExtractedAt(t time.Time) *SpiderPayloadUpdate {
	spu.mutation.SetExtractedAt(t)
	return spu
}

// SetNillableExtractedAt sets the "extracted_at" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillableExtractedAt(t *time.Time) *SpiderPayloadUpdate {
	if t != nil {
		spu.SetExtractedAt(*t)
	}
	return spu
}

// SetURL sets the "url" field.
func (spu *SpiderPayloadUpdate) SetURL(s string) *SpiderPayloadUpdate {
	spu.mutation.SetURL(s)
	return spu
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillableURL(s *string) *SpiderPayloadUpdate {
	if s != nil {
		spu.SetURL(*s)
	}
	return spu
}

// ClearURL clears the value of the "url" field.
func (spu *SpiderPayloadUpdate) ClearURL() *SpiderPayloadUpdate {
	spu.mutation.ClearURL()
	return spu
}

// SetPath sets the "path" field.
func (spu *SpiderPayloadUpdate) SetPath(s string) *SpiderPayloadUpdate {
	spu.mutation.SetPath(s)
	return spu
}

// SetNillablePath sets the "path" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillablePath(s *string) *SpiderPayloadUpdate {
	if s != nil {
		spu.SetPath(*s)
	}
	return spu
}

// ClearPath clears the value of the "path" field.
func (spu *SpiderPayloadUpdate) ClearPath() *SpiderPayloadUpdate {
	spu.mutation.ClearPath()
	return spu
}

// SetStatus sets the "status" field.
func (spu *SpiderPayloadUpdate) SetStatus(u uint8) *SpiderPayloadUpdate {
	spu.mutation.ResetStatus()
	spu.mutation.SetStatus(u)
	return spu
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillableStatus(u *uint8) *SpiderPayloadUpdate {
	if u != nil {
		spu.SetStatus(*u)
	}
	return spu
}

// AddStatus adds u to the "status" field.
func (spu *SpiderPayloadUpdate) AddStatus(u int8) *SpiderPayloadUpdate {
	spu.mutation.AddStatus(u)
	return spu
}

// SetTitle sets the "title" field.
func (spu *SpiderPayloadUpdate) SetTitle(s string) *SpiderPayloadUpdate {
	spu.mutation.SetTitle(s)
	return spu
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillableTitle(s *string) *SpiderPayloadUpdate {
	if s != nil {
		spu.SetTitle(*s)
	}
	return spu
}

// SetJobProvider sets the "job_provider" field.
func (spu *SpiderPayloadUpdate) SetJobProvider(s string) *SpiderPayloadUpdate {
	spu.mutation.SetJobProvider(s)
	return spu
}

// SetNillableJobProvider sets the "job_provider" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillableJobProvider(s *string) *SpiderPayloadUpdate {
	if s != nil {
		spu.SetJobProvider(*s)
	}
	return spu
}

// ClearJobProvider clears the value of the "job_provider" field.
func (spu *SpiderPayloadUpdate) ClearJobProvider() *SpiderPayloadUpdate {
	spu.mutation.ClearJobProvider()
	return spu
}

// SetJobID sets the "job_id" field.
func (spu *SpiderPayloadUpdate) SetJobID(u uuid.UUID) *SpiderPayloadUpdate {
	spu.mutation.SetJobID(u)
	return spu
}

// SetNillableJobID sets the "job_id" field if the given value is not nil.
func (spu *SpiderPayloadUpdate) SetNillableJobID(u *uuid.UUID) *SpiderPayloadUpdate {
	if u != nil {
		spu.SetJobID(*u)
	}
	return spu
}

// ClearJobID clears the value of the "job_id" field.
func (spu *SpiderPayloadUpdate) ClearJobID() *SpiderPayloadUpdate {
	spu.mutation.ClearJobID()
	return spu
}

// Mutation returns the SpiderPayloadMutation object of the builder.
func (spu *SpiderPayloadUpdate) Mutation() *SpiderPayloadMutation {
	return spu.mutation
}

// Save executes the query and returns the number of nodes affected by the update operation.
func (spu *SpiderPayloadUpdate) Save(ctx context.Context) (int, error) {
	return withHooks(ctx, spu.sqlSave, spu.mutation, spu.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (spu *SpiderPayloadUpdate) SaveX(ctx context.Context) int {
	affected, err := spu.Save(ctx)
	if err != nil {
		panic(err)
	}
	return affected
}

// Exec executes the query.
func (spu *SpiderPayloadUpdate) Exec(ctx context.Context) error {
	_, err := spu.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (spu *SpiderPayloadUpdate) ExecX(ctx context.Context) {
	if err := spu.Exec(ctx); err != nil {
		panic(err)
	}
}

func (spu *SpiderPayloadUpdate) sqlSave(ctx context.Context) (n int, err error) {
	_spec := sqlgraph.NewUpdateSpec(spiderpayload.Table, spiderpayload.Columns, sqlgraph.NewFieldSpec(spiderpayload.FieldID, field.TypeUUID))
	if ps := spu.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := spu.mutation.SpiderID(); ok {
		_spec.SetField(spiderpayload.FieldSpiderID, field.TypeUUID, value)
	}
	if value, ok := spu.mutation.ExtractedAt(); ok {
		_spec.SetField(spiderpayload.FieldExtractedAt, field.TypeTime, value)
	}
	if value, ok := spu.mutation.URL(); ok {
		_spec.SetField(spiderpayload.FieldURL, field.TypeString, value)
	}
	if spu.mutation.URLCleared() {
		_spec.ClearField(spiderpayload.FieldURL, field.TypeString)
	}
	if value, ok := spu.mutation.Path(); ok {
		_spec.SetField(spiderpayload.FieldPath, field.TypeString, value)
	}
	if spu.mutation.PathCleared() {
		_spec.ClearField(spiderpayload.FieldPath, field.TypeString)
	}
	if value, ok := spu.mutation.Status(); ok {
		_spec.SetField(spiderpayload.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := spu.mutation.AddedStatus(); ok {
		_spec.AddField(spiderpayload.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := spu.mutation.Title(); ok {
		_spec.SetField(spiderpayload.FieldTitle, field.TypeString, value)
	}
	if value, ok := spu.mutation.JobProvider(); ok {
		_spec.SetField(spiderpayload.FieldJobProvider, field.TypeString, value)
	}
	if spu.mutation.JobProviderCleared() {
		_spec.ClearField(spiderpayload.FieldJobProvider, field.TypeString)
	}
	if value, ok := spu.mutation.JobID(); ok {
		_spec.SetField(spiderpayload.FieldJobID, field.TypeUUID, value)
	}
	if spu.mutation.JobIDCleared() {
		_spec.ClearField(spiderpayload.FieldJobID, field.TypeUUID)
	}
	if n, err = sqlgraph.UpdateNodes(ctx, spu.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{spiderpayload.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return 0, err
	}
	spu.mutation.done = true
	return n, nil
}

// SpiderPayloadUpdateOne is the builder for updating a single SpiderPayload entity.
type SpiderPayloadUpdateOne struct {
	config
	fields   []string
	hooks    []Hook
	mutation *SpiderPayloadMutation
}

// SetSpiderID sets the "spider_id" field.
func (spuo *SpiderPayloadUpdateOne) SetSpiderID(u uuid.UUID) *SpiderPayloadUpdateOne {
	spuo.mutation.SetSpiderID(u)
	return spuo
}

// SetNillableSpiderID sets the "spider_id" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillableSpiderID(u *uuid.UUID) *SpiderPayloadUpdateOne {
	if u != nil {
		spuo.SetSpiderID(*u)
	}
	return spuo
}

// SetExtractedAt sets the "extracted_at" field.
func (spuo *SpiderPayloadUpdateOne) SetExtractedAt(t time.Time) *SpiderPayloadUpdateOne {
	spuo.mutation.SetExtractedAt(t)
	return spuo
}

// SetNillableExtractedAt sets the "extracted_at" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillableExtractedAt(t *time.Time) *SpiderPayloadUpdateOne {
	if t != nil {
		spuo.SetExtractedAt(*t)
	}
	return spuo
}

// SetURL sets the "url" field.
func (spuo *SpiderPayloadUpdateOne) SetURL(s string) *SpiderPayloadUpdateOne {
	spuo.mutation.SetURL(s)
	return spuo
}

// SetNillableURL sets the "url" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillableURL(s *string) *SpiderPayloadUpdateOne {
	if s != nil {
		spuo.SetURL(*s)
	}
	return spuo
}

// ClearURL clears the value of the "url" field.
func (spuo *SpiderPayloadUpdateOne) ClearURL() *SpiderPayloadUpdateOne {
	spuo.mutation.ClearURL()
	return spuo
}

// SetPath sets the "path" field.
func (spuo *SpiderPayloadUpdateOne) SetPath(s string) *SpiderPayloadUpdateOne {
	spuo.mutation.SetPath(s)
	return spuo
}

// SetNillablePath sets the "path" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillablePath(s *string) *SpiderPayloadUpdateOne {
	if s != nil {
		spuo.SetPath(*s)
	}
	return spuo
}

// ClearPath clears the value of the "path" field.
func (spuo *SpiderPayloadUpdateOne) ClearPath() *SpiderPayloadUpdateOne {
	spuo.mutation.ClearPath()
	return spuo
}

// SetStatus sets the "status" field.
func (spuo *SpiderPayloadUpdateOne) SetStatus(u uint8) *SpiderPayloadUpdateOne {
	spuo.mutation.ResetStatus()
	spuo.mutation.SetStatus(u)
	return spuo
}

// SetNillableStatus sets the "status" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillableStatus(u *uint8) *SpiderPayloadUpdateOne {
	if u != nil {
		spuo.SetStatus(*u)
	}
	return spuo
}

// AddStatus adds u to the "status" field.
func (spuo *SpiderPayloadUpdateOne) AddStatus(u int8) *SpiderPayloadUpdateOne {
	spuo.mutation.AddStatus(u)
	return spuo
}

// SetTitle sets the "title" field.
func (spuo *SpiderPayloadUpdateOne) SetTitle(s string) *SpiderPayloadUpdateOne {
	spuo.mutation.SetTitle(s)
	return spuo
}

// SetNillableTitle sets the "title" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillableTitle(s *string) *SpiderPayloadUpdateOne {
	if s != nil {
		spuo.SetTitle(*s)
	}
	return spuo
}

// SetJobProvider sets the "job_provider" field.
func (spuo *SpiderPayloadUpdateOne) SetJobProvider(s string) *SpiderPayloadUpdateOne {
	spuo.mutation.SetJobProvider(s)
	return spuo
}

// SetNillableJobProvider sets the "job_provider" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillableJobProvider(s *string) *SpiderPayloadUpdateOne {
	if s != nil {
		spuo.SetJobProvider(*s)
	}
	return spuo
}

// ClearJobProvider clears the value of the "job_provider" field.
func (spuo *SpiderPayloadUpdateOne) ClearJobProvider() *SpiderPayloadUpdateOne {
	spuo.mutation.ClearJobProvider()
	return spuo
}

// SetJobID sets the "job_id" field.
func (spuo *SpiderPayloadUpdateOne) SetJobID(u uuid.UUID) *SpiderPayloadUpdateOne {
	spuo.mutation.SetJobID(u)
	return spuo
}

// SetNillableJobID sets the "job_id" field if the given value is not nil.
func (spuo *SpiderPayloadUpdateOne) SetNillableJobID(u *uuid.UUID) *SpiderPayloadUpdateOne {
	if u != nil {
		spuo.SetJobID(*u)
	}
	return spuo
}

// ClearJobID clears the value of the "job_id" field.
func (spuo *SpiderPayloadUpdateOne) ClearJobID() *SpiderPayloadUpdateOne {
	spuo.mutation.ClearJobID()
	return spuo
}

// Mutation returns the SpiderPayloadMutation object of the builder.
func (spuo *SpiderPayloadUpdateOne) Mutation() *SpiderPayloadMutation {
	return spuo.mutation
}

// Where appends a list predicates to the SpiderPayloadUpdate builder.
func (spuo *SpiderPayloadUpdateOne) Where(ps ...predicate.SpiderPayload) *SpiderPayloadUpdateOne {
	spuo.mutation.Where(ps...)
	return spuo
}

// Select allows selecting one or more fields (columns) of the returned entity.
// The default is selecting all fields defined in the entity schema.
func (spuo *SpiderPayloadUpdateOne) Select(field string, fields ...string) *SpiderPayloadUpdateOne {
	spuo.fields = append([]string{field}, fields...)
	return spuo
}

// Save executes the query and returns the updated SpiderPayload entity.
func (spuo *SpiderPayloadUpdateOne) Save(ctx context.Context) (*SpiderPayload, error) {
	return withHooks(ctx, spuo.sqlSave, spuo.mutation, spuo.hooks)
}

// SaveX is like Save, but panics if an error occurs.
func (spuo *SpiderPayloadUpdateOne) SaveX(ctx context.Context) *SpiderPayload {
	node, err := spuo.Save(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// Exec executes the query on the entity.
func (spuo *SpiderPayloadUpdateOne) Exec(ctx context.Context) error {
	_, err := spuo.Save(ctx)
	return err
}

// ExecX is like Exec, but panics if an error occurs.
func (spuo *SpiderPayloadUpdateOne) ExecX(ctx context.Context) {
	if err := spuo.Exec(ctx); err != nil {
		panic(err)
	}
}

func (spuo *SpiderPayloadUpdateOne) sqlSave(ctx context.Context) (_node *SpiderPayload, err error) {
	_spec := sqlgraph.NewUpdateSpec(spiderpayload.Table, spiderpayload.Columns, sqlgraph.NewFieldSpec(spiderpayload.FieldID, field.TypeUUID))
	id, ok := spuo.mutation.ID()
	if !ok {
		return nil, &ValidationError{Name: "id", err: errors.New(`ent: missing "SpiderPayload.id" for update`)}
	}
	_spec.Node.ID.Value = id
	if fields := spuo.fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, spiderpayload.FieldID)
		for _, f := range fields {
			if !spiderpayload.ValidColumn(f) {
				return nil, &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
			}
			if f != spiderpayload.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, f)
			}
		}
	}
	if ps := spuo.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if value, ok := spuo.mutation.SpiderID(); ok {
		_spec.SetField(spiderpayload.FieldSpiderID, field.TypeUUID, value)
	}
	if value, ok := spuo.mutation.ExtractedAt(); ok {
		_spec.SetField(spiderpayload.FieldExtractedAt, field.TypeTime, value)
	}
	if value, ok := spuo.mutation.URL(); ok {
		_spec.SetField(spiderpayload.FieldURL, field.TypeString, value)
	}
	if spuo.mutation.URLCleared() {
		_spec.ClearField(spiderpayload.FieldURL, field.TypeString)
	}
	if value, ok := spuo.mutation.Path(); ok {
		_spec.SetField(spiderpayload.FieldPath, field.TypeString, value)
	}
	if spuo.mutation.PathCleared() {
		_spec.ClearField(spiderpayload.FieldPath, field.TypeString)
	}
	if value, ok := spuo.mutation.Status(); ok {
		_spec.SetField(spiderpayload.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := spuo.mutation.AddedStatus(); ok {
		_spec.AddField(spiderpayload.FieldStatus, field.TypeUint8, value)
	}
	if value, ok := spuo.mutation.Title(); ok {
		_spec.SetField(spiderpayload.FieldTitle, field.TypeString, value)
	}
	if value, ok := spuo.mutation.JobProvider(); ok {
		_spec.SetField(spiderpayload.FieldJobProvider, field.TypeString, value)
	}
	if spuo.mutation.JobProviderCleared() {
		_spec.ClearField(spiderpayload.FieldJobProvider, field.TypeString)
	}
	if value, ok := spuo.mutation.JobID(); ok {
		_spec.SetField(spiderpayload.FieldJobID, field.TypeUUID, value)
	}
	if spuo.mutation.JobIDCleared() {
		_spec.ClearField(spiderpayload.FieldJobID, field.TypeUUID)
	}
	_node = &SpiderPayload{config: spuo.config}
	_spec.Assign = _node.assignValues
	_spec.ScanValues = _node.scanValues
	if err = sqlgraph.UpdateNode(ctx, spuo.driver, _spec); err != nil {
		if _, ok := err.(*sqlgraph.NotFoundError); ok {
			err = &NotFoundError{spiderpayload.Label}
		} else if sqlgraph.IsConstraintError(err) {
			err = &ConstraintError{msg: err.Error(), wrap: err}
		}
		return nil, err
	}
	spuo.mutation.done = true
	return _node, nil
}
