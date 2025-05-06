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
	model        *Model           // A WhereBuilder should be bound to certain Model.
	whereHolder  []WhereHolder    // Condition strings for where operation.
	handlers     []BuilderHandler // Handlers for where operation.
	handlerIndex int
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

// Clone clones and returns a WhereBuilder that is a copy of current one.
func (b *WhereBuilder) Clone() *WhereBuilder {
	newBuilder := &WhereBuilder{
		model:    b.model,
		handlers: make([]BuilderHandler, len(b.handlers)),
	}
	copy(newBuilder.handlers, b.handlers)
	return newBuilder
}

// Handler registers handlers for where operation.
func (b *WhereBuilder) Handler(handlers ...BuilderHandler) *WhereBuilder {
	b.handlers = append(b.handlers, handlers...)
	return b
}

func (b *WhereBuilder) callHandlers(ctx context.Context) *WhereBuilder {
	var (
		builder           = b
		oldHandlersLength = len(builder.handlers)
		newHandlersLength = oldHandlersLength
	)
	for {
		// Exit the loop if all handlers have been processed
		if builder.handlerIndex >= len(builder.handlers) {
			break
		}

		// Record the current length of handlers
		oldHandlersLength = len(builder.handlers)

		// Execute the current handler
		builder = builder.handlers[builder.handlerIndex](ctx, builder)

		// Check if new handlers were added
		newHandlersLength = len(builder.handlers)
		if newHandlersLength > oldHandlersLength {
			var (
				addedCount = newHandlersLength - oldHandlersLength
				targetPos  = builder.handlerIndex + 1
			)

			// Insert newly added handlers into the target position using element swapping technique
			// Example of the swapping logic:
			// 1. We have an array of digits: 123456
			// 2. We're at position 2 and add two new digits 7,8 at the end: 12345678
			// 3. We want to insert these new digits after position 2, so we perform these swaps:
			//    - Swap 3 and 7: 12745638
			//    - Swap 4 and 8: 12785634
			//    - Swap 5 and 3: 12783654
			//    - Swap 6 and 4: 12783456
			// 4. Result: 12783456 - new elements are inserted after position 2
			for i := 0; i < addedCount; i++ {
				// Start from each new element and swap it with preceding elements until it reaches target position
				newItemPos := oldHandlersLength + i
				for j := newItemPos; j > targetPos+i; j-- {
					// Swap elements at positions j and j-1
					builder.handlers[j], builder.handlers[j-1] = builder.handlers[j-1], builder.handlers[j]
				}
			}
		}

		// Move to the next handler
		builder.handlerIndex++
	}
	return builder
}

func (b *WhereBuilder) doCallHandlers(ctx context.Context, builder *WhereBuilder, handlers []BuilderHandler) *WhereBuilder {
	if len(builder.handlers) == 0 {
		return builder
	}
	for builder.handlerIndex < len(builder.handlers) {
		builder = builder.handlers[builder.handlerIndex](ctx, builder)
		builder.handlerIndex++
	}
	return builder
}

// Build builds current WhereBuilder and returns the condition string and parameters.
func (b *WhereBuilder) Build(ctx context.Context) (conditionWhere string, conditionArgs []any) {
	var (
		model                       = b.model.callHandlers(ctx)
		builder                     = b.callHandlers(ctx)
		autoPrefix                  = model.getAutoPrefix()
		tableForMappingAndFiltering = model.tables
	)
	if len(builder.whereHolder) > 0 {
		for _, holder := range builder.whereHolder {
			if holder.Prefix == "" {
				holder.Prefix = autoPrefix
			}
			switch holder.Operator {
			case whereHolderOperatorWhere, whereHolderOperatorAnd:
				newWhere, newArgs := formatWhereHolder(ctx, model.db, formatWhereHolderInput{
					WhereHolder: holder,
					OmitNil:     model.option&optionOmitNilWhere > 0,
					OmitEmpty:   model.option&optionOmitEmptyWhere > 0,
					Schema:      model.schema,
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
				newWhere, newArgs := formatWhereHolder(ctx, model.db, formatWhereHolderInput{
					WhereHolder: holder,
					OmitNil:     model.option&optionOmitNilWhere > 0,
					OmitEmpty:   model.option&optionOmitEmptyWhere > 0,
					Schema:      model.schema,
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
