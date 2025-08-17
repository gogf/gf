// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
//

package gdb

import (
	"context"
	"database/sql"
)

type localStatsItem struct {
	node  *ConfigNode
	stats sql.DBStats
}

// Node returns the configuration node info.
func (item *localStatsItem) Node() ConfigNode {
	return *item.node
}

// Stats returns the connection stat for current node.
func (item *localStatsItem) Stats() sql.DBStats {
	return item.stats
}

// Stats retrieves and returns the pool stat for all nodes that have been established.
func (c *Core) Stats(ctx context.Context) []StatsItem {
	var items = make([]StatsItem, 0)
	c.links.Iterator(func(k, v any) bool {
		var (
			node  = k.(ConfigNode)
			sqlDB = v.(*sql.DB)
		)
		items = append(items, &localStatsItem{
			node:  &node,
			stats: sqlDB.Stats(),
		})
		return true
	})
	return items
}
