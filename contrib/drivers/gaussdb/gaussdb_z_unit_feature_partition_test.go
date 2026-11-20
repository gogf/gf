// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gaussdb_test

import (
	"testing"
)

// GaussDB does support partitioned tables and the `FROM t PARTITION (name)`
// selection clause, accepting a single partition name per clause.
//
// These tests are skipped because the GaussDB driver does not override
// DB.FormatPartitionClause: Model.Partition() names are therefore left unused
// and never reach the generated SQL. The driver cannot render them yet, as it
// rejects the comma-separated list that Model.Partition() accepts, and the clause
// formatter has no way to report that per-call.
//
// The function names are kept aligned with the MySQL/MariaDB baseline so the
// cases can be filled in once the driver implements FormatPartitionClause.

func Test_Partition_Range_Insert_And_Query(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}

func Test_Partition_Range_PartitionQuery(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}

func Test_Partition_Hash_Insert_And_Distribution(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}

func Test_Partition_List_Insert_And_Query(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}

func Test_Partition_Range_Update(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}

func Test_Partition_Range_Delete(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}

func Test_Partition_Transaction(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}

func Test_Partition_Range_Count_And_Sum(t *testing.T) {
	t.Skip("GaussDB driver does not support Model.Partition() yet")
}
