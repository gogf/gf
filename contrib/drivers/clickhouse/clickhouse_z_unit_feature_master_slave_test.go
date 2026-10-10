// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package clickhouse_test

import (
	"context"
	"fmt"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
)

// chDMasterSlaveNode returns a configuration node of the local server using database `name` in `role`.
func chDMasterSlaveNode(name, role string) gdb.ConfigNode {
	return gdb.ConfigNode{
		Host:   "127.0.0.1",
		Port:   "9000",
		User:   "default",
		Name:   name,
		Type:   "clickhouse",
		Role:   gdb.Role(role),
		Weight: 100,
		Extra:  "mutations_sync=1",
	}
}

// chDNewMasterSlaveDb creates the databases of a master node and a slave node on the local server and
// returns the db of a configuration group of both nodes with the names of both databases.
func chDNewMasterSlaveDb() (masterSlaveDB gdb.DB, master, slave string) {
	master = chDCreateDatabase("master")
	slave = chDCreateDatabase("slave")
	configKey := guid.S()
	if err := gdb.SetConfigGroup(configKey, gdb.ConfigGroup{
		chDMasterSlaveNode(master, "master"),
		chDMasterSlaveNode(slave, "slave"),
	}); err != nil {
		gtest.Fatal(err)
	}
	masterSlaveDB = g.DB(configKey)
	return
}

// Test_Master_Slave tests that writes go to the master node and reads go to the slave node by default.
func Test_Master_Slave(t *testing.T) {
	var err error

	masterSlaveDB, master, slave := chDNewMasterSlaveDb()
	defer chDDropDatabase(master, slave)
	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema(master), table)
		createTableWithDb(masterSlaveDB.Schema(slave), table)
		defer dropTableWithDb(masterSlaveDB.Schema(master), table)
		defer dropTableWithDb(masterSlaveDB.Schema(slave), table)

		// Data insert to master.
		array := garray.New(true)
		for i := 1; i <= TableSize; i++ {
			array.Append(g.Map{
				"id":          i,
				"passport":    fmt.Sprintf(`user_%d`, i),
				"password":    fmt.Sprintf(`pass_%d`, i),
				"nickname":    fmt.Sprintf(`name_%d`, i),
				"create_time": gtime.NewFromStr(chDCreateTime).String(),
			})
		}
		_, err = masterSlaveDB.Model(table).Data(array).Insert()
		t.AssertNil(err)

		var count int
		// Auto slave.
		count, err = masterSlaveDB.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// slave.
		count, err = masterSlaveDB.Model(table).Slave().Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		// master.
		count, err = masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

// Test_Master_Slave_Concurrent_ReadWrite tests concurrent read/write routing
func Test_Master_Slave_Concurrent_ReadWrite(t *testing.T) {
	masterSlaveDB, master, slave := chDNewMasterSlaveDb()
	defer chDDropDatabase(master, slave)

	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema(master), table)
		createTableWithDb(masterSlaveDB.Schema(slave), table)
		defer dropTableWithDb(masterSlaveDB.Schema(master), table)
		defer dropTableWithDb(masterSlaveDB.Schema(slave), table)

		var wg sync.WaitGroup
		concurrency := 10

		// Concurrent writes to master
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				// Note: ClickHouse has no auto increment, so the id is given.
				_, err := masterSlaveDB.Model(table).Insert(g.Map{
					"id":       id,
					"passport": fmt.Sprintf("concurrent_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
				t.AssertNil(err)
			}(i)
		}
		wg.Wait()

		// Verify writes went to master
		count, err := masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, concurrency)

		count, err = masterSlaveDB.Model(table).Slave().Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Master_Slave_Transaction_Routing tests that a transaction on a master-slave db fails as unsupported.
func Test_Master_Slave_Transaction_Routing(t *testing.T) {
	masterSlaveDB, master, slave := chDNewMasterSlaveDb()
	defer chDDropDatabase(master, slave)

	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema(master), table)
		createTableWithDb(masterSlaveDB.Schema(slave), table)
		defer dropTableWithDb(masterSlaveDB.Schema(master), table)
		defer dropTableWithDb(masterSlaveDB.Schema(slave), table)

		// Note: ClickHouse has no transactions, so the driver rejects Transaction before calling the function.
		var called bool
		err := masterSlaveDB.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			called = true
			_, err := tx.Model(table).Insert(g.Map{
				"id":       1,
				"passport": "tx_user",
				"password": "tx_pass",
				"nickname": "tx_name",
			})
			return err
		})
		t.AssertNE(err, nil)
		t.Assert(err.Error(), "unsupported method:Transaction")
		t.Assert(called, false)

		count, err := masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Master_Slave_Explicit_Selection tests explicit master/slave selection
func Test_Master_Slave_Explicit_Selection(t *testing.T) {
	masterSlaveDB, master, slave := chDNewMasterSlaveDb()
	defer chDDropDatabase(master, slave)

	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema(master), table)
		createTableWithDb(masterSlaveDB.Schema(slave), table)
		defer dropTableWithDb(masterSlaveDB.Schema(master), table)
		defer dropTableWithDb(masterSlaveDB.Schema(slave), table)

		// Insert to master
		_, err := masterSlaveDB.Model(table).Master().Insert(g.Map{
			"id":       1,
			"passport": "explicit_test",
			"password": "pass",
			"nickname": "name",
		})
		t.AssertNil(err)

		// Explicitly read from slave (should be empty)
		count, err := masterSlaveDB.Model(table).Slave().Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		// Explicitly read from master (should have data)
		count, err = masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}
