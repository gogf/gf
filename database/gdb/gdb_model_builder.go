// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
)

// WhereBuilder holds multiple where conditions in a group.
type WhereBuilder struct {
	model       *Model           // A WhereBuilder should be bound to certain Model.
	whereHolder []WhereHolder    // Condition strings for where operation.
	handlers    []BuilderHandler // Handlers for where operation.
}

// BuilderHandler is the handler for where operation.
type BuilderHandler func(ctx context.Context, builder *WhereBuilder) *WhereBuilder

// WhereHolder is the holder for where condition preparing.
type WhereHolder struct {
	Type     string // Type of this holder.
	Operator int    // Operator for this holder.
	Where    any    // Where parameter, which can commonly be type of string/map/struct.
	Args     []any  // Arguments for where parameter.
	Prefix   string // Field prefix, eg: "user.", "order.".
}

// Builder creates and returns a WhereBuilder.
func (m *Model) Builder() *WhereBuilder {
	b := &WhereBuilder{
		model:    m,
		handlers: make([]BuilderHandler, 0),
	}
	return b
}

// Handler registers handlers for where operation.
func (b *WhereBuilder) Handler(handlers ...BuilderHandler) *WhereBuilder {
	b.handlers = append(b.handlers, handlers...)
	return b
}

func (b *WhereBuilder) execHandlers(ctx context.Context) {
	for _, handler := range b.handlers {
		handler(ctx, b)
	}
}

// Clone clones and returns a WhereBuilder that is a copy of current one.
func (b *WhereBuilder) Clone() *WhereBuilder {
	newBuilder := b.model.Builder()
	copy(newBuilder.handlers, b.handlers)
	return newBuilder
}

// Build builds current WhereBuilder and returns the condition string and parameters.
func (b *WhereBuilder) Build(ctx context.Context) (conditionWhere string, conditionArgs []any) {
	// use Clone here to guarantee repeatable build.
	b.Clone().execHandlers(ctx)

	var (
		autoPrefix                  = b.model.getAutoPrefix()
		tableForMappingAndFiltering = b.model.tables
	)
	if len(b.whereHolder) > 0 {
		for _, holder := range b.whereHolder {
			if holder.Prefix == "" {
				holder.Prefix = autoPrefix
			}
			switch holder.Operator {
			case whereHolderOperatorWhere, whereHolderOperatorAnd:
				newWhere, newArgs := formatWhereHolder(ctx, b.model.db, formatWhereHolderInput{
					WhereHolder: holder,
					OmitNil:     b.model.option&optionOmitNilWhere > 0,
					OmitEmpty:   b.model.option&optionOmitEmptyWhere > 0,
					Schema:      b.model.schema,
					Table:       tableForMappingAndFiltering,
				})
				if len(newWhere) > 0 {
					if len(conditionWhere) == 0 {
						conditionWhere = newWhere
					} else if conditionWhere[0] == '(' {
						conditionWhere = fmt.Sprintf(`%s AND (%s)`, conditionWhere, newWhere)
					} else {
						conditionWhere = fmt.Sprintf(`(%s) AND (%s)`, conditionWhere, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}

			case whereHolderOperatorOr:
				newWhere, newArgs := formatWhereHolder(ctx, b.model.db, formatWhereHolderInput{
					WhereHolder: holder,
					OmitNil:     b.model.option&optionOmitNilWhere > 0,
					OmitEmpty:   b.model.option&optionOmitEmptyWhere > 0,
					Schema:      b.model.schema,
					Table:       tableForMappingAndFiltering,
				})
				if len(newWhere) > 0 {
					if len(conditionWhere) == 0 {
						conditionWhere = newWhere
					} else if conditionWhere[0] == '(' {
						conditionWhere = fmt.Sprintf(`%s OR (%s)`, conditionWhere, newWhere)
					} else {
						conditionWhere = fmt.Sprintf(`(%s) OR (%s)`, conditionWhere, newWhere)
					}
					conditionArgs = append(conditionArgs, newArgs...)
				}
			}
		}
	}
	return
}

// convertWhereBuilder converts parameter `where` to condition string and parameters if `where` is also a WhereBuilder.
func (b *WhereBuilder) convertWhereBuilder(ctx context.Context, where any, args []any) (newWhere any, newArgs []any) {
	var builder *WhereBuilder
	switch v := where.(type) {
	case WhereBuilder:
		builder = &v

	case *WhereBuilder:
		builder = v
	}
	if builder != nil {
		conditionWhere, conditionArgs := builder.Build(ctx)
		if conditionWhere != "" && (len(b.whereHolder) == 0 || len(builder.whereHolder) > 1) {
			conditionWhere = "(" + conditionWhere + ")"
		}
		return conditionWhere, conditionArgs
	}
	return where, args
}
