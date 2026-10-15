// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

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

const (
	masterSlaveUser          = "gf_ms_slave"
	masterSlaveErrUserExists = 1920
)

var (
	masterSlaveUserOnce sync.Once
	masterSlaveUserErr  error
)

func masterSlavePrepareUser(t *gtest.T) {
	masterSlaveUserOnce.Do(func() {
		masterSlaveUserErr = masterSlaveCreateUser()
	})
	t.AssertNil(masterSlaveUserErr)
}

func masterSlaveCreateUser() error {
	_, err := db.Exec(ctx, fmt.Sprintf("CREATE USER %s IDENTIFIED BY %s", masterSlaveUser, TestDbPass))
	if err != nil && !isOracleError(err, masterSlaveErrUserExists) {
		return err
	}
	_, err = db.Exec(ctx, fmt.Sprintf("GRANT CONNECT, RESOURCE TO %s", masterSlaveUser))
	return err
}

func masterSlaveConfigGroup() gdb.ConfigGroup {
	return gdb.ConfigGroup{
		gdb.ConfigNode{
			Host:   TestDbIP,
			Port:   TestDbPort,
			User:   TestDbUser,
			Pass:   TestDbPass,
			Name:   TestDbName,
			Type:   TestDbType,
			Role:   "master",
			Weight: 100,
		},
		gdb.ConfigNode{
			Host:   TestDbIP,
			Port:   TestDbPort,
			User:   masterSlaveUser,
			Pass:   TestDbPass,
			Name:   TestDbName,
			Type:   TestDbType,
			Role:   "slave",
			Weight: 100,
		},
	}
}

func masterSlaveTableName() string {
	return fmt.Sprintf("t_ms_%d", gtime.TimestampMicro()%1e9)
}

func Test_Master_Slave(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		masterSlavePrepareUser(t)
	})

	var (
		configKey   = guid.S()
		configGroup = masterSlaveConfigGroup()
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)
	masterSlaveDB.SetDebug(true)
	defer masterSlaveDB.Close(ctx)
	gtest.C(t, func(t *gtest.T) {
		slaveDB, err := gdb.New(configGroup[1])
		t.AssertNil(err)
		defer slaveDB.Close(ctx)

		table := masterSlaveTableName()
		createTableWithDb(db, table)
		createTableWithDb(slaveDB, table)
		defer dropTableWithDb(db, table)
		defer dropTableWithDb(slaveDB, table)

		array := garray.New(true)
		for i := 1; i <= TableSize; i++ {
			array.Append(g.Map{
				"id":          i,
				"passport":    fmt.Sprintf(`user_%d`, i),
				"password":    fmt.Sprintf(`pass_%d`, i),
				"nickname":    fmt.Sprintf(`name_%d`, i),
				"create_time": gtime.Now().String(),
			})
		}
		_, err = masterSlaveDB.Model(table).Data(array).Insert()
		t.AssertNil(err)

		var count int
		count, err = masterSlaveDB.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = masterSlaveDB.Model(table).Slave().Count()
		t.AssertNil(err)
		t.Assert(count, int64(0))

		count, err = masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, int64(TableSize))
	})
}

// Test_Master_Slave_Concurrent_ReadWrite tests concurrent read/write routing
func Test_Master_Slave_Concurrent_ReadWrite(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		masterSlavePrepareUser(t)
	})

	var (
		configKey   = guid.S()
		configGroup = masterSlaveConfigGroup()
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)
	defer masterSlaveDB.Close(ctx)

	gtest.C(t, func(t *gtest.T) {
		slaveDB, err := gdb.New(configGroup[1])
		t.AssertNil(err)
		defer slaveDB.Close(ctx)

		table := masterSlaveTableName()
		createTableWithDb(db, table)
		createTableWithDb(slaveDB, table)
		defer dropTableWithDb(db, table)
		defer dropTableWithDb(slaveDB, table)

		var (
			wg          sync.WaitGroup
			concurrency = 10
			errs        = make([]error, concurrency)
		)
		wg.Add(concurrency)
		for i := 0; i < concurrency; i++ {
			go func(id int) {
				defer wg.Done()
				_, errs[id] = masterSlaveDB.Model(table).Insert(g.Map{
					"passport": fmt.Sprintf("concurrent_%d", id),
					"password": fmt.Sprintf("pass_%d", id),
					"nickname": fmt.Sprintf("name_%d", id),
				})
			}(i)
		}
		wg.Wait()
		for _, err := range errs {
			t.AssertNil(err)
		}

		count, err := masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, concurrency)

		count, err = masterSlaveDB.Model(table).Slave().Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Master_Slave_Transaction_Routing tests transaction routing to master
func Test_Master_Slave_Transaction_Routing(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		masterSlavePrepareUser(t)
	})

	var (
		configKey   = guid.S()
		configGroup = masterSlaveConfigGroup()
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)
	defer masterSlaveDB.Close(ctx)

	gtest.C(t, func(t *gtest.T) {
		slaveDB, err := gdb.New(configGroup[1])
		t.AssertNil(err)
		defer slaveDB.Close(ctx)

		table := masterSlaveTableName()
		createTableWithDb(db, table)
		createTableWithDb(slaveDB, table)
		defer dropTableWithDb(db, table)
		defer dropTableWithDb(slaveDB, table)

		err = masterSlaveDB.Transaction(ctx, func(ctx context.Context, tx gdb.TX) error {
			_, err := tx.Model(table).Insert(g.Map{
				"passport": "tx_user",
				"password": "tx_pass",
				"nickname": "tx_name",
			})
			if err != nil {
				return err
			}

			count, err := tx.Model(table).Count()
			t.AssertNil(err)
			t.Assert(count, 1)

			return nil
		})
		t.AssertNil(err)

		count, err := masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		count, err = masterSlaveDB.Model(table).Slave().Count()
		t.AssertNil(err)
		t.Assert(count, 0)
	})
}

// Test_Master_Slave_Explicit_Selection tests explicit master/slave selection
func Test_Master_Slave_Explicit_Selection(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		masterSlavePrepareUser(t)
	})

	var (
		configKey   = guid.S()
		configGroup = masterSlaveConfigGroup()
	)
	gdb.SetConfigGroup(configKey, configGroup)
	masterSlaveDB := g.DB(configKey)
	defer masterSlaveDB.Close(ctx)

	gtest.C(t, func(t *gtest.T) {
		slaveDB, err := gdb.New(configGroup[1])
		t.AssertNil(err)
		defer slaveDB.Close(ctx)

		table := masterSlaveTableName()
		createTableWithDb(db, table)
		createTableWithDb(slaveDB, table)
		defer dropTableWithDb(db, table)
		defer dropTableWithDb(slaveDB, table)

		_, err = masterSlaveDB.Model(table).Master().Insert(g.Map{
			"passport": "explicit_test",
			"password": "pass",
			"nickname": "name",
		})
		t.AssertNil(err)

		count, err := masterSlaveDB.Model(table).Slave().Count()
		t.AssertNil(err)
		t.Assert(count, 0)

		count, err = masterSlaveDB.Model(table).Master().Count()
		t.AssertNil(err)
		t.Assert(count, 1)
	})
}
