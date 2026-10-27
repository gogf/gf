// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package tidb_test

import (
	"fmt"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gstr"
)

// autoRandomShardBits is the default AUTO_RANDOM shard bit width, so the shard
// prefix of a generated id is id>>(63-5).
const autoRandomShardBits = 5

// tidbFeatureConn opens a dedicated connection whose pool holds exactly one
// physical connection, so session variables set through it stay in effect for
// every following statement.
func tidbFeatureConn(extra string) (gdb.DB, error) {
	node := gdb.ConfigNode{
		Host:             "127.0.0.1",
		Port:             "4000",
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             TestSchema1,
		Type:             "tidb",
		Extra:            extra,
		MaxIdleConnCount: 1,
		MaxOpenConnCount: 1,
		TranTimeout:      time.Second * 5,
	}
	return gdb.New(node)
}

// tidbFeatureTable creates a uniquely named table from the given column
// definition and options.
func tidbFeatureTable(t *gtest.T, prefix, definition string) string {
	name := fmt.Sprintf(`%s_%d`, prefix, gtime.TimestampNano())
	_, err := db.Exec(ctx, fmt.Sprintf(`CREATE TABLE %s %s`, name, definition))
	t.AssertNil(err)
	return name
}

// =============================================================================
// AUTO_RANDOM and row id allocation
// =============================================================================

// Test_TiDB_AutoRandom_NonSequentialIds asserts AUTO_RANDOM ids are scattered
// across shards instead of counting up from 1.
func Test_TiDB_AutoRandom_NonSequentialIds(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "auto_random", `(
			id   bigint PRIMARY KEY AUTO_RANDOM,
			name varchar(50)
		)`)
		defer dropTable(table)

		var (
			ids    []int64
			shards = map[int64]struct{}{}
		)
		for i := 0; i < 6; i++ {
			conn, err := tidbFeatureConn("")
			t.AssertNil(err)
			id, err := conn.Model(table).Data(g.Map{"name": fmt.Sprintf("n_%d", i)}).InsertAndGetId()
			t.AssertNil(err)
			t.AssertNil(conn.Close(ctx))

			ids = append(ids, id)
			shards[id>>(63-autoRandomShardBits)] = struct{}{}
		}

		// Each transaction picks its own shard, so the ids are not monotonic.
		// A shard of 0 is a legal draw, hence this checks the set, not each id.
		t.AssertGT(len(shards), 1)

		// At least one draw lands on a non-zero shard, which puts the value far
		// above what a plain counter would have produced.
		var maxID int64
		for _, id := range ids {
			if id > maxID {
				maxID = id
			}
		}
		t.AssertGT(maxID, int64(1)<<58)

		// Ids are unique and every row is readable back by its id.
		seen := map[int64]struct{}{}
		for _, id := range ids {
			_, dup := seen[id]
			t.Assert(dup, false)
			seen[id] = struct{}{}

			one, err := db.Model(table).Where("id", id).One()
			t.AssertNil(err)
			t.Assert(one.IsEmpty(), false)
		}

		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 6)
	})
}

// Test_TiDB_AutoRandom_SameTransactionShard asserts AUTO_RANDOM keeps the shard
// fixed inside one transaction and increments only the low bits.
func Test_TiDB_AutoRandom_SameTransactionShard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "auto_random_one", `(
			id   bigint PRIMARY KEY AUTO_RANDOM,
			name varchar(50)
		)`)
		defer dropTable(table)

		// A single multi row insert is one transaction.
		_, err := db.Model(table).Data(g.List{
			{"name": "n_0"},
			{"name": "n_1"},
			{"name": "n_2"},
			{"name": "n_3"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf("SELECT id FROM %s ORDER BY id", table))
		t.AssertNil(err)
		t.Assert(len(all), 4)

		ids := make([]int64, 0, 4)
		for _, row := range all {
			ids = append(ids, row["id"].Int64())
		}

		// One transaction keeps one shard prefix.
		shard := ids[0] >> (63 - autoRandomShardBits)
		for _, id := range ids {
			t.Assert(id>>(63-autoRandomShardBits), shard)
		}

		// The low bits are a plain counter.
		t.Assert(ids[1]-ids[0], int64(1))
		t.Assert(ids[2]-ids[1], int64(1))
		t.Assert(ids[3]-ids[2], int64(1))
	})
}

// Test_TiDB_AutoRandom_DDLRoundTrip asserts the AUTO_RANDOM attribute survives
// in the table definition read back through the ORM.
func Test_TiDB_AutoRandom_DDLRoundTrip(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "auto_random_ddl", `(
			id   bigint PRIMARY KEY AUTO_RANDOM,
			name varchar(50)
		)`)
		defer dropTable(table)

		one, err := db.GetOne(ctx, fmt.Sprintf("SHOW CREATE TABLE %s", table))
		t.AssertNil(err)
		t.Assert(gstr.Contains(one["Create Table"].String(), "AUTO_RANDOM"), true)

		// The id column is not reported as auto_increment.
		fields, err := db.TableFields(ctx, table)
		t.AssertNil(err)
		t.Assert(fields["id"].Extra, "")
	})
}

// Test_TiDB_ShardRowIdBits_Scatters asserts SHARD_ROW_ID_BITS scatters
// _tidb_rowid while a plain table allocates it sequentially.
func Test_TiDB_ShardRowIdBits_Scatters(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		sharded := tidbFeatureTable(t, "shard_bits", `(
			id int,
			v  varchar(10)
		) SHARD_ROW_ID_BITS=4`)
		defer dropTable(sharded)

		plain := tidbFeatureTable(t, "shard_none", `(
			id int,
			v  varchar(10)
		)`)
		defer dropTable(plain)

		for i := 1; i <= 5; i++ {
			_, err := db.Model(sharded).Data(g.Map{"id": i, "v": "x"}).Insert()
			t.AssertNil(err)
			_, err = db.Model(plain).Data(g.Map{"id": i, "v": "x"}).Insert()
			t.AssertNil(err)
		}

		// The plain table numbers rows 1..5.
		all, err := db.GetAll(ctx, fmt.Sprintf("SELECT _tidb_rowid, id FROM %s ORDER BY id", plain))
		t.AssertNil(err)
		t.Assert(len(all), 5)
		for i, row := range all {
			t.Assert(row["_tidb_rowid"].Int64(), int64(i+1))
		}

		// The sharded table spreads rowids over several shards.
		all, err = db.GetAll(ctx, fmt.Sprintf("SELECT _tidb_rowid, id FROM %s ORDER BY id", sharded))
		t.AssertNil(err)
		t.Assert(len(all), 5)
		var (
			shards   = map[int64]struct{}{}
			maxRowID int64
		)
		for _, row := range all {
			rowID := row["_tidb_rowid"].Int64()
			shards[rowID>>59] = struct{}{}
			if rowID > maxRowID {
				maxRowID = rowID
			}
		}
		// A shard of 0 is a legal draw, so this checks the set rather than each row.
		t.AssertGT(len(shards), 1)
		t.AssertGT(maxRowID, int64(1)<<58)

		// The DDL keeps the option.
		one, err := db.GetOne(ctx, fmt.Sprintf("SHOW CREATE TABLE %s", sharded))
		t.AssertNil(err)
		t.Assert(gstr.Contains(one["Create Table"].String(), "SHARD_ROW_ID_BITS=4"), true)
	})
}

// Test_TiDB_AutoIdCache_Contiguous asserts AUTO_ID_CACHE 1 keeps auto increment
// ids contiguous across reconnects.
func Test_TiDB_AutoIdCache_Contiguous(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "auto_id_cache", `(
			id bigint PRIMARY KEY AUTO_INCREMENT,
			v  varchar(10)
		) AUTO_ID_CACHE 1`)
		defer dropTable(table)

		for i := 1; i <= 5; i++ {
			conn, err := tidbFeatureConn("")
			t.AssertNil(err)
			id, err := conn.Model(table).Data(g.Map{"v": "x"}).InsertAndGetId()
			t.AssertNil(err)
			t.AssertNil(conn.Close(ctx))
			t.Assert(id, int64(i))
		}

		one, err := db.GetOne(ctx, fmt.Sprintf("SHOW CREATE TABLE %s", table))
		t.AssertNil(err)
		t.Assert(gstr.Contains(one["Create Table"].String(), "AUTO_ID_CACHE=1"), true)
	})
}

// =============================================================================
// Clustered and nonclustered primary keys
// =============================================================================

// Test_TiDB_ClusteredIndex_Clustered asserts a CLUSTERED primary key table is
// fully usable through the ORM and exposes no _tidb_rowid.
func Test_TiDB_ClusteredIndex_Clustered(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "clustered", `(
			id int PRIMARY KEY CLUSTERED,
			v  varchar(20)
		)`)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"id": 1, "v": "a"},
			{"id": 2, "v": "b"},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{"v": "a2"}).Where("id", 1).Update()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["v"].String(), "a2")

		_, err = db.Model(table).Where("id", 2).Delete()
		t.AssertNil(err)
		count, err := db.Model(table).Count()
		t.AssertNil(err)
		t.Assert(count, 1)

		clustered, err := db.GetValue(ctx, fmt.Sprintf(
			"SELECT CLUSTERED FROM information_schema.tidb_indexes WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'",
			TestSchema1, table,
		))
		t.AssertNil(err)
		t.Assert(clustered.String(), "YES")

		// The primary key is the row handle, so there is no separate rowid column.
		_, err = db.GetAll(ctx, fmt.Sprintf("SELECT _tidb_rowid FROM %s", table))
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "_tidb_rowid"), true)
	})
}

// Test_TiDB_ClusteredIndex_NonClustered asserts a NONCLUSTERED primary key
// table works through the ORM and exposes a readable _tidb_rowid.
func Test_TiDB_ClusteredIndex_NonClustered(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "nonclustered", `(
			id int PRIMARY KEY NONCLUSTERED,
			v  varchar(20)
		)`)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"id": 10, "v": "a"},
			{"id": 20, "v": "b"},
		}).Insert()
		t.AssertNil(err)

		_, err = db.Model(table).Data(g.Map{"v": "b2"}).Where("id", 20).Update()
		t.AssertNil(err)
		one, err := db.Model(table).Where("id", 20).One()
		t.AssertNil(err)
		t.Assert(one["v"].String(), "b2")

		clustered, err := db.GetValue(ctx, fmt.Sprintf(
			"SELECT CLUSTERED FROM information_schema.tidb_indexes WHERE TABLE_SCHEMA='%s' AND TABLE_NAME='%s'",
			TestSchema1, table,
		))
		t.AssertNil(err)
		t.Assert(clustered.String(), "NO")

		// The hidden handle column is allocated sequentially and is readable.
		all, err := db.GetAll(ctx, fmt.Sprintf("SELECT _tidb_rowid, id FROM %s ORDER BY id", table))
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["_tidb_rowid"].Int64(), int64(1))
		t.Assert(all[0]["id"].Int(), 10)
		t.Assert(all[1]["_tidb_rowid"].Int64(), int64(2))
		t.Assert(all[1]["id"].Int(), 20)
	})
}

// =============================================================================
// Stale read
// =============================================================================

// Test_TiDB_StaleRead_AsOfTimestamp asserts SELECT ... AS OF TIMESTAMP returns
// the value from before a later update.
func Test_TiDB_StaleRead_AsOfTimestamp(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"id": 1, "nickname": "old"}).Insert()
		t.AssertNil(err)

		time.Sleep(2 * time.Second)

		_, err = db.Model(table).Data(g.Map{"nickname": "new"}).Where("id", 1).Update()
		t.AssertNil(err)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "new")

		all, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT nickname FROM %s AS OF TIMESTAMP NOW() - INTERVAL 1 SECOND WHERE id=1", table,
		))
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["nickname"].String(), "old")
	})
}

// Test_TiDB_StaleRead_StartTransactionAsOf asserts a read only transaction
// opened at a past timestamp sees the old value.
func Test_TiDB_StaleRead_StartTransactionAsOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"id": 1, "nickname": "old"}).Insert()
		t.AssertNil(err)

		time.Sleep(2 * time.Second)

		_, err = db.Model(table).Data(g.Map{"nickname": "new"}).Where("id", 1).Update()
		t.AssertNil(err)

		conn, err := tidbFeatureConn("")
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "START TRANSACTION READ ONLY AS OF TIMESTAMP NOW() - INTERVAL 1 SECOND")
		t.AssertNil(err)
		value, err := conn.GetValue(ctx, fmt.Sprintf("SELECT nickname FROM %s WHERE id=1", table))
		t.AssertNil(err)
		t.Assert(value.String(), "old")
		_, err = conn.Exec(ctx, "COMMIT")
		t.AssertNil(err)

		// After the stale transaction ends the session is back on current data.
		value, err = conn.GetValue(ctx, fmt.Sprintf("SELECT nickname FROM %s WHERE id=1", table))
		t.AssertNil(err)
		t.Assert(value.String(), "new")
	})
}

// Test_TiDB_StaleRead_SetTransactionAsOf asserts SET TRANSACTION READ ONLY AS OF
// TIMESTAMP applies to the next transaction only.
func Test_TiDB_StaleRead_SetTransactionAsOf(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createTable()
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{"id": 1, "nickname": "old"}).Insert()
		t.AssertNil(err)

		time.Sleep(2 * time.Second)

		_, err = db.Model(table).Data(g.Map{"nickname": "new"}).Where("id", 1).Update()
		t.AssertNil(err)

		conn, err := tidbFeatureConn("")
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "SET TRANSACTION READ ONLY AS OF TIMESTAMP NOW() - INTERVAL 1 SECOND")
		t.AssertNil(err)
		_, err = conn.Exec(ctx, "BEGIN")
		t.AssertNil(err)
		value, err := conn.GetValue(ctx, fmt.Sprintf("SELECT nickname FROM %s WHERE id=1", table))
		t.AssertNil(err)
		t.Assert(value.String(), "old")
		_, err = conn.Exec(ctx, "COMMIT")
		t.AssertNil(err)

		// The next transaction is a normal one again.
		_, err = conn.Exec(ctx, "BEGIN")
		t.AssertNil(err)
		value, err = conn.GetValue(ctx, fmt.Sprintf("SELECT nickname FROM %s WHERE id=1", table))
		t.AssertNil(err)
		t.Assert(value.String(), "new")
		_, err = conn.Exec(ctx, "COMMIT")
		t.AssertNil(err)
	})
}

// =============================================================================
// Transaction modes
// =============================================================================

// Test_TiDB_TxnMode_SessionVariable asserts tidb_txn_mode defaults to
// pessimistic and can be switched per connection.
func Test_TiDB_TxnMode_SessionVariable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		def, err := tidbFeatureConn("")
		t.AssertNil(err)
		defer def.Close(ctx)
		mode, err := def.GetValue(ctx, "SELECT @@tidb_txn_mode")
		t.AssertNil(err)
		t.Assert(mode.String(), "pessimistic")

		opt, err := tidbFeatureConn("tidb_txn_mode=optimistic")
		t.AssertNil(err)
		defer opt.Close(ctx)
		mode, err = opt.GetValue(ctx, "SELECT @@tidb_txn_mode")
		t.AssertNil(err)
		t.Assert(mode.String(), "optimistic")
	})
}

// Test_TiDB_TxnMode_OptimisticWriteConflict asserts the loser of two concurrent
// optimistic transactions fails at commit with a write conflict.
func Test_TiDB_TxnMode_OptimisticWriteConflict(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		conn1, err := tidbFeatureConn("tidb_txn_mode=optimistic")
		t.AssertNil(err)
		defer conn1.Close(ctx)
		conn2, err := tidbFeatureConn("tidb_txn_mode=optimistic")
		t.AssertNil(err)
		defer conn2.Close(ctx)

		tx1, err := conn1.Begin(ctx)
		t.AssertNil(err)
		tx2, err := conn2.Begin(ctx)
		t.AssertNil(err)

		// Optimistic transactions do not lock, so both updates buffer locally.
		_, err = tx1.Model(table).Data(g.Map{"nickname": "tx1"}).Where("id", 1).Update()
		t.AssertNil(err)
		_, err = tx2.Model(table).Data(g.Map{"nickname": "tx2"}).Where("id", 1).Update()
		t.AssertNil(err)

		// The first commit wins.
		t.AssertNil(tx1.Commit())

		// The second one detects the conflict only at commit time.
		err = tx2.Commit()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "Write conflict"), true)
		t.Assert(gstr.Contains(err.Error(), "Error 9007"), true)

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "tx1")
	})
}

// Test_TiDB_TxnMode_PessimisticBlocks asserts a pessimistic transaction blocks
// on a locked row and succeeds once the lock holder commits.
func Test_TiDB_TxnMode_PessimisticBlocks(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		conn1, err := tidbFeatureConn("tidb_txn_mode=pessimistic")
		t.AssertNil(err)
		defer conn1.Close(ctx)
		conn2, err := tidbFeatureConn("tidb_txn_mode=pessimistic")
		t.AssertNil(err)
		defer conn2.Close(ctx)

		tx1, err := conn1.Begin(ctx)
		t.AssertNil(err)
		_, err = tx1.Model(table).Data(g.Map{"nickname": "tx1"}).Where("id", 1).Update()
		t.AssertNil(err)

		done := make(chan error, 1)
		go func() {
			tx2, e := conn2.Begin(ctx)
			if e != nil {
				done <- e
				return
			}
			if _, e = tx2.Model(table).Data(g.Map{"nickname": "tx2"}).Where("id", 1).Update(); e != nil {
				done <- e
				return
			}
			done <- tx2.Commit()
		}()

		// The row is locked, so the second writer cannot get through.
		time.Sleep(time.Second)
		select {
		case e := <-done:
			t.Errorf("pessimistic transaction was not blocked, finished with: %v", e)
		default:
		}

		t.AssertNil(tx1.Commit())

		// Releasing the lock lets the waiter finish.
		select {
		case e := <-done:
			t.AssertNil(e)
		case <-time.After(5 * time.Second):
			t.Error("blocked transaction did not finish after the lock was released")
		}

		one, err := db.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "tx2")
	})
}

// =============================================================================
// Locking clauses
// =============================================================================

// Test_TiDB_ForUpdate_Nowait asserts FOR UPDATE NOWAIT fails immediately on a
// row locked by another transaction.
func Test_TiDB_ForUpdate_Nowait(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		conn1, err := tidbFeatureConn("tidb_txn_mode=pessimistic")
		t.AssertNil(err)
		defer conn1.Close(ctx)
		conn2, err := tidbFeatureConn("tidb_txn_mode=pessimistic")
		t.AssertNil(err)
		defer conn2.Close(ctx)

		tx1, err := conn1.Begin(ctx)
		t.AssertNil(err)
		_, err = tx1.Model(table).Data(g.Map{"nickname": "locked"}).Where("id", 1).Update()
		t.AssertNil(err)

		tx2, err := conn2.Begin(ctx)
		t.AssertNil(err)

		start := gtime.TimestampMilli()
		_, err = tx2.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=1 FOR UPDATE NOWAIT", table))
		elapsed := gtime.TimestampMilli() - start

		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "Error 3572"), true)
		t.Assert(gstr.Contains(err.Error(), "NOWAIT"), true)
		// It must fail immediately rather than after the lock wait timeout.
		t.AssertLT(elapsed, int64(2000))

		t.AssertNil(tx2.Rollback())

		// An unlocked row is still lockable with NOWAIT.
		tx3, err := conn2.Begin(ctx)
		t.AssertNil(err)
		rows, err := tx3.GetAll(fmt.Sprintf("SELECT * FROM %s WHERE id=2 FOR UPDATE NOWAIT", table))
		t.AssertNil(err)
		t.Assert(len(rows), 1)
		t.AssertNil(tx3.Rollback())

		t.AssertNil(tx1.Commit())
	})
}

// Test_TiDB_NoopFunctions_LockShared asserts LOCK IN SHARE MODE is rejected
// until tidb_enable_noop_functions is turned on for the session.
func Test_TiDB_NoopFunctions_LockShared(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		conn, err := tidbFeatureConn("")
		t.AssertNil(err)
		defer conn.Close(ctx)

		value, err := conn.GetValue(ctx, "SELECT @@tidb_enable_noop_functions")
		t.AssertNil(err)
		t.Assert(value.String(), "OFF")

		_, err = conn.Model(table).LockShared().Where("id", 1).One()
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "Error 1235"), true)
		t.Assert(gstr.Contains(err.Error(), "LOCK IN SHARE MODE"), true)

		_, err = conn.Exec(ctx, "SET SESSION tidb_enable_noop_functions=1")
		t.AssertNil(err)
		value, err = conn.GetValue(ctx, "SELECT @@tidb_enable_noop_functions")
		t.AssertNil(err)
		t.Assert(value.String(), "ON")

		// The clause is now accepted and the query returns the row.
		one, err := conn.Model(table).LockShared().Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["nickname"].String(), "name_1")

		all, err := conn.Model(table).LockShared().Where("id<=?", 3).Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 3)
	})
}

// =============================================================================
// Isolation levels
// =============================================================================

// Test_TiDB_Isolation_ReadCommitted asserts a READ COMMITTED transaction sees
// values committed by another connection while it is open.
func Test_TiDB_Isolation_ReadCommitted(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		conn, err := tidbFeatureConn("tidb_txn_mode=pessimistic")
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "SET SESSION TRANSACTION ISOLATION LEVEL READ COMMITTED")
		t.AssertNil(err)
		level, err := conn.GetValue(ctx, "SELECT @@transaction_isolation")
		t.AssertNil(err)
		t.Assert(level.String(), "READ-COMMITTED")

		tx, err := conn.Begin(ctx)
		t.AssertNil(err)
		one, err := tx.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "name_1")

		_, err = db.Model(table).Data(g.Map{"nickname": "committed"}).Where("id", 1).Update()
		t.AssertNil(err)

		// READ COMMITTED takes a fresh snapshot for every statement.
		one, err = tx.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "committed")
		t.AssertNil(tx.Commit())
	})
}

// Test_TiDB_Isolation_RepeatableRead asserts a REPEATABLE READ transaction keeps
// its snapshot when another connection commits.
func Test_TiDB_Isolation_RepeatableRead(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		conn, err := tidbFeatureConn("tidb_txn_mode=pessimistic")
		t.AssertNil(err)
		defer conn.Close(ctx)

		level, err := conn.GetValue(ctx, "SELECT @@transaction_isolation")
		t.AssertNil(err)
		t.Assert(level.String(), "REPEATABLE-READ")

		tx, err := conn.Begin(ctx)
		t.AssertNil(err)
		one, err := tx.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "name_1")

		_, err = db.Model(table).Data(g.Map{"nickname": "committed"}).Where("id", 1).Update()
		t.AssertNil(err)

		// The snapshot taken at the first read stays in effect.
		one, err = tx.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "name_1")
		t.AssertNil(tx.Commit())

		// Outside the transaction the new value is visible.
		one, err = conn.Model(table).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["nickname"].String(), "committed")
	})
}

// Test_TiDB_Isolation_SerializableRejected asserts SERIALIZABLE is refused with
// TiDB error 8048.
func Test_TiDB_Isolation_SerializableRejected(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := tidbFeatureConn("")
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE")
		t.AssertNE(err, nil)
		t.Assert(gstr.Contains(err.Error(), "Error 8048"), true)
		t.Assert(gstr.Contains(err.Error(), "SERIALIZABLE"), true)

		// The session keeps its previous level.
		level, err := conn.GetValue(ctx, "SELECT @@transaction_isolation")
		t.AssertNil(err)
		t.Assert(level.String(), "REPEATABLE-READ")
	})
}

// Test_TiDB_Isolation_SkipLevelCheck asserts tidb_skip_isolation_level_check
// lets the unsupported level through.
func Test_TiDB_Isolation_SkipLevelCheck(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := tidbFeatureConn("")
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "SET SESSION tidb_skip_isolation_level_check=1")
		t.AssertNil(err)
		_, err = conn.Exec(ctx, "SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE")
		t.AssertNil(err)

		level, err := conn.GetValue(ctx, "SELECT @@transaction_isolation")
		t.AssertNil(err)
		t.Assert(level.String(), "SERIALIZABLE")
	})
}

// =============================================================================
// Error mapping
// =============================================================================

// Test_TiDB_ErrorMapping_NoopFunction asserts TiDB error 1235 reaches the caller
// through gerror with its number and gf error code intact.
func Test_TiDB_ErrorMapping_NoopFunction(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		_, err := db.Model(table).LockShared().Where("id", 1).One()
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeDbOperationError)
		t.Assert(gstr.Contains(err.Error(), "Error 1235"), true)

		// The wrapped cause carries the raw driver message.
		t.AssertNE(gerror.Unwrap(err), nil)
		t.Assert(gstr.Contains(gerror.Unwrap(err).Error(), "1235"), true)

		// The failing statement is kept as the outermost message.
		t.Assert(gstr.Contains(gerror.Current(err).Error(), "LOCK IN SHARE MODE"), true)
	})
}

// Test_TiDB_ErrorMapping_Serializable asserts TiDB error 8048 reaches the caller
// through gerror with its number and gf error code intact.
func Test_TiDB_ErrorMapping_Serializable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		conn, err := tidbFeatureConn("")
		t.AssertNil(err)
		defer conn.Close(ctx)

		_, err = conn.Exec(ctx, "SET SESSION TRANSACTION ISOLATION LEVEL SERIALIZABLE")
		t.AssertNE(err, nil)
		t.Assert(gerror.Code(err), gcode.CodeDbOperationError)
		t.Assert(gstr.Contains(err.Error(), "Error 8048"), true)
		t.AssertNE(gerror.Unwrap(err), nil)
		t.Assert(gstr.Contains(gerror.Unwrap(err).Error(), "8048"), true)
	})
}

// =============================================================================
// JSON
// =============================================================================

// Test_TiDB_JSON_ArrowOperators asserts -> keeps the JSON quoting and ->>
// unquotes the value.
func Test_TiDB_JSON_ArrowOperators(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "json_arrow", `(
			id int PRIMARY KEY,
			d  json
		)`)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.Map{
			"id": 1,
			"d":  `{"name":"alice","age":30,"tags":["go","tidb"]}`,
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			`d->'$.name' AS quoted`,
			`d->>'$.name' AS unquoted`,
			`d->>'$.tags[1]' AS second_tag`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["quoted"].String(), `"alice"`)
		t.Assert(one["unquoted"].String(), "alice")
		t.Assert(one["second_tag"].String(), "tidb")
	})
}

// Test_TiDB_JSON_Extract asserts JSON_EXTRACT reads scalars and nested paths.
func Test_TiDB_JSON_Extract(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "json_extract", `(
			id int PRIMARY KEY,
			d  json
		)`)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"id": 1, "d": `{"user":{"name":"alice","age":30}}`},
			{"id": 2, "d": `{"user":{"name":"bob","age":25}}`},
		}).Insert()
		t.AssertNil(err)

		one, err := db.Model(table).Fields(
			`JSON_EXTRACT(d, '$.user.age') AS age`,
			`JSON_UNQUOTE(JSON_EXTRACT(d, '$.user.name')) AS name`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["age"].Int(), 30)
		t.Assert(one["name"].String(), "alice")

		// The extracted value is usable as a predicate.
		all, err := db.Model(table).
			Where("JSON_EXTRACT(d, '$.user.age') > ?", 28).
			Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 1)
		t.Assert(all[0]["id"].Int(), 1)
	})
}

// Test_TiDB_JSON_Contains asserts JSON_CONTAINS matches array members and
// object fragments.
func Test_TiDB_JSON_Contains(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "json_contains", `(
			id int PRIMARY KEY,
			d  json
		)`)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"id": 1, "d": `{"role":"admin","tags":["go","tidb"]}`},
			{"id": 2, "d": `{"role":"user","tags":["go"]}`},
			{"id": 3, "d": `{"role":"admin","tags":["rust"]}`},
		}).Insert()
		t.AssertNil(err)

		all, err := db.Model(table).
			Where(`JSON_CONTAINS(d->'$.tags', ?)`, `"go"`).
			Order("id").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[1]["id"].Int(), 2)

		count, err := db.Model(table).
			Where(`JSON_CONTAINS(d, ?)`, `{"role":"admin"}`).Count()
		t.AssertNil(err)
		t.Assert(count, 2)

		one, err := db.Model(table).Fields(
			`JSON_CONTAINS(d->'$.tags', '"tidb"') AS has_tidb`,
		).Where("id", 1).One()
		t.AssertNil(err)
		t.Assert(one["has_tidb"].Int(), 1)
	})
}

// =============================================================================
// CTE and window functions
// =============================================================================

// Test_TiDB_CTE_Basic asserts a WITH clause narrows a result set.
func Test_TiDB_CTE_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			WITH first_five AS (
				SELECT id, nickname FROM %s ORDER BY id ASC LIMIT 5
			)
			SELECT * FROM first_five WHERE id > 2 ORDER BY id
		`, table))
		t.AssertNil(err)
		t.Assert(len(all), 3)
		t.Assert(all[0]["id"].Int(), 3)
		t.Assert(all[0]["nickname"].String(), "name_3")
		t.Assert(all[2]["id"].Int(), 5)
	})
}

// Test_TiDB_CTE_Recursive asserts a recursive CTE walks a parent/child tree.
func Test_TiDB_CTE_Recursive(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "cte_tree", `(
			id        int PRIMARY KEY,
			parent_id int,
			name      varchar(50)
		)`)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"id": 1, "parent_id": nil, "name": "root"},
			{"id": 2, "parent_id": 1, "name": "child_a"},
			{"id": 3, "parent_id": 1, "name": "child_b"},
			{"id": 4, "parent_id": 2, "name": "grandchild"},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			WITH RECURSIVE tree AS (
				SELECT id, parent_id, name, 0 AS depth FROM %s WHERE id = 1
				UNION ALL
				SELECT t.id, t.parent_id, t.name, tree.depth + 1
				FROM %s t JOIN tree ON t.parent_id = tree.id
			)
			SELECT * FROM tree ORDER BY depth, id
		`, table, table))
		t.AssertNil(err)
		t.Assert(len(all), 4)
		t.Assert(all[0]["name"].String(), "root")
		t.Assert(all[0]["depth"].Int(), 0)
		t.Assert(all[3]["name"].String(), "grandchild")
		t.Assert(all[3]["depth"].Int(), 2)
	})
}

// Test_TiDB_Window_RowNumber asserts ROW_NUMBER numbers rows in window order.
func Test_TiDB_Window_RowNumber(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT id, ROW_NUMBER() OVER (ORDER BY id DESC) AS rn
			FROM %s ORDER BY id
		`, table))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		// Ascending output over a descending window: rn counts back down.
		t.Assert(all[0]["id"].Int(), 1)
		t.Assert(all[0]["rn"].Int(), TableSize)
		t.Assert(all[TableSize-1]["id"].Int(), TableSize)
		t.Assert(all[TableSize-1]["rn"].Int(), 1)
	})
}

// Test_TiDB_Window_RankAndPartition asserts RANK and DENSE_RANK differ on ties
// inside a partition.
func Test_TiDB_Window_RankAndPartition(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := tidbFeatureTable(t, "window_rank", `(
			id     int PRIMARY KEY AUTO_INCREMENT,
			dept   varchar(20),
			salary int
		)`)
		defer dropTable(table)

		_, err := db.Model(table).Data(g.List{
			{"dept": "eng", "salary": 200},
			{"dept": "eng", "salary": 100},
			{"dept": "eng", "salary": 100},
			{"dept": "sales", "salary": 250},
			{"dept": "sales", "salary": 150},
		}).Insert()
		t.AssertNil(err)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT dept, salary,
				RANK()       OVER (PARTITION BY dept ORDER BY salary DESC) AS rnk,
				DENSE_RANK() OVER (PARTITION BY dept ORDER BY salary DESC) AS dense_rnk
			FROM %s ORDER BY dept, salary DESC, rnk
		`, table))
		t.AssertNil(err)
		t.Assert(len(all), 5)

		t.Assert(all[0]["dept"].String(), "eng")
		t.Assert(all[0]["salary"].Int(), 200)
		t.Assert(all[0]["rnk"].Int(), 1)
		// The tie shares rank 2 and DENSE_RANK agrees here.
		t.Assert(all[1]["salary"].Int(), 100)
		t.Assert(all[1]["rnk"].Int(), 2)
		t.Assert(all[1]["dense_rnk"].Int(), 2)
		t.Assert(all[2]["rnk"].Int(), 2)
		t.Assert(all[2]["dense_rnk"].Int(), 2)
		// A new partition restarts the numbering.
		t.Assert(all[3]["dept"].String(), "sales")
		t.Assert(all[3]["rnk"].Int(), 1)
	})
}

// Test_TiDB_Window_RunningTotal asserts SUM OVER produces a running total.
func Test_TiDB_Window_RunningTotal(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf(`
			SELECT id,
				SUM(id)  OVER (ORDER BY id) AS running_total,
				LAG(id)  OVER (ORDER BY id) AS prev_id,
				LEAD(id) OVER (ORDER BY id) AS next_id
			FROM %s ORDER BY id
		`, table))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)

		t.Assert(all[0]["running_total"].Int(), 1)
		t.Assert(all[4]["running_total"].Int(), 15)
		t.Assert(all[9]["running_total"].Int(), 55)

		t.Assert(all[0]["prev_id"].IsNil() || all[0]["prev_id"].IsEmpty(), true)
		t.Assert(all[0]["next_id"].Int(), 2)
		t.Assert(all[9]["prev_id"].Int(), 9)
		t.Assert(all[9]["next_id"].IsNil() || all[9]["next_id"].IsEmpty(), true)
	})
}

// =============================================================================
// Optimizer hints
// =============================================================================

// Test_TiDB_Hint_UseIndex asserts USE_INDEX is honoured and shows up in the plan.
func Test_TiDB_Hint_UseIndex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		_, err := db.Exec(ctx, fmt.Sprintf("CREATE INDEX idx_passport ON %s (passport)", table))
		t.AssertNil(err)

		plan, err := db.GetAll(ctx, fmt.Sprintf(
			"EXPLAIN SELECT /*+ USE_INDEX(%s, idx_passport) */ id FROM %s WHERE passport = 'user_3'",
			table, table,
		))
		t.AssertNil(err)
		t.AssertGT(len(plan), 0)

		var planText string
		for _, row := range plan {
			planText += row["id"].String() + " "
		}
		// The hint forces the index path rather than a full table scan.
		t.Assert(gstr.Contains(planText, "IndexLookUp") || gstr.Contains(planText, "IndexReader"), true)
		t.Assert(gstr.Contains(planText, "TableFullScan"), false)

		// The hint is accepted, so no warning is raised.
		warnings, err := db.GetAll(ctx, "SHOW WARNINGS")
		t.AssertNil(err)
		t.Assert(len(warnings), 0)

		// The hinted query returns the right row.
		value, err := db.GetValue(ctx, fmt.Sprintf(
			"SELECT /*+ USE_INDEX(%s, idx_passport) */ id FROM %s WHERE passport = 'user_3'",
			table, table,
		))
		t.AssertNil(err)
		t.Assert(value.Int(), 3)
	})
}

// Test_TiDB_Hint_InljJoin asserts TIDB_INLJ selects an index join in the plan.
func Test_TiDB_Hint_InljJoin(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		left := createInitTable()
		defer dropTable(left)
		right := createInitTable()
		defer dropTable(right)

		plan, err := db.GetAll(ctx, fmt.Sprintf(
			"EXPLAIN SELECT /*+ TIDB_INLJ(t1) */ t1.id FROM %s t1 JOIN %s t2 ON t1.id = t2.id",
			left, right,
		))
		t.AssertNil(err)
		t.AssertGT(len(plan), 0)

		var planText string
		for _, row := range plan {
			planText += row["id"].String() + " "
		}
		t.Assert(gstr.Contains(planText, "IndexJoin"), true)

		warnings, err := db.GetAll(ctx, "SHOW WARNINGS")
		t.AssertNil(err)
		t.Assert(len(warnings), 0)

		// The hinted join still returns every matching row.
		all, err := db.GetAll(ctx, fmt.Sprintf(
			"SELECT /*+ TIDB_INLJ(t1) */ t1.id FROM %s t1 JOIN %s t2 ON t1.id = t2.id ORDER BY t1.id",
			left, right,
		))
		t.AssertNil(err)
		t.Assert(len(all), TableSize)
		t.Assert(all[0]["id"].Int(), 1)
	})
}

// =============================================================================
// TiDB metadata
// =============================================================================

// Test_TiDB_Version asserts tidb_version() reports a TiDB release.
func Test_TiDB_Version(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		value, err := db.GetValue(ctx, "SELECT tidb_version()")
		t.AssertNil(err)
		text := value.String()
		t.Assert(gstr.Contains(text, "Release Version"), true)
		t.Assert(gstr.Contains(text, "Git Commit Hash"), true)
		t.Assert(gstr.Contains(text, "Store"), true)

		// The MySQL compatibility version is reported separately.
		value, err = db.GetValue(ctx, "SELECT VERSION()")
		t.AssertNil(err)
		t.Assert(gstr.Contains(value.String(), "TiDB"), true)
	})
}

// Test_TiDB_ShowTableRegions asserts SHOW TABLE ... REGIONS returns region rows.
func Test_TiDB_ShowTableRegions(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		all, err := db.GetAll(ctx, fmt.Sprintf("SHOW TABLE %s REGIONS", table))
		t.AssertNil(err)
		t.AssertGT(len(all), 0)

		row := all[0]
		_, hasRegionID := row["REGION_ID"]
		t.Assert(hasRegionID, true)
		_, hasStartKey := row["START_KEY"]
		t.Assert(hasStartKey, true)
		_, hasLeaderID := row["LEADER_ID"]
		t.Assert(hasLeaderID, true)
		t.AssertGT(row["REGION_ID"].Int64(), int64(0))
	})
}

// Test_TiDB_InformationSchema_TidbIndexes asserts the TiDB index view reports
// the indexes of a table including the clustered flag.
func Test_TiDB_InformationSchema_TidbIndexes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		table := createInitTable()
		defer dropTable(table)

		_, err := db.Exec(ctx, fmt.Sprintf("CREATE INDEX idx_nickname ON %s (nickname)", table))
		t.AssertNil(err)

		all, err := db.Model("information_schema.tidb_indexes").
			Where("TABLE_SCHEMA", TestSchema1).
			Where("TABLE_NAME", table).
			Order("KEY_NAME").All()
		t.AssertNil(err)
		t.Assert(len(all), 2)

		t.Assert(all[0]["KEY_NAME"].String(), "PRIMARY")
		t.Assert(all[0]["COLUMN_NAME"].String(), "id")
		t.Assert(all[0]["NON_UNIQUE"].Int(), 0)
		t.Assert(all[0]["CLUSTERED"].String(), "YES")

		t.Assert(all[1]["KEY_NAME"].String(), "idx_nickname")
		t.Assert(all[1]["COLUMN_NAME"].String(), "nickname")
		t.Assert(all[1]["NON_UNIQUE"].Int(), 1)
		t.Assert(all[1]["CLUSTERED"].String(), "NO")
	})
}
