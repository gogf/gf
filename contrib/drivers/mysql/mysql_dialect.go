// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql

import (
	"github.com/gogf/gf/v2/text/gstr"
)

// FormatPartitionClause returns the MySQL PARTITION clause.
// MySQL supports PARTITION (p0,p1) syntax for partition pruning.
func (d *Driver) FormatPartitionClause(table string, partitions []string) string {
	if len(partitions) == 0 {
		return table
	}
	return table + " PARTITION (" + gstr.Join(partitions, ",") + ")"
}
