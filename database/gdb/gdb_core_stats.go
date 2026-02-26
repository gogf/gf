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
	c.links.Iterator(func(k ConfigNode, v *sql.DB) bool {
		// Create a local copy of k to avoid loop variable address re-use issue
		// In Go, loop variables are re-used with the same memory address across iterations,
		// directly using &k would cause all localStatsItem instances to share the same address
		node := k
		items = append(items, &localStatsItem{
			node:  &node,
			stats: v.Stats(),
		})
		return true
	})
	return items
}
