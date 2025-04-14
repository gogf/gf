// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "context"

// WhereOr adds "OR" condition to the where statement.
// See WhereBuilder.WhereOr.
func (m *Model) WhereOr(where any, args ...any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOr(where, args...))
	})
}

// WhereOrf builds `OR` condition string using fmt.Sprintf and arguments.
// See WhereBuilder.WhereOrf.
func (m *Model) WhereOrf(format string, args ...any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrf(format, args...))
	})
}

// WhereOrLT builds `column < value` statement in `OR` conditions.
// See WhereBuilder.WhereOrLT.
func (m *Model) WhereOrLT(column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrLT(column, value))
	})
}

// WhereOrLTE builds `column <= value` statement in `OR` conditions.
// See WhereBuilder.WhereOrLTE.
func (m *Model) WhereOrLTE(column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrLTE(column, value))
	})
}

// WhereOrGT builds `column > value` statement in `OR` conditions.
// See WhereBuilder.WhereOrGT.
func (m *Model) WhereOrGT(column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrGT(column, value))
	})
}

// WhereOrGTE builds `column >= value` statement in `OR` conditions.
// See WhereBuilder.WhereOrGTE.
func (m *Model) WhereOrGTE(column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrGTE(column, value))
	})
}

// WhereOrBetween builds `column BETWEEN min AND max` statement in `OR` conditions.
// See WhereBuilder.WhereOrBetween.
func (m *Model) WhereOrBetween(column string, min, max any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrBetween(column, min, max))
	})
}

// WhereOrLike builds `column LIKE like` statement in `OR` conditions.
// See WhereBuilder.WhereOrLike.
func (m *Model) WhereOrLike(column string, like any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrLike(column, like))
	})
}

// WhereOrIn builds `column IN (in)` statement in `OR` conditions.
// See WhereBuilder.WhereOrIn.
func (m *Model) WhereOrIn(column string, in any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrIn(column, in))
	})
}

// WhereOrNull builds `columns[0] IS NULL OR columns[1] IS NULL ...` statement in `OR` conditions.
// See WhereBuilder.WhereOrNull.
func (m *Model) WhereOrNull(columns ...string) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrNull(columns...))
	})
}

// WhereOrNotBetween builds `column NOT BETWEEN min AND max` statement in `OR` conditions.
// See WhereBuilder.WhereOrNotBetween.
func (m *Model) WhereOrNotBetween(column string, min, max any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrNotBetween(column, min, max))
	})
}

// WhereOrNotLike builds `column NOT LIKE 'like'` statement in `OR` conditions.
// See WhereBuilder.WhereOrNotLike.
func (m *Model) WhereOrNotLike(column string, like any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrNotLike(column, like))
	})
}

// WhereOrNot builds `column != value` statement.
// See WhereBuilder.WhereOrNot.
func (m *Model) WhereOrNot(column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrNot(column, value))
	})
}

// WhereOrNotIn builds `column NOT IN (in)` statement.
// See WhereBuilder.WhereOrNotIn.
func (m *Model) WhereOrNotIn(column string, in any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrNotIn(column, in))
	})
}

// WhereOrNotNull builds `columns[0] IS NOT NULL OR columns[1] IS NOT NULL ...` statement in `OR` conditions.
// See WhereBuilder.WhereOrNotNull.
func (m *Model) WhereOrNotNull(columns ...string) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrNotNull(columns...))
	})
}
