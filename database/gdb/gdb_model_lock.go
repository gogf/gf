// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import "context"

// LockUpdate sets the lock for update for current operation.
func (m *Model) LockUpdate() *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.lockInfo = "FOR UPDATE"
		return model
	})
}

// LockShared sets the lock in share mode for current operation.
func (m *Model) LockShared() *Model {
	return m.Handler(func(ctx context.Context, model *Model) *Model {
		model.lockInfo = "LOCK IN SHARE MODE"
		return model
	})
}
