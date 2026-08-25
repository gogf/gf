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
// These tests are skipped because Model.Partition() itself is a no-op on every
// driver: the partition names are stored on the model and never read, so they
// never reach the generated SQL. See https://github.com/gogf/gf/issues/4826
//
// The function names are kept aligned with the MySQL/MariaDB baseline so the
// cases can be filled in once Model.Partition() is wired up.

func Test_Partition_Range_Insert_And_Query(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}

func Test_Partition_Range_PartitionQuery(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}

func Test_Partition_Hash_Insert_And_Distribution(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}

func Test_Partition_List_Insert_And_Query(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}

func Test_Partition_Range_Update(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}

func Test_Partition_Range_Delete(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}

func Test_Partition_Transaction(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}

func Test_Partition_Range_Count_And_Sum(t *testing.T) {
	t.Skip("Model.Partition() is a no-op on all drivers; see https://github.com/gogf/gf/issues/4826")
}
