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

// FormatPartitionClause returns the SQL clause selecting the given table partitions.
// MySQL accepts a comma-separated list, like "PARTITION (p0,p1)"; drivers embedding
// *mysql.Driver, such as MariaDB, OceanBase and TiDB, inherit this.
func (d *Driver) FormatPartitionClause(partitions []string) string {
	return fmt.Sprintf("PARTITION (%s)", gstr.Join(partitions, ","))
}
