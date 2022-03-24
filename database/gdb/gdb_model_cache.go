// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"time"

	"github.com/gogf/gf/v2/internal/intlog"
)

type CacheOption struct {
	// Duration is the TTL for the cache.
	// If the parameter `Duration` < 0, which means it clear the cache with given `Name`.
	// If the parameter `Duration` = 0, which means it never expires.
	// If the parameter `Duration` > 0, which means it expires after `Duration`.
	Duration time.Duration

	// Name is an optional unique name for the cache.
	// The Name is used to bind a name to the cache, which means you can later control the cache
	// like changing the `duration` or clearing the cache with specified Name.
	Name string

	// Force caches the query result whatever the result is nil or not.
	// It is used to avoid Cache Penetration.
	Force bool
}

// Cache sets the cache feature for the model. It caches the result of the sql, which means
// if there's another same sql request, it just reads and returns the result from cache, it
// but not committed and executed into the database.
//
// Note that, the cache feature is disabled if the model is performing select statement
// on a transaction.
func (m *Model) Cache(option CacheOption) *Model {
	model := m.getModel()
	model.cacheOption = option
	model.cacheEnabled = true
	return model
}

// checkAndRemoveCache checks and removes the cache in insert/update/delete statement if
// cache feature is enabled.
func (m *Model) checkAndRemoveCache(ctx context.Context) {
	if m.cacheEnabled && m.cacheOption.Duration < 0 && len(m.cacheOption.Name) > 0 {
		if _, err := m.db.GetCache().Remove(ctx, m.cacheOption.Name); err != nil {
			intlog.Errorf(ctx, `%+v`, err)
		}
	}
}
