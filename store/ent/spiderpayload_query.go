// Code generated by ent, DO NOT EDIT.

package ent

import (
	"context"
	"fmt"
	"math"

	"entgo.io/ent"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/dialect/sql/sqlgraph"
	"entgo.io/ent/schema/field"
	"github.com/editorpost/spider/store/ent/predicate"
	"github.com/editorpost/spider/store/ent/spiderpayload"
	"github.com/google/uuid"
)

// SpiderPayloadQuery is the builder for querying SpiderPayload entities.
type SpiderPayloadQuery struct {
	config
	ctx        *QueryContext
	order      []spiderpayload.OrderOption
	inters     []Interceptor
	predicates []predicate.SpiderPayload
	// intermediate query (i.e. traversal path).
	sql  *sql.Selector
	path func(context.Context) (*sql.Selector, error)
}

// Where adds a new predicate for the SpiderPayloadQuery builder.
func (spq *SpiderPayloadQuery) Where(ps ...predicate.SpiderPayload) *SpiderPayloadQuery {
	spq.predicates = append(spq.predicates, ps...)
	return spq
}

// Limit the number of records to be returned by this query.
func (spq *SpiderPayloadQuery) Limit(limit int) *SpiderPayloadQuery {
	spq.ctx.Limit = &limit
	return spq
}

// Offset to start from.
func (spq *SpiderPayloadQuery) Offset(offset int) *SpiderPayloadQuery {
	spq.ctx.Offset = &offset
	return spq
}

// Unique configures the query builder to filter duplicate records on query.
// By default, unique is set to true, and can be disabled using this method.
func (spq *SpiderPayloadQuery) Unique(unique bool) *SpiderPayloadQuery {
	spq.ctx.Unique = &unique
	return spq
}

// Order specifies how the records should be ordered.
func (spq *SpiderPayloadQuery) Order(o ...spiderpayload.OrderOption) *SpiderPayloadQuery {
	spq.order = append(spq.order, o...)
	return spq
}

// First returns the first SpiderPayload entity from the query.
// Returns a *NotFoundError when no SpiderPayload was found.
func (spq *SpiderPayloadQuery) First(ctx context.Context) (*SpiderPayload, error) {
	nodes, err := spq.Limit(1).All(setContextOp(ctx, spq.ctx, ent.OpQueryFirst))
	if err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nil, &NotFoundError{spiderpayload.Label}
	}
	return nodes[0], nil
}

// FirstX is like First, but panics if an error occurs.
func (spq *SpiderPayloadQuery) FirstX(ctx context.Context) *SpiderPayload {
	node, err := spq.First(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return node
}

// FirstID returns the first SpiderPayload ID from the query.
// Returns a *NotFoundError when no SpiderPayload ID was found.
func (spq *SpiderPayloadQuery) FirstID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = spq.Limit(1).IDs(setContextOp(ctx, spq.ctx, ent.OpQueryFirstID)); err != nil {
		return
	}
	if len(ids) == 0 {
		err = &NotFoundError{spiderpayload.Label}
		return
	}
	return ids[0], nil
}

// FirstIDX is like FirstID, but panics if an error occurs.
func (spq *SpiderPayloadQuery) FirstIDX(ctx context.Context) uuid.UUID {
	id, err := spq.FirstID(ctx)
	if err != nil && !IsNotFound(err) {
		panic(err)
	}
	return id
}

// Only returns a single SpiderPayload entity found by the query, ensuring it only returns one.
// Returns a *NotSingularError when more than one SpiderPayload entity is found.
// Returns a *NotFoundError when no SpiderPayload entities are found.
func (spq *SpiderPayloadQuery) Only(ctx context.Context) (*SpiderPayload, error) {
	nodes, err := spq.Limit(2).All(setContextOp(ctx, spq.ctx, ent.OpQueryOnly))
	if err != nil {
		return nil, err
	}
	switch len(nodes) {
	case 1:
		return nodes[0], nil
	case 0:
		return nil, &NotFoundError{spiderpayload.Label}
	default:
		return nil, &NotSingularError{spiderpayload.Label}
	}
}

// OnlyX is like Only, but panics if an error occurs.
func (spq *SpiderPayloadQuery) OnlyX(ctx context.Context) *SpiderPayload {
	node, err := spq.Only(ctx)
	if err != nil {
		panic(err)
	}
	return node
}

// OnlyID is like Only, but returns the only SpiderPayload ID in the query.
// Returns a *NotSingularError when more than one SpiderPayload ID is found.
// Returns a *NotFoundError when no entities are found.
func (spq *SpiderPayloadQuery) OnlyID(ctx context.Context) (id uuid.UUID, err error) {
	var ids []uuid.UUID
	if ids, err = spq.Limit(2).IDs(setContextOp(ctx, spq.ctx, ent.OpQueryOnlyID)); err != nil {
		return
	}
	switch len(ids) {
	case 1:
		id = ids[0]
	case 0:
		err = &NotFoundError{spiderpayload.Label}
	default:
		err = &NotSingularError{spiderpayload.Label}
	}
	return
}

// OnlyIDX is like OnlyID, but panics if an error occurs.
func (spq *SpiderPayloadQuery) OnlyIDX(ctx context.Context) uuid.UUID {
	id, err := spq.OnlyID(ctx)
	if err != nil {
		panic(err)
	}
	return id
}

// All executes the query and returns a list of SpiderPayloads.
func (spq *SpiderPayloadQuery) All(ctx context.Context) ([]*SpiderPayload, error) {
	ctx = setContextOp(ctx, spq.ctx, ent.OpQueryAll)
	if err := spq.prepareQuery(ctx); err != nil {
		return nil, err
	}
	qr := querierAll[[]*SpiderPayload, *SpiderPayloadQuery]()
	return withInterceptors[[]*SpiderPayload](ctx, spq, qr, spq.inters)
}

// AllX is like All, but panics if an error occurs.
func (spq *SpiderPayloadQuery) AllX(ctx context.Context) []*SpiderPayload {
	nodes, err := spq.All(ctx)
	if err != nil {
		panic(err)
	}
	return nodes
}

// IDs executes the query and returns a list of SpiderPayload IDs.
func (spq *SpiderPayloadQuery) IDs(ctx context.Context) (ids []uuid.UUID, err error) {
	if spq.ctx.Unique == nil && spq.path != nil {
		spq.Unique(true)
	}
	ctx = setContextOp(ctx, spq.ctx, ent.OpQueryIDs)
	if err = spq.Select(spiderpayload.FieldID).Scan(ctx, &ids); err != nil {
		return nil, err
	}
	return ids, nil
}

// IDsX is like IDs, but panics if an error occurs.
func (spq *SpiderPayloadQuery) IDsX(ctx context.Context) []uuid.UUID {
	ids, err := spq.IDs(ctx)
	if err != nil {
		panic(err)
	}
	return ids
}

// Count returns the count of the given query.
func (spq *SpiderPayloadQuery) Count(ctx context.Context) (int, error) {
	ctx = setContextOp(ctx, spq.ctx, ent.OpQueryCount)
	if err := spq.prepareQuery(ctx); err != nil {
		return 0, err
	}
	return withInterceptors[int](ctx, spq, querierCount[*SpiderPayloadQuery](), spq.inters)
}

// CountX is like Count, but panics if an error occurs.
func (spq *SpiderPayloadQuery) CountX(ctx context.Context) int {
	count, err := spq.Count(ctx)
	if err != nil {
		panic(err)
	}
	return count
}

// Exist returns true if the query has elements in the graph.
func (spq *SpiderPayloadQuery) Exist(ctx context.Context) (bool, error) {
	ctx = setContextOp(ctx, spq.ctx, ent.OpQueryExist)
	switch _, err := spq.FirstID(ctx); {
	case IsNotFound(err):
		return false, nil
	case err != nil:
		return false, fmt.Errorf("ent: check existence: %w", err)
	default:
		return true, nil
	}
}

// ExistX is like Exist, but panics if an error occurs.
func (spq *SpiderPayloadQuery) ExistX(ctx context.Context) bool {
	exist, err := spq.Exist(ctx)
	if err != nil {
		panic(err)
	}
	return exist
}

// Clone returns a duplicate of the SpiderPayloadQuery builder, including all associated steps. It can be
// used to prepare common query builders and use them differently after the clone is made.
func (spq *SpiderPayloadQuery) Clone() *SpiderPayloadQuery {
	if spq == nil {
		return nil
	}
	return &SpiderPayloadQuery{
		config:     spq.config,
		ctx:        spq.ctx.Clone(),
		order:      append([]spiderpayload.OrderOption{}, spq.order...),
		inters:     append([]Interceptor{}, spq.inters...),
		predicates: append([]predicate.SpiderPayload{}, spq.predicates...),
		// clone intermediate query.
		sql:  spq.sql.Clone(),
		path: spq.path,
	}
}

// GroupBy is used to group vertices by one or more fields/columns.
// It is often used with aggregate functions, like: count, max, mean, min, sum.
//
// Example:
//
//	var v []struct {
//		SpiderID uuid.UUID `json:"spider_id,omitempty"`
//		Count int `json:"count,omitempty"`
//	}
//
//	client.SpiderPayload.Query().
//		GroupBy(spiderpayload.FieldSpiderID).
//		Aggregate(ent.Count()).
//		Scan(ctx, &v)
func (spq *SpiderPayloadQuery) GroupBy(field string, fields ...string) *SpiderPayloadGroupBy {
	spq.ctx.Fields = append([]string{field}, fields...)
	grbuild := &SpiderPayloadGroupBy{build: spq}
	grbuild.flds = &spq.ctx.Fields
	grbuild.label = spiderpayload.Label
	grbuild.scan = grbuild.Scan
	return grbuild
}

// Select allows the selection one or more fields/columns for the given query,
// instead of selecting all fields in the entity.
//
// Example:
//
//	var v []struct {
//		SpiderID uuid.UUID `json:"spider_id,omitempty"`
//	}
//
//	client.SpiderPayload.Query().
//		Select(spiderpayload.FieldSpiderID).
//		Scan(ctx, &v)
func (spq *SpiderPayloadQuery) Select(fields ...string) *SpiderPayloadSelect {
	spq.ctx.Fields = append(spq.ctx.Fields, fields...)
	sbuild := &SpiderPayloadSelect{SpiderPayloadQuery: spq}
	sbuild.label = spiderpayload.Label
	sbuild.flds, sbuild.scan = &spq.ctx.Fields, sbuild.Scan
	return sbuild
}

// Aggregate returns a SpiderPayloadSelect configured with the given aggregations.
func (spq *SpiderPayloadQuery) Aggregate(fns ...AggregateFunc) *SpiderPayloadSelect {
	return spq.Select().Aggregate(fns...)
}

func (spq *SpiderPayloadQuery) prepareQuery(ctx context.Context) error {
	for _, inter := range spq.inters {
		if inter == nil {
			return fmt.Errorf("ent: uninitialized interceptor (forgotten import ent/runtime?)")
		}
		if trv, ok := inter.(Traverser); ok {
			if err := trv.Traverse(ctx, spq); err != nil {
				return err
			}
		}
	}
	for _, f := range spq.ctx.Fields {
		if !spiderpayload.ValidColumn(f) {
			return &ValidationError{Name: f, err: fmt.Errorf("ent: invalid field %q for query", f)}
		}
	}
	if spq.path != nil {
		prev, err := spq.path(ctx)
		if err != nil {
			return err
		}
		spq.sql = prev
	}
	return nil
}

func (spq *SpiderPayloadQuery) sqlAll(ctx context.Context, hooks ...queryHook) ([]*SpiderPayload, error) {
	var (
		nodes = []*SpiderPayload{}
		_spec = spq.querySpec()
	)
	_spec.ScanValues = func(columns []string) ([]any, error) {
		return (*SpiderPayload).scanValues(nil, columns)
	}
	_spec.Assign = func(columns []string, values []any) error {
		node := &SpiderPayload{config: spq.config}
		nodes = append(nodes, node)
		return node.assignValues(columns, values)
	}
	for i := range hooks {
		hooks[i](ctx, _spec)
	}
	if err := sqlgraph.QueryNodes(ctx, spq.driver, _spec); err != nil {
		return nil, err
	}
	if len(nodes) == 0 {
		return nodes, nil
	}
	return nodes, nil
}

func (spq *SpiderPayloadQuery) sqlCount(ctx context.Context) (int, error) {
	_spec := spq.querySpec()
	_spec.Node.Columns = spq.ctx.Fields
	if len(spq.ctx.Fields) > 0 {
		_spec.Unique = spq.ctx.Unique != nil && *spq.ctx.Unique
	}
	return sqlgraph.CountNodes(ctx, spq.driver, _spec)
}

func (spq *SpiderPayloadQuery) querySpec() *sqlgraph.QuerySpec {
	_spec := sqlgraph.NewQuerySpec(spiderpayload.Table, spiderpayload.Columns, sqlgraph.NewFieldSpec(spiderpayload.FieldID, field.TypeUUID))
	_spec.From = spq.sql
	if unique := spq.ctx.Unique; unique != nil {
		_spec.Unique = *unique
	} else if spq.path != nil {
		_spec.Unique = true
	}
	if fields := spq.ctx.Fields; len(fields) > 0 {
		_spec.Node.Columns = make([]string, 0, len(fields))
		_spec.Node.Columns = append(_spec.Node.Columns, spiderpayload.FieldID)
		for i := range fields {
			if fields[i] != spiderpayload.FieldID {
				_spec.Node.Columns = append(_spec.Node.Columns, fields[i])
			}
		}
	}
	if ps := spq.predicates; len(ps) > 0 {
		_spec.Predicate = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	if limit := spq.ctx.Limit; limit != nil {
		_spec.Limit = *limit
	}
	if offset := spq.ctx.Offset; offset != nil {
		_spec.Offset = *offset
	}
	if ps := spq.order; len(ps) > 0 {
		_spec.Order = func(selector *sql.Selector) {
			for i := range ps {
				ps[i](selector)
			}
		}
	}
	return _spec
}

func (spq *SpiderPayloadQuery) sqlQuery(ctx context.Context) *sql.Selector {
	builder := sql.Dialect(spq.driver.Dialect())
	t1 := builder.Table(spiderpayload.Table)
	columns := spq.ctx.Fields
	if len(columns) == 0 {
		columns = spiderpayload.Columns
	}
	selector := builder.Select(t1.Columns(columns...)...).From(t1)
	if spq.sql != nil {
		selector = spq.sql
		selector.Select(selector.Columns(columns...)...)
	}
	if spq.ctx.Unique != nil && *spq.ctx.Unique {
		selector.Distinct()
	}
	for _, p := range spq.predicates {
		p(selector)
	}
	for _, p := range spq.order {
		p(selector)
	}
	if offset := spq.ctx.Offset; offset != nil {
		// limit is mandatory for offset clause. We start
		// with default value, and override it below if needed.
		selector.Offset(*offset).Limit(math.MaxInt32)
	}
	if limit := spq.ctx.Limit; limit != nil {
		selector.Limit(*limit)
	}
	return selector
}

// SpiderPayloadGroupBy is the group-by builder for SpiderPayload entities.
type SpiderPayloadGroupBy struct {
	selector
	build *SpiderPayloadQuery
}

// Aggregate adds the given aggregation functions to the group-by query.
func (spgb *SpiderPayloadGroupBy) Aggregate(fns ...AggregateFunc) *SpiderPayloadGroupBy {
	spgb.fns = append(spgb.fns, fns...)
	return spgb
}

// Scan applies the selector query and scans the result into the given value.
func (spgb *SpiderPayloadGroupBy) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, spgb.build.ctx, ent.OpQueryGroupBy)
	if err := spgb.build.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SpiderPayloadQuery, *SpiderPayloadGroupBy](ctx, spgb.build, spgb, spgb.build.inters, v)
}

func (spgb *SpiderPayloadGroupBy) sqlScan(ctx context.Context, root *SpiderPayloadQuery, v any) error {
	selector := root.sqlQuery(ctx).Select()
	aggregation := make([]string, 0, len(spgb.fns))
	for _, fn := range spgb.fns {
		aggregation = append(aggregation, fn(selector))
	}
	if len(selector.SelectedColumns()) == 0 {
		columns := make([]string, 0, len(*spgb.flds)+len(spgb.fns))
		for _, f := range *spgb.flds {
			columns = append(columns, selector.C(f))
		}
		columns = append(columns, aggregation...)
		selector.Select(columns...)
	}
	selector.GroupBy(selector.Columns(*spgb.flds...)...)
	if err := selector.Err(); err != nil {
		return err
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := spgb.build.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}

// SpiderPayloadSelect is the builder for selecting fields of SpiderPayload entities.
type SpiderPayloadSelect struct {
	*SpiderPayloadQuery
	selector
}

// Aggregate adds the given aggregation functions to the selector query.
func (sps *SpiderPayloadSelect) Aggregate(fns ...AggregateFunc) *SpiderPayloadSelect {
	sps.fns = append(sps.fns, fns...)
	return sps
}

// Scan applies the selector query and scans the result into the given value.
func (sps *SpiderPayloadSelect) Scan(ctx context.Context, v any) error {
	ctx = setContextOp(ctx, sps.ctx, ent.OpQuerySelect)
	if err := sps.prepareQuery(ctx); err != nil {
		return err
	}
	return scanWithInterceptors[*SpiderPayloadQuery, *SpiderPayloadSelect](ctx, sps.SpiderPayloadQuery, sps, sps.inters, v)
}

func (sps *SpiderPayloadSelect) sqlScan(ctx context.Context, root *SpiderPayloadQuery, v any) error {
	selector := root.sqlQuery(ctx)
	aggregation := make([]string, 0, len(sps.fns))
	for _, fn := range sps.fns {
		aggregation = append(aggregation, fn(selector))
	}
	switch n := len(*sps.selector.flds); {
	case n == 0 && len(aggregation) > 0:
		selector.Select(aggregation...)
	case n != 0 && len(aggregation) > 0:
		selector.AppendSelect(aggregation...)
	}
	rows := &sql.Rows{}
	query, args := selector.Query()
	if err := sps.driver.Query(ctx, query, args, rows); err != nil {
		return err
	}
	defer rows.Close()
	return sql.ScanSlice(rows, v)
}