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
	TableNamePrefix1 = "gf_"
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

	// db.SetDebug(true)
}

func dropTable(table string) {
	dropTableWithDb(db, table)
}

// dropTableWithDb drops `table` in the schema of `db`, both as the upper case name an unquoted
// identifier folds to and as written, with the sequence that createUserTableWithDb creates for
// it, and clears its cached fields, as a table of the same name may be created again.
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
	if !strings.Contains(table, `"`) {
		if _, err := db.Exec(ctx, fmt.Sprintf(`DROP SEQUENCE IF EXISTS %s_ID_SEQ`, table)); err != nil {
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

// createUserTable creates the table that the ported test suites use, see createUserTableWithDb.
func createUserTable(table ...string) string {
	return createUserTableWithDb(db, table...)
}

// createInitUserTable creates the table that the ported test suites use with TableSize records.
func createInitUserTable(table ...string) string {
	return createInitUserTableWithDb(db, table...)
}

// createUserTableWithDb creates, in the schema of `db`, the table that the ported test suites
// use, named `table` or else a unique name, with a sequence and a trigger that fill in its ID
// when an insert gives none.
func createUserTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
	} else {
		name = fmt.Sprintf("user_%d", gtime.TimestampMicro())
	}
	dropTableWithDb(db, name)
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        ID NUMBER(10) NOT NULL,
        PASSPORT VARCHAR(45) NOT NULL,
        PASSWORD CHAR(32) NOT NULL,
        NICKNAME VARCHAR(45) NOT NULL,
        CREATE_TIME VARCHAR(45),
        SALARY NUMBER(18,2),
        PRIMARY KEY (ID)
    )`, name)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE SEQUENCE %s_ID_SEQ
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9999999999
    NOCYCLE
    NOCACHE`, name)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE OR REPLACE TRIGGER %s_ID_TRG
BEFORE INSERT ON %s
FOR EACH ROW
BEGIN
    IF :NEW.ID IS NULL THEN
        :NEW.ID := %s_ID_SEQ.NEXTVAL;
    END IF;
END;`, name, name, name)); err != nil {
		gtest.Fatal(err)
	}
	return
}

// createInitUserTableWithDb creates the table of createUserTableWithDb with TableSize records.
func createInitUserTableWithDb(db gdb.DB, table ...string) (name string) {
	name = createUserTableWithDb(db, table...)
	var (
		values = make([]string, 0, TableSize)
		args   = make([]any, 0, TableSize*5)
	)
	for i := 1; i <= TableSize; i++ {
		values = append(values, fmt.Sprintf(
			"INTO %s (ID, PASSPORT, PASSWORD, NICKNAME, CREATE_TIME) VALUES (?, ?, ?, ?, ?)", name,
		))
		args = append(args,
			i,
			fmt.Sprintf(`user_%d`, i),
			fmt.Sprintf(`pass_%d`, i),
			fmt.Sprintf(`name_%d`, i),
			gtime.Now().String(),
		)
	}
	result, err := db.Exec(ctx, "INSERT ALL "+strings.Join(values, " ")+" SELECT 1 FROM DUAL", args...)
	gtest.AssertNil(err)
	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)
	_, err = db.TableFields(ctx, name)
	gtest.AssertNil(err)
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
