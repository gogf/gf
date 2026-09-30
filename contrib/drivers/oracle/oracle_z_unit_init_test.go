// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

import (
	"context"
	"errors"
	"fmt"
	"strings"

	_ "github.com/sijms/go-ora/v2"

	"github.com/sijms/go-ora/v2/network"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

var (
	db     gdb.DB
	dblink gdb.DB
	dbErr  gdb.DB
	ctx    context.Context
)

const (
	TableSize        = 10
	TableName        = "t_user"
	TestSchema1      = "test1"
	TestSchema2      = "test2"
	TableNamePrefix1 = "gf_"
	TestSchema       = "XE"
)

const (
	oracleMaxIdentifierLength = 30
	oracleErrTableNotExist    = 942
	oracleErrSequenceNotExist = 2289
)

const (
	TestDbIP   = "127.0.0.1"
	TestDbPort = "1521"
	TestDbUser = "system"
	TestDbPass = "oracle"
	TestDbName = "XE"
	TestDbType = "oracle"
)

func init() {
	node := gdb.ConfigNode{
		Host:             TestDbIP,
		Port:             TestDbPort,
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             TestDbName,
		Type:             TestDbType,
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
	}

	nodeLink := gdb.ConfigNode{
		Type: TestDbType,
		Name: TestDbName,
		Link: fmt.Sprintf("%s:%s:%s@tcp(%s:%s)/%s",
			TestDbType, TestDbUser, TestDbPass, TestDbIP, TestDbPort, TestDbName,
		),
	}

	nodeErr := gdb.ConfigNode{
		Host:    TestDbIP,
		Port:    TestDbPort,
		User:    TestDbUser,
		Pass:    "1234",
		Name:    TestDbName,
		Type:    TestDbType,
		Role:    "master",
		Charset: "utf8",
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
}

func createTable(table ...string) string {
	return createTableWithDb(db, table...)
}

func createInitTable(table ...string) string {
	return createInitTableWithDb(db, table...)
}

func dropTable(table string) {
	dropTableWithDb(db, table)
}

func createTableWithDb(db gdb.DB, table ...string) (name string) {
	if len(table) > 0 {
		name = table[0]
		dropTableWithDb(db, name)
	} else {
		name = fmt.Sprintf("user_%d", gtime.TimestampMicro())
	}

	// Step 1: Create table
	createTableSQL := fmt.Sprintf(`
    CREATE TABLE %s (
        ID NUMBER(10) NOT NULL,
        PASSPORT VARCHAR(45) NOT NULL,
        PASSWORD CHAR(32) NOT NULL,
        NICKNAME VARCHAR(45) NOT NULL,
        CREATE_TIME VARCHAR(45),
        SALARY NUMBER(18,2),
        PRIMARY KEY (ID)
    )`, name)

	if _, err := db.Exec(ctx, createTableSQL); err != nil {
		gtest.Fatal(err)
	}

	// Step 2: Create sequence
	createSeqSQL := fmt.Sprintf(`
    CREATE SEQUENCE %s_ID_SEQ
    START WITH 1
    INCREMENT BY 1
    MINVALUE 1
    MAXVALUE 9999999999
    NOCYCLE
    NOCACHE`, name)

	if _, err := db.Exec(ctx, createSeqSQL); err != nil {
		gtest.Fatal(err)
	}

	// Step 3: Create trigger - only set ID from sequence when it's NULL
	createTriggerSQL := fmt.Sprintf(`
CREATE OR REPLACE TRIGGER %s_ID_TRG
BEFORE INSERT ON %s
FOR EACH ROW
BEGIN
    IF :NEW.ID IS NULL THEN
        :NEW.ID := %s_ID_SEQ.NEXTVAL;
    END IF;
END;`, name, name, name)

	if _, err := db.Exec(ctx, createTriggerSQL); err != nil {
		gtest.Fatal(err)
	}

	// db.Schema("test")
	return
}

func createInitTableWithDb(db gdb.DB, table ...string) (name string) {
	name = createTableWithDb(db, table...)
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
	result, err := db.Exec(context.Background(), "INSERT ALL "+strings.Join(values, " ")+" SELECT 1 FROM DUAL", args...)
	gtest.AssertNil(err)

	n, e := result.RowsAffected()
	gtest.Assert(e, nil)
	gtest.Assert(n, TableSize)

	_, err = db.TableFields(context.Background(), name)
	gtest.AssertNil(err)
	return
}

func dropTableWithDb(db gdb.DB, table string) {
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE %s", table)); err != nil && !isOracleError(err, oracleErrTableNotExist) {
		gtest.Fatal(err)
	}
	sequence := table + "_ID_SEQ"
	if len(sequence) > oracleMaxIdentifierLength {
		return
	}
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP SEQUENCE %s", sequence)); err != nil && !isOracleError(err, oracleErrSequenceNotExist) {
		gtest.Fatal(err)
	}
}

func isOracleError(err error, code int) bool {
	var oracleError *network.OracleError
	return errors.As(err, &oracleError) && oracleError.ErrCode == code
}
