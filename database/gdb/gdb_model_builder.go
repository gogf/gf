// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

type WhereBuilder struct {
	model       *Model        // Bound to parent model.
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
	b := &WhereBuilder{
		model:       m,
		whereHolder: make([]WhereHolder, 0),
	}
	return b
}

// getBuilder creates and returns a cloned WhereBuilder of current WhereBuilder if `safe` is true,
// or else it returns the current WhereBuilder.
func (b *WhereBuilder) getBuilder() *WhereBuilder {
	if !b.model.safe {
		return b
	} else {
		return b.Clone()
	}
}

func (b *WhereBuilder) Clone() *WhereBuilder {
	newBuilder := b.model.Builder()
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
			case whereHolderOperatorWhere:
				if conditionWhere == "" {
					newWhere, newArgs := formatWhereHolder(ctx, b.model.db, formatWhereHolderInput{
						WhereHolder: holder,
						OmitNil:     b.model.option&optionOmitNilWhere > 0,
						OmitEmpty:   b.model.option&optionOmitEmptyWhere > 0,
						Schema:      b.model.schema,
						Table:       tableForMappingAndFiltering,
					})
					if len(newWhere) > 0 {
						conditionWhere = newWhere
						conditionArgs = newArgs
					}
					continue
				}
				fallthrough

			case whereHolderOperatorAnd:
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
	// Soft deletion.
	softDeletingCondition := b.model.getConditionForSoftDeleting()
	if b.model.rawSql != "" && conditionWhere != "" {
		if gstr.ContainsI(b.model.rawSql, " WHERE ") {
			conditionWhere = " AND " + conditionWhere
		} else {
			conditionWhere = " WHERE " + conditionWhere
		}
	} else if !b.model.unscoped && softDeletingCondition != "" {
		if conditionWhere == "" {
			conditionWhere = fmt.Sprintf(` WHERE %s`, softDeletingCondition)
		} else {
			conditionWhere = fmt.Sprintf(` WHERE (%s) AND %s`, conditionWhere, softDeletingCondition)
		}
	} else {
		if conditionWhere != "" {
			conditionWhere = " WHERE " + conditionWhere
		}
	}
	return
}
