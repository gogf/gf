// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mariadb_test

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

func Test_Master_Slave(t *testing.T) {
	var err error

	gtest.C(t, func(t *gtest.T) {
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `master` CHARACTER SET UTF8")
		t.AssertNil(err)
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `slave` CHARACTER SET UTF8")
		t.AssertNil(err)
	})
	defer func() {
		_, _ = db.Exec(ctx, "DROP DATABASE `master`")
		_, _ = db.Exec(ctx, "DROP DATABASE `slave`")
	}()
	var (
		configKey   = guid.S()
		configGroup = gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "master",
				Type:   "mariadb",
				Role:   "master",
				Debug:  true,
				Weight: 100,
			},
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "slave",
				Type:   "mariadb",
				Role:   "slave",
				Debug:  true,
				Weight: 100,
			},
		}
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)
	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema("master"), table)
		createTableWithDb(masterSlaveDB.Schema("slave"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("master"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("slave"), table)

		// Data insert to master.
		array := garray.New(true)
		for i := 1; i <= TableSize; i++ {
			array.Append(g.Map{
				"id":          i,
				"passport":    fmt.Sprintf(`user_%d`, i),
				"password":    fmt.Sprintf(`pass_%d`, i),
				"nickname":    fmt.Sprintf(`name_%d`, i),
				"create_time": gtime.NewFromStr(CreateTime).String(),
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
	var err error

	gtest.C(t, func(t *gtest.T) {
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `master` CHARACTER SET UTF8")
		t.AssertNil(err)
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `slave` CHARACTER SET UTF8")
		t.AssertNil(err)
	})
	defer func() {
		_, _ = db.Exec(ctx, "DROP DATABASE `master`")
		_, _ = db.Exec(ctx, "DROP DATABASE `slave`")
	}()

	var (
		configKey   = guid.S()
		configGroup = gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "master",
				Type:   "mariadb",
				Role:   "master",
				Weight: 100,
			},
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "slave",
				Type:   "mariadb",
				Role:   "slave",
				Weight: 100,
			},
		}
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)

	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema("master"), table)
		createTableWithDb(masterSlaveDB.Schema("slave"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("master"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("slave"), table)

		var wg sync.WaitGroup
		concurrency := 10

		// Concurrent writes to master
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				_, err := masterSlaveDB.Model(table).Insert(g.Map{
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
	})
}

// Test_Master_Slave_Transaction_Routing tests transaction routing to master
func Test_Master_Slave_Transaction_Routing(t *testing.T) {
	var err error

	gtest.C(t, func(t *gtest.T) {
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `master` CHARACTER SET UTF8")
		t.AssertNil(err)
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `slave` CHARACTER SET UTF8")
		t.AssertNil(err)
	})
	defer func() {
		_, _ = db.Exec(ctx, "DROP DATABASE `master`")
		_, _ = db.Exec(ctx, "DROP DATABASE `slave`")
	}()

	var (
		configKey   = guid.S()
		configGroup = gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "master",
				Type:   "mariadb",
				Role:   "master",
				Weight: 100,
			},
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "slave",
				Type:   "mariadb",
				Role:   "slave",
				Weight: 100,
			},
		}
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)

	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema("master"), table)
		createTableWithDb(masterSlaveDB.Schema("slave"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("master"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("slave"), table)

		// Transaction should route to master
		err := masterSlaveDB.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Insert(g.Map{
				"passport": "tx_user",
				"password": "tx_pass",
				"nickname": "tx_name",
			})
			if err != nil {
				return err
			}

			// Read within transaction should also use master
			count, err := tx.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count, 1)

			return nil
		})
		t.AssertNil(err)

		// Verify data is in master
		count, err := masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}

// Test_Master_Slave_Explicit_Selection tests explicit master/slave selection
func Test_Master_Slave_Explicit_Selection(t *testing.T) {
	var err error

	gtest.C(t, func(t *gtest.T) {
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `master` CHARACTER SET UTF8")
		t.AssertNil(err)
		_, err = db.Exec(ctx, "CREATE DATABASE IF NOT EXISTS `slave` CHARACTER SET UTF8")
		t.AssertNil(err)
	})
	defer func() {
		_, _ = db.Exec(ctx, "DROP DATABASE `master`")
		_, _ = db.Exec(ctx, "DROP DATABASE `slave`")
	}()

	var (
		configKey   = guid.S()
		configGroup = gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "master",
				Type:   "mariadb",
				Role:   "master",
				Weight: 100,
			},
			gdb.ConfigNode{
				Host:   "127.0.0.1",
				Port:   "3307",
				User:   "root",
				Pass:   "12345678",
				Name:   "slave",
				Type:   "mariadb",
				Role:   "slave",
				Weight: 100,
			},
		}
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)

	gtest.C(t, func(t *gtest.T) {
		table := "table_" + guid.S()
		createTableWithDb(masterSlaveDB.Schema("master"), table)
		createTableWithDb(masterSlaveDB.Schema("slave"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("master"), table)
		defer dropTableWithDb(masterSlaveDB.Schema("slave"), table)

		// Insert to master
		_, err := masterSlaveDB.Model(table).Master().Insert(g.Map{
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
