// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql

import (
	"fmt"

	"github.com/gogf/gf/v2/text/gstr"
)

// FormatPartitionClause returns the SQL clause selecting the given table partitions, like
// "PARTITION (`p0`,`p1`)"; drivers embedding *mysql.Driver, such as MariaDB, OceanBase and
// TiDB, inherit this.
//
// Each name is emitted as a quoted identifier with an embedded backtick doubled, so a name
// that carries ")" or "`" stays inside the clause instead of ending it. Names are taken as
// bare identifiers, so pass "p0" and not "`p0`".
//
// An empty name is emitted as an empty identifier, which the server rejects: leaving it out
// instead would silently read every partition, the very behaviour Partition() exists to stop.
func (d *Driver) FormatPartitionClause(partitions []string) string {
	if len(partitions) == 0 {
		return ""
	}
	var names = make([]string, len(partitions))
	for i, partition := range partitions {
		names[i] = "`" + gstr.Replace(gstr.Trim(partition), "`", "``") + "`"
	}
	return fmt.Sprintf("PARTITION (%s)", gstr.Join(names, ","))
}
