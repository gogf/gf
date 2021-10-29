// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"github.com/gogf/gf/v2/text/gstr"
	"github.com/gogf/gf/v2/util/gconv"
	"strings"
)

// Where sets the condition statement for the model. The parameter `where` can be type of
// string/map/gmap/slice/struct/*struct, etc. Note that, if it's called more than one times,
// multiple conditions will be joined into where statement using "AND".
// Eg:
// Where("uid=10000")
// Where("uid", 10000)
// Where("money>? AND name like ?", 99999, "vip_%")
// Where("uid", 1).Where("name", "john")
// Where("status IN (?)", g.Slice{1,2,3})
// Where("age IN(?,?)", 18, 50)
// Where(User{ Id : 1, UserName : "john"})
func (m *Model) Where(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]ModelWhereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, ModelWhereHolder{
		Operator: whereHolderOperatorWhere,
		Where:    where,
		Args:     args,
	})
	return model
}

// Having sets the having statement for the model.
// The parameters of this function usage are as the same as function Where.
// See Where.
func (m *Model) Having(having interface{}, args ...interface{}) *Model {
	model := m.getModel()
	model.having = []interface{}{
		having, args,
	}
	return model
}

// WherePri does the same logic as Model.Where except that if the parameter `where`
// is a single condition like int/string/float/slice, it treats the condition as the primary
// key value. That is, if primary key is "id" and given `where` parameter as "123", the
// WherePri function treats the condition as "id=123", but Model.Where treats the condition
// as string "123".
func (m *Model) WherePri(where interface{}, args ...interface{}) *Model {
	if len(args) > 0 {
		return m.Where(where, args...)
	}
	newWhere := GetPrimaryKeyCondition(m.getPrimaryKey(), where)
	return m.Where(newWhere[0], newWhere[1:]...)
}

// Wheref builds condition string using fmt.Sprintf and arguments.
// Note that if the number of `args` is more than the placeholder in `format`,
// the extra `args` will be used as the where condition arguments of the Model.
func (m *Model) Wheref(format string, args ...interface{}) *Model {
	var (
		placeHolderCount = gstr.Count(format, "?")
		conditionStr     = fmt.Sprintf(format, args[:len(args)-placeHolderCount]...)
	)
	return m.Where(conditionStr, args[len(args)-placeHolderCount:]...)
}

// WhereLT builds `column < value` statement.
func (m *Model) WhereLT(column string, value interface{}) *Model {
	return m.Wheref(`%s < ?`, column, value)
}

// WhereLTE builds `column <= value` statement.
func (m *Model) WhereLTE(column string, value interface{}) *Model {
	return m.Wheref(`%s <= ?`, column, value)
}

// WhereGT builds `column > value` statement.
func (m *Model) WhereGT(column string, value interface{}) *Model {
	return m.Wheref(`%s > ?`, column, value)
}

// WhereGTE builds `column >= value` statement.
func (m *Model) WhereGTE(column string, value interface{}) *Model {
	return m.Wheref(`%s >= ?`, column, value)
}

// WhereBetween builds `column BETWEEN min AND max` statement.
func (m *Model) WhereBetween(column string, min, max interface{}) *Model {
	return m.Wheref(`%s BETWEEN ? AND ?`, m.db.GetCore().QuoteWord(column), min, max)
}

// WhereLike builds `column LIKE like` statement.
func (m *Model) WhereLike(column string, like interface{}) *Model {
	return m.Wheref(`%s LIKE ?`, m.db.GetCore().QuoteWord(column), like)
}

// WhereIn builds `column IN (in)` statement.
func (m *Model) WhereIn(column string, in interface{}) *Model {
	return m.Wheref(`%s IN (?)`, m.db.GetCore().QuoteWord(column), in)
}

// WhereNull builds `columns[0] IS NULL AND columns[1] IS NULL ...` statement.
func (m *Model) WhereNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.Wheref(`%s IS NULL`, m.db.GetCore().QuoteWord(column))
	}
	return model
}

// WhereNotBetween builds `column NOT BETWEEN min AND max` statement.
func (m *Model) WhereNotBetween(column string, min, max interface{}) *Model {
	return m.Wheref(`%s NOT BETWEEN ? AND ?`, m.db.GetCore().QuoteWord(column), min, max)
}

// WhereNotLike builds `column NOT LIKE like` statement.
func (m *Model) WhereNotLike(column string, like interface{}) *Model {
	return m.Wheref(`%s NOT LIKE ?`, m.db.GetCore().QuoteWord(column), like)
}

// WhereNot builds `column != value` statement.
func (m *Model) WhereNot(column string, value interface{}) *Model {
	return m.Wheref(`%s != ?`, m.db.GetCore().QuoteWord(column), value)
}

// WhereNotIn builds `column NOT IN (in)` statement.
func (m *Model) WhereNotIn(column string, in interface{}) *Model {
	return m.Wheref(`%s NOT IN (?)`, m.db.GetCore().QuoteWord(column), in)
}

// WhereNotNull builds `columns[0] IS NOT NULL AND columns[1] IS NOT NULL ...` statement.
func (m *Model) WhereNotNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.Wheref(`%s IS NOT NULL`, m.db.GetCore().QuoteWord(column))
	}
	return model
}

// WhereOr adds "OR" condition to the where statement.
func (m *Model) WhereOr(where interface{}, args ...interface{}) *Model {
	model := m.getModel()
	if model.whereHolder == nil {
		model.whereHolder = make([]ModelWhereHolder, 0)
	}
	model.whereHolder = append(model.whereHolder, ModelWhereHolder{
		Operator: whereHolderOperatorOr,
		Where:    where,
		Args:     args,
	})
	return model
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
func (m *Model) WhereOrf(format string, args ...interface{}) *Model {
	var (
		placeHolderCount = gstr.Count(format, "?")
		conditionStr     = fmt.Sprintf(format, args[:len(args)-placeHolderCount]...)
	)
	return m.WhereOr(conditionStr, args[len(args)-placeHolderCount:]...)
}

// WhereOrLT builds `column < value` statement in `OR` conditions..
func (m *Model) WhereOrLT(column string, value interface{}) *Model {
	return m.WhereOrf(`%s < ?`, column, value)
}

// WhereOrLTE builds `column <= value` statement in `OR` conditions..
func (m *Model) WhereOrLTE(column string, value interface{}) *Model {
	return m.WhereOrf(`%s <= ?`, column, value)
}

// WhereOrGT builds `column > value` statement in `OR` conditions..
func (m *Model) WhereOrGT(column string, value interface{}) *Model {
	return m.WhereOrf(`%s > ?`, column, value)
}

// WhereOrGTE builds `column >= value` statement in `OR` conditions..
func (m *Model) WhereOrGTE(column string, value interface{}) *Model {
	return m.WhereOrf(`%s >= ?`, column, value)
}

// WhereOrBetween builds `column BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrBetween(column string, min, max interface{}) *Model {
	return m.WhereOrf(`%s BETWEEN ? AND ?`, m.db.GetCore().QuoteWord(column), min, max)
}

// WhereOrLike builds `column LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrLike(column string, like interface{}) *Model {
	return m.WhereOrf(`%s LIKE ?`, m.db.GetCore().QuoteWord(column), like)
}

// WhereOrIn builds `column IN (in)` statement in `OR` conditions.
func (m *Model) WhereOrIn(column string, in interface{}) *Model {
	return m.WhereOrf(`%s IN (?)`, m.db.GetCore().QuoteWord(column), in)
}

// WhereOrNull builds `columns[0] IS NULL OR columns[1] IS NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.WhereOrf(`%s IS NULL`, m.db.GetCore().QuoteWord(column))
	}
	return model
}

// WhereOrNotBetween builds `column NOT BETWEEN min AND max` statement in `OR` conditions.
func (m *Model) WhereOrNotBetween(column string, min, max interface{}) *Model {
	return m.WhereOrf(`%s NOT BETWEEN ? AND ?`, m.db.GetCore().QuoteWord(column), min, max)
}

// WhereOrNotLike builds `column NOT LIKE like` statement in `OR` conditions.
func (m *Model) WhereOrNotLike(column string, like interface{}) *Model {
	return m.WhereOrf(`%s NOT LIKE ?`, m.db.GetCore().QuoteWord(column), like)
}

// WhereOrNotIn builds `column NOT IN (in)` statement.
func (m *Model) WhereOrNotIn(column string, in interface{}) *Model {
	return m.WhereOrf(`%s NOT IN (?)`, m.db.GetCore().QuoteWord(column), in)
}

// WhereOrNotNull builds `columns[0] IS NOT NULL OR columns[1] IS NOT NULL ...` statement in `OR` conditions.
func (m *Model) WhereOrNotNull(columns ...string) *Model {
	model := m
	for _, column := range columns {
		model = m.WhereOrf(`%s IS NOT NULL`, m.db.GetCore().QuoteWord(column))
	}
	return model
}

// Group sets the "GROUP BY" statement for the model.
func (m *Model) Group(groupBy ...string) *Model {
	if len(groupBy) == 0 {
		return m
	}
	model := m.getModel()
	if model.groupBy != "" {
		model.groupBy += ","
	}
	model.groupBy = model.db.GetCore().QuoteString(strings.Join(groupBy, ","))
	return model
}

// Order sets the "ORDER BY" statement for the model.
func (m *Model) Order(orderBy ...string) *Model {
	if len(orderBy) == 0 {
		return m
	}
	model := m.getModel()
	if model.orderBy != "" {
		model.orderBy += ","
	}
	model.orderBy = model.db.GetCore().QuoteString(strings.Join(orderBy, " "))
	return model
}

// OrderAsc sets the "ORDER BY xxx ASC" statement for the model.
func (m *Model) OrderAsc(column string) *Model {
	if len(column) == 0 {
		return m
	}
	return m.Order(column + " ASC")
}

// OrderDesc sets the "ORDER BY xxx DESC" statement for the model.
func (m *Model) OrderDesc(column string) *Model {
	if len(column) == 0 {
		return m
	}
	return m.Order(column + " DESC")
}

// OrderRandom sets the "ORDER BY RANDOM()" statement for the model.
func (m *Model) OrderRandom() *Model {
	model := m.getModel()
	model.orderBy = "RAND()"
	return model
}

// Limit sets the "LIMIT" statement for the model.
// The parameter `limit` can be either one or two number, if passed two number is passed,
// it then sets "LIMIT limit[0],limit[1]" statement for the model, or else it sets "LIMIT limit[0]"
// statement.
func (m *Model) Limit(limit ...int) *Model {
	model := m.getModel()
	switch len(limit) {
	case 1:
		model.limit = limit[0]
	case 2:
		model.start = limit[0]
		model.limit = limit[1]
	}
	return model
}

// Offset sets the "OFFSET" statement for the model.
// It only makes sense for some databases like SQLServer, PostgreSQL, etc.
func (m *Model) Offset(offset int) *Model {
	model := m.getModel()
	model.offset = offset
	return model
}

// Distinct forces the query to only return distinct results.
func (m *Model) Distinct() *Model {
	model := m.getModel()
	model.distinct = "DISTINCT "
	return model
}

// Page sets the paging number for the model.
// The parameter `page` is started from 1 for paging.
// Note that, it differs that the Limit function starts from 0 for "LIMIT" statement.
func (m *Model) Page(page, limit int) *Model {
	model := m.getModel()
	if page <= 0 {
		page = 1
	}
	model.start = (page - 1) * limit
	model.limit = limit
	return model
}

// formatCondition formats where arguments of the model and returns a new condition sql and its arguments.
// Note that this function does not change any attribute value of the `m`.
//
// The parameter `limit1` specifies whether limits querying only one record if m.limit is not set.
func (m *Model) formatCondition(limit1 bool, isCountStatement bool) (conditionWhere string, conditionExtra string, conditionArgs []interface{}) {
	var (
		prefix = ""
	)
	if gstr.Contains(m.tables, " JOIN ") {
		prefix = m.db.GetCore().QuoteWord(m.tablesInit)
	}
	if len(m.whereHolder) > 0 {
		for _, v := range m.whereHolder {
			switch v.Operator {
			case whereHolderOperatorWhere:
				if conditionWhere == "" {
					newWhere, newArgs := formatWhere(m.db, formatWhereInput{
						Where:     v.Where,
						Args:      v.Args,
						OmitNil:   m.option&optionOmitNilWhere > 0,
						OmitEmpty: m.option&optionOmitEmptyWhere > 0,
						Schema:    m.schema,
						Table:     m.tables,
						Prefix:    prefix,
					})
					if len(newWhere) > 0 {
						conditionWhere = newWhere
						conditionArgs = newArgs
					}
					continue
				}
				fallthrough

			case whereHolderOperatorAnd:
				newWhere, newArgs := formatWhere(m.db, formatWhereInput{
					Where:     v.Where,
					Args:      v.Args,
					OmitNil:   m.option&optionOmitNilWhere > 0,
					OmitEmpty: m.option&optionOmitEmptyWhere > 0,
					Schema:    m.schema,
					Table:     m.tables,
					Prefix:    prefix,
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
				newWhere, newArgs := formatWhere(m.db, formatWhereInput{
					Where:     v.Where,
					Args:      v.Args,
					OmitNil:   m.option&optionOmitNilWhere > 0,
					OmitEmpty: m.option&optionOmitEmptyWhere > 0,
					Schema:    m.schema,
					Table:     m.tables,
					Prefix:    prefix,
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
	softDeletingCondition := m.getConditionForSoftDeleting()
	if m.rawSql != "" && conditionWhere != "" {
		if gstr.ContainsI(m.rawSql, " WHERE ") {
			conditionWhere = " AND " + conditionWhere
		} else {
			conditionWhere = " WHERE " + conditionWhere
		}
	} else if !m.unscoped && softDeletingCondition != "" {
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

	// GROUP BY.
	if m.groupBy != "" {
		conditionExtra += " GROUP BY " + m.groupBy
	}
	// HAVING.
	if len(m.having) > 0 {
		havingStr, havingArgs := formatWhere(m.db, formatWhereInput{
			Where:     m.having[0],
			Args:      gconv.Interfaces(m.having[1]),
			OmitNil:   m.option&optionOmitNilWhere > 0,
			OmitEmpty: m.option&optionOmitEmptyWhere > 0,
			Schema:    m.schema,
			Table:     m.tables,
			Prefix:    prefix,
		})
		if len(havingStr) > 0 {
			conditionExtra += " HAVING " + havingStr
			conditionArgs = append(conditionArgs, havingArgs...)
		}
	}
	// ORDER BY.
	if m.orderBy != "" {
		conditionExtra += " ORDER BY " + m.orderBy
	}
	// LIMIT.
	if !isCountStatement {
		if m.limit != 0 {
			if m.start >= 0 {
				conditionExtra += fmt.Sprintf(" LIMIT %d,%d", m.start, m.limit)
			} else {
				conditionExtra += fmt.Sprintf(" LIMIT %d", m.limit)
			}
		} else if limit1 {
			conditionExtra += " LIMIT 1"
		}

		if m.offset >= 0 {
			conditionExtra += fmt.Sprintf(" OFFSET %d", m.offset)
		}
	}

	if m.lockInfo != "" {
		conditionExtra += " " + m.lockInfo
	}
	return
}
