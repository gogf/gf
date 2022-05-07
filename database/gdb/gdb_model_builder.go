// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
)

type WhereBuilder struct {
	safe        *bool         // If nil, it uses the safe attribute of its model.
	model       *Model        // A WhereBuilder should be bound to certain Model.
	whereHolder []WhereHolder // Condition strings for where operation.
}

// WhereHolder is the holder for where condition preparing.
type WhereHolder struct {
	Type     string        // Type of this holder.
	Operator int           // Operator for this holder.
	Where    interface{}   // Where parameter, which can commonly be type of string/map/struct.
	Args     []interface{} // Arguments for where parameter.
	Prefix   string        // Field prefix, eg: "user.", "order.".
}

func (m *Model) Builder() *WhereBuilder {
	// The WhereBuilder is safe in default when it is created using Builder().
	var isSafe = true
	b := &WhereBuilder{
		safe:        &isSafe,
		model:       m,
		whereHolder: make([]WhereHolder, 0),
	}
	return b
}

// getBuilder creates and returns a cloned WhereBuilder of current WhereBuilder if `safe` is true,
// or else it returns the current WhereBuilder.
func (b *WhereBuilder) getBuilder() *WhereBuilder {
	var isSafe bool
	if b.safe != nil {
		isSafe = *b.safe
	} else {
		isSafe = b.model.safe
	}
	if !isSafe {
		return b
	} else {
		return b.Clone()
	}
}

func (b *WhereBuilder) Clone() *WhereBuilder {
	newBuilder := b.model.Builder()
	newBuilder.safe = b.safe
	newBuilder.whereHolder = make([]WhereHolder, len(b.whereHolder))
	copy(newBuilder.whereHolder, b.whereHolder)
	return newBuilder
}

func (b *WhereBuilder) Build() (conditionWhere string, conditionArgs []interface{}) {
	var (
		ctx                         = b.model.GetCtx()
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

func (b *WhereBuilder) convertWrappedBuilder(where interface{}, args []interface{}) (newWhere interface{}, newArgs []interface{}) {
	var builder *WhereBuilder
	switch v := where.(type) {
	case WhereBuilder:
		builder = &v
	case *WhereBuilder:
		builder = v
	}
	if builder != nil {
		conditionWhere, conditionArgs := builder.Build()
		if len(b.whereHolder) == 0 {
			conditionWhere = "(" + conditionWhere + ")"
		}
		return conditionWhere, conditionArgs
	}
	return where, args
}
