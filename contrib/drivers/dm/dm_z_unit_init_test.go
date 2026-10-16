// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"gitee.com/chunanyong/dm"

	"github.com/gogf/gf/v2/container/garray"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	db        gdb.DB
	db2       gdb.DB
	dbPrefix  gdb.DB
	dblink    gdb.DB
	dbErr     gdb.DB
	ctx       context.Context
	TableSize = 10
)

const (
	TestDBHost  = "127.0.0.1"
	TestDBPort  = "5236"
	TestDBUser  = "SYSDBA"
	TestDBPass  = "SYSDBA001"
	TestDBName  = "SYSDBA"
	TestDBType  = "dm"
	TestCharset = "utf8"
)

const (
	TableName        = "t_user"
	TestSchema1      = "test1"
	TestSchema2      = "test2"
	TableNamePrefix1 = "gf_"
	CreateTime       = "2018-10-24 10:00:00"
)

const (
	dmErrObjectExists = -2124
)

type User struct {
	ID          int64     `orm:"id"`
	AccountName string    `orm:"account_name"`
	PwdReset    int64     `orm:"pwd_reset"`
	AttrIndex   int64     `orm:"attr_index"`
	Enabled     int64     `orm:"enabled"`
	Deleted     int64     `orm:"deleted"`
	CreatedBy   string    `orm:"created_by"`
	CreatedTime time.Time `orm:"created_time"`
	UpdatedBy   string    `orm:"updated_by"`
	UpdatedTime time.Time `orm:"updated_time"`
}

func init() {
	node := gdb.ConfigNode{
		Host:             TestDBHost,
		Port:             TestDBPort,
		User:             TestDBUser,
		Pass:             TestDBPass,
		Name:             TestDBName,
		Type:             TestDBType,
		Role:             "master",
		Charset:          TestCharset,
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		// CreatedAt:        "created_time",
		// UpdatedAt:        "updated_time",
	}

	nodeLink := gdb.ConfigNode{
		Type: TestDBType,
		Name: TestDBName,
		Link: fmt.Sprintf(
			"dm:%s:%s@tcp(%s:%s)/%s?charset=%s",
			TestDBUser, TestDBPass, TestDBHost, TestDBPort, TestDBName, TestCharset,
		),
	}

	nodeErr := gdb.ConfigNode{
		Host:    TestDBHost,
		Port:    TestDBPort,
		User:    TestDBUser,
		Pass:    "1234",
		Name:    TestDBName,
		Type:    TestDBType,
		Role:    "master",
		Charset: TestCharset,
		Weight:  1,
	}

	gdb.AddConfigNode(gdb.DefaultGroupName, node)
	if r, err := gdb.New(node); err != nil {
		gtest.Fatal(err)
	} else {
		db = r
	}

	gdb.AddConfigNode("dblink", nodeLink)
	if r, err := gdb.New(nodeLink); err != nil {
		gtest.Fatal(err)
	} else {
		dblink = r
	}

	gdb.AddConfigNode("dbErr", nodeErr)
	if r, err := gdb.New(nodeErr); err != nil {
		gtest.Fatal(err)
	} else {
		dbErr = r
	}

	ctx = context.Background()

	for _, schema := range []string{TestSchema1, TestSchema2} {
		if _, err := db.Exec(ctx, fmt.Sprintf(`CREATE SCHEMA %s`, schema)); err != nil && !isDmError(err, dmErrObjectExists) {
			gtest.Fatal(err)
		}
	}
	db2 = db.Schema(TestSchema2)

	nodePrefix := node
	nodePrefix.Prefix = TableNamePrefix1
	gdb.AddConfigNode("prefix", nodePrefix)
	if r, err := gdb.New(nodePrefix); err != nil {
		gtest.Fatal(err)
	} else {
		dbPrefix = r
	}

	// db.SetDebug(true)
}

func dropTable(table string) {
	dropTableWithDb(db, table)
}

// dropTableWithDb drops `table` in the schema of `db`, both as the upper case name an unquoted
// identifier folds to and as written, and clears its cached fields, as a table of the same name
// may be created again.
func dropTableWithDb(db gdb.DB, table string) {
	metadataTable := strings.Trim(table, `"`)
	tableNameMap := map[string]struct{}{
		strings.ToUpper(metadataTable): {},
		metadataTable:                  {},
	}
	for tableName := range tableNameMap {
		if _, err := db.Exec(ctx, fmt.Sprintf(`DROP TABLE IF EXISTS "%s"`, strings.ReplaceAll(tableName, `"`, `""`))); err != nil {
			gtest.Fatal(err)
		}
	}
	if err := db.GetCore().ClearTableFields(ctx, table); err != nil {
		gtest.Fatal(err)
	}
}

// isDmError reports whether `err` is an error of the DM server with the error code `code`.
func isDmError(err error, code int32) bool {
	var dmError *dm.DmError
	return errors.As(err, &dmError) && dmError.ErrCode == code
}

// createUserTable creates the table that the MySQL test suites use, see createUserTableWithDb.
func createUserTable(table ...string) string {
	return createUserTableWithDb(db, table...)
}

// createInitUserTable creates the table that the MySQL test suites use with TableSize records.
func createInitUserTable(table ...string) string {
	return createInitUserTableWithDb(db, table...)
}

// createUserTableWithDb creates, in the schema of `db`, the table that the MySQL test suites use,
// named `table` or else a unique name.
func createUserTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf(`%s_%d`, TableName, gtime.TimestampNano())
	}
	dropTableWithDb(db, name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s (
		ID          INT AUTO_INCREMENT NOT NULL,
		PASSPORT    VARCHAR(45) NULL,
		PASSWORD    CHAR(32) NULL,
		NICKNAME    VARCHAR(45) NULL,
		CREATE_TIME TIMESTAMP(6) NULL,
		CREATE_DATE DATE NULL,
		PRIMARY KEY (ID)
	)`, name)); err != nil {
		gtest.Fatal(err)
	}
	return name
}

// createInitUserTableWithDb creates the table of createUserTableWithDb with TableSize records.
func createInitUserTableWithDb(db gdb.DB, table ...string) (name string) {
	name = createUserTableWithDb(db, table...)
	list := make(g.List, 0, TableSize)
	for i := 1; i <= TableSize; i++ {
		list = append(list, g.Map{
			"id":          i,
			"passport":    fmt.Sprintf(`user_%d`, i),
			"password":    fmt.Sprintf(`pass_%d`, i),
			"nickname":    fmt.Sprintf(`name_%d`, i),
			"create_time": gtime.NewFromStr(CreateTime).String(),
		})
	}
	result, err := db.Insert(ctx, name, list)
	gtest.AssertNil(err)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func createTable(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("random_%d", gtime.Timestamp())
	}

	dropTable(name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s
(
"ID" BIGINT NOT NULL,
"ACCOUNT_NAME" VARCHAR(128) DEFAULT '' NOT NULL COMMENT 'Account Name',
"PWD_RESET" TINYINT DEFAULT 0 NOT NULL,
"ENABLED" INT DEFAULT 1 NOT NULL,
"DELETED" INT DEFAULT 0 NOT NULL,
"ATTR_INDEX" INT DEFAULT 0 ,
"CREATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"CREATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
"UPDATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"UPDATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
NOT CLUSTER PRIMARY KEY("ID")) STORAGE(ON "MAIN", CLUSTERBTR) ;
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return
}

func createInitTable(table ...string) (name string) {
	name = createTable(table...)
	array := garray.New(true)
	for i := 1; i <= TableSize; i++ {
		array.Append(g.Map{
			"id":           i,
			"account_name": fmt.Sprintf(`name_%d`, i),
			"pwd_reset":    0,
			"attr_index":   i,
			"created_time": gtime.Now(),
		})
	}
	result, err := db.Schema(TestDBName).Insert(context.Background(), name, array.Slice())
	gtest.AssertNil(err)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	return
}

func createTableFalse(table ...string) (name string, err error) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("random_%d", gtime.Timestamp())
	}

	dropTable(name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s
(
"ID" BIGINT NOT NULL,
"ACCOUNT_NAME" VARCHAR(128) DEFAULT '' NOT NULL,
"PWD_RESET" TINYINT DEFAULT 0 NOT NULL,
"ENABLED" INT DEFAULT 1 NOT NULL,
"DELETED" INT DEFAULT 0 NOT NULL,
INDEX INT DEFAULT 0 ,
"ATTR_INDEX" INT DEFAULT 0 ,
"CREATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"CREATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
"UPDATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"UPDATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
NOT CLUSTER PRIMARY KEY("ID")) STORAGE(ON "MAIN", CLUSTERBTR) ;
	`, name)); err != nil {
		// gtest.Fatal(err)
		return name, fmt.Errorf("createTableFalse")
	}

	return name, nil
}

func createInitTables(len int) []string {
	tables := make([]string, 0, len)
	for range len {
		tables = append(tables, createInitTable())
	}
	return tables
}

// createTableWithIdentity creates a table with IDENTITY column for LastInsertId testing
func createTableWithIdentity(table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("random_%d", gtime.Timestamp())
	}

	dropTable(name)

	if _, err := db.Exec(ctx, fmt.Sprintf(`
	CREATE TABLE %s
(
"ID" BIGINT IDENTITY(1, 1) NOT NULL,
"ACCOUNT_NAME" VARCHAR(128) DEFAULT '' NOT NULL COMMENT 'Account Name',
"PWD_RESET" TINYINT DEFAULT 0 NOT NULL,
"ENABLED" INT DEFAULT 1 NOT NULL,
"DELETED" INT DEFAULT 0 NOT NULL,
"ATTR_INDEX" INT DEFAULT 0 ,
"CREATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"CREATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
"UPDATED_BY" VARCHAR(32) DEFAULT '' NOT NULL,
"UPDATED_TIME" TIMESTAMP(6) DEFAULT CURRENT_TIMESTAMP() NOT NULL,
NOT CLUSTER PRIMARY KEY("ID")) STORAGE(ON "MAIN", CLUSTERBTR) ;
	`, name)); err != nil {
		gtest.Fatal(err)
	}
	return
}
