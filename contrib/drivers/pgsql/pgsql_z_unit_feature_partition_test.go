// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package pgsql_test

import (
	"testing"
)

// PostgreSQL supports declarative partitioning with a different syntax than
// MySQL/MariaDB (PARTITION BY RANGE/LIST/HASH at CREATE TABLE level, plus
// separate partition child tables). The Model.Partition("<name>") API used by
// the MySQL baseline tests targets MySQL partition-name syntax
// (`... PARTITION (p0)`), which PostgreSQL does not accept. Until gdb provides
// dialect-specific partition routing for PostgreSQL, the following tests are
// skipped to preserve function-name parity with the MySQL/MariaDB baseline.

func Test_Partition_Range_Insert_And_Query(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}

func Test_Partition_Range_PartitionQuery(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}

func Test_Partition_Hash_Insert_And_Distribution(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}

func Test_Partition_List_Insert_And_Query(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}

func Test_Partition_Range_Update(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}

func Test_Partition_Range_Delete(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}

func Test_Partition_Transaction(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}

func Test_Partition_Range_Count_And_Sum(t *testing.T) {
	t.Skip("PostgreSQL partition syntax differs from MySQL; Model.Partition() not supported on pgsql")
}
