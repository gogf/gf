// Copyright GoFrame Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"time"
)

// Cache sets the cache feature for the model. It caches the result of the sql, which means
// if there's another same sql request, it just reads and returns the result from cache, it
// but not committed and executed into the database.
//
// If the parameter <duration> < 0, which means it clear the cache with given <name>.
// If the parameter <duration> = 0, which means it never expires.
// If the parameter <duration> > 0, which means it expires after <duration>.
//
// The optional parameter <name> is used to bind a name to the cache, which means you can
// later control the cache like changing the <duration> or clearing the cache with specified
// <name>.
//
// Note that, the cache feature is disabled if the model is performing select statement
// on a transaction.
func (m *Model) Cache(duration time.Duration, name ...string) *Model {
	model := m.getModel()
	model.cacheDuration = duration
	if len(name) > 0 {
		model.cacheName = name[0]
	}
	model.cacheEnabled = true
	return model
}

// checkAndRemoveCache checks and removes the cache in insert/update/delete statement if
// cache feature is enabled.
func (m *Model) checkAndRemoveCache() {
	if m.cacheEnabled && m.cacheDuration < 0 && len(m.cacheName) > 0 {
		m.db.GetCache().Remove(m.cacheName)
	}
}
