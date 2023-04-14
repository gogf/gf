// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
)

// Transaction wraps the transaction logic using function `f`.
// It rollbacks the transaction and returns the error from function `f` if
// it returns non-nil error. It commits the transaction and returns nil if
// function `f` returns nil.
//
// Note that, you should not Commit or Rollback the transaction in function `f`
// as it is automatically handled by this function.
func (m *Model) Transaction(ctx context.Context, f func(ctx context.Context, tx TX) error) (err error) {
	if ctx == nil {
		ctx = m.GetCtx()
	}
	if m.tx != nil {
		return m.tx.Transaction(ctx, f)
	}
	return m.db.Transaction(ctx, f)
}
