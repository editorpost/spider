// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"

	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/editorpost/spider/store/ent/predicate"
	"github.com/editorpost/spider/store/ent/spiderpayload"
)

// SpiderPayloadDelete is the builder for deleting a SpiderPayload entity.
type SpiderPayloadDelete struct {
	config
	hooks    []Hook
	mutation *SpiderPayloadMutation
}

// Where appends a list predicates to the SpiderPayloadDelete builder.
func (spd *SpiderPayloadDelete) Where(ps ...predicate.SpiderPayload) *SpiderPayloadDelete {
	spd.mutation.Where(ps...)
	return spd
}

// Exec executes the deletion query and returns how many vertices were deleted.
func (spd *SpiderPayloadDelete) Exec(ctx context.Context) (int, error) {
	return withHooks(ctx, spd.sqlExec, spd.mutation, spd.hooks)
}

// ExecX is like Exec, but panics if an error occurs.
func (spd *SpiderPayloadDelete) ExecX(ctx context.Context) int {
	n, err := spd.Exec(ctx)
	if err != nil {
		panic(err)
	}
	return n
}

func (spd *SpiderPayloadDelete) sqlExec(ctx context.Context) (int, error) {
	_spec := sqlgraph.NewDeleteSpec(spiderpayload.Table, sqlgraph.NewFieldSpec(spiderpayload.FieldID, field.TypeUUID))
	if ps := spd.mutation.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	affected, err := sqlgraph.DeleteNodes(ctx, spd.driver, _spec)
	if err != nil && sqlgraph.IsConstraintError(err) {
		err = &ConstraintError{msg: err.Error(), wrap: err}
	}
	spd.mutation.done = true
	return affected, err
}

// SpiderPayloadDeleteOne is the builder for deleting a single SpiderPayload entity.
type SpiderPayloadDeleteOne struct {
	spd *SpiderPayloadDelete
}

// Where appends a list predicates to the SpiderPayloadDelete builder.
func (spdo *SpiderPayloadDeleteOne) Where(ps ...predicate.SpiderPayload) *SpiderPayloadDeleteOne {
	spdo.spd.mutation.Where(ps...)
	return spdo
}

// Exec executes the deletion query.
func (spdo *SpiderPayloadDeleteOne) Exec(ctx context.Context) error {
	n, err := spdo.spd.Exec(ctx)
	switch {
	case err != nil:
		return err
	case n == 0:
		return &NotFoundError{spiderpayload.Label}
	default:
		return nil
	}
}

// ExecX is like Exec, but panics if an error occurs.
func (spdo *SpiderPayloadDeleteOne) ExecX(ctx context.Context) {
	if err := spdo.Exec(ctx); err != nil {
		panic(err)
	}
}
