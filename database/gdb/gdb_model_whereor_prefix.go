// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "context"

// WhereOrPrefix performs as WhereOr, but it adds prefix to each field in where statement.
// See WhereBuilder.WhereOrPrefix.
func (m *Model) WhereOrPrefix(prefix string, where any, args ...any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefix(prefix, where, args...))
	})
}

// WhereOrPrefixLT builds `prefix.column < value` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixLT.
func (m *Model) WhereOrPrefixLT(prefix string, column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixLT(prefix, column, value))
	})
}

// WhereOrPrefixLTE builds `prefix.column <= value` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixLTE.
func (m *Model) WhereOrPrefixLTE(prefix string, column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixLTE(prefix, column, value))
	})
}

// WhereOrPrefixGT builds `prefix.column > value` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixGT.
func (m *Model) WhereOrPrefixGT(prefix string, column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixGT(prefix, column, value))
	})
}

// WhereOrPrefixGTE builds `prefix.column >= value` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixGTE.
func (m *Model) WhereOrPrefixGTE(prefix string, column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixGTE(prefix, column, value))
	})
}

// WhereOrPrefixBetween builds `prefix.column BETWEEN min AND max` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixBetween.
func (m *Model) WhereOrPrefixBetween(prefix string, column string, min, max any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixBetween(prefix, column, min, max))
	})
}

// WhereOrPrefixLike builds `prefix.column LIKE like` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixLike.
func (m *Model) WhereOrPrefixLike(prefix string, column string, like any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixLike(prefix, column, like))
	})
}

// WhereOrPrefixIn builds `prefix.column IN (in)` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixIn.
func (m *Model) WhereOrPrefixIn(prefix string, column string, in any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixIn(prefix, column, in))
	})
}

// WhereOrPrefixNull builds `prefix.columns[0] IS NULL OR prefix.columns[1] IS NULL ...` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixNull.
func (m *Model) WhereOrPrefixNull(prefix string, columns ...string) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixNull(prefix, columns...))
	})
}

// WhereOrPrefixNotBetween builds `prefix.column NOT BETWEEN min AND max` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixNotBetween.
func (m *Model) WhereOrPrefixNotBetween(prefix string, column string, min, max any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixNotBetween(prefix, column, min, max))
	})
}

// WhereOrPrefixNotLike builds `prefix.column NOT LIKE like` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixNotLike.
func (m *Model) WhereOrPrefixNotLike(prefix string, column string, like any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixNotLike(prefix, column, like))
	})
}

// WhereOrPrefixNotIn builds `prefix.column NOT IN (in)` statement.
// See WhereBuilder.WhereOrPrefixNotIn.
func (m *Model) WhereOrPrefixNotIn(prefix string, column string, in any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixNotIn(prefix, column, in))
	})
}

// WhereOrPrefixNotNull builds `prefix.columns[0] IS NOT NULL OR prefix.columns[1] IS NOT NULL ...` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixNotNull.
func (m *Model) WhereOrPrefixNotNull(prefix string, columns ...string) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixNotNull(prefix, columns...))
	})
}

// WhereOrPrefixNot builds `prefix.column != value` statement in `OR` conditions.
// See WhereBuilder.WhereOrPrefixNot.
func (m *Model) WhereOrPrefixNot(prefix string, column string, value any) *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		return model.callWhereBuilder(model.whereBuilder.WhereOrPrefixNot(prefix, column, value))
	})
}
