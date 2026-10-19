// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package dm_test

import (
	"fmt"
	"strings"

	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const CreateTime = "2018-10-24 10:00:00"

func setCreateTime(table string) string {
	if _, err := db.Model(table).Data("create_time", CreateTime).Where("id>?", 0).Update(); err != nil {
		gtest.Fatal(err)
	}
	return table
}

func createNullableTable() string {
	table := fmt.Sprintf("user_%d", gtime.TimestampMicro())
	if _, err := db.Exec(ctx, fmt.Sprintf(`
    CREATE TABLE %s (
        ID NUMBER(10) NOT NULL,
        PASSPORT VARCHAR(45),
        PASSWORD CHAR(32),
        NICKNAME VARCHAR(45),
        CREATE_TIME VARCHAR(45),
        SALARY NUMBER(18,2),
        PRIMARY KEY (ID)
    )`, table)); err != nil {
		gtest.Fatal(err)
	}
	return table
}

func makeNullable(t *gtest.T, table string, columns ...string) {
	modifications := make([]string, len(columns))
	for i, column := range columns {
		modifications[i] = column + " NULL"
	}
	_, err := db.Exec(ctx, fmt.Sprintf("ALTER TABLE %s MODIFY (%s)", table, strings.Join(modifications, ", ")))
	t.AssertNil(err)
}

func createAutoIncrement(table, column string, start int) {
	if _, err := db.Exec(ctx, fmt.Sprintf(
		`CREATE SEQUENCE %s_ID_SEQ START WITH %d INCREMENT BY 1 NOCACHE`, table, start,
	)); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE OR REPLACE TRIGGER %s_ID_TRG
BEFORE INSERT ON %s
FOR EACH ROW
BEGIN
    IF :NEW.%s IS NULL THEN
        :NEW.%s := %s_ID_SEQ.NEXTVAL;
    END IF;
END;`, table, table, column, column, table)); err != nil {
		gtest.Fatal(err)
	}
}

// The codes of the DM server errors that the tests expect.
const (
	dmErrSyntax              = -2007
	dmErrInvalidTable        = -2106
	dmErrInvalidColumn       = -2111
	dmErrInvalidSavepoint    = -2121
	dmErrObjectExists        = -2124
	dmErrInvalidMember       = -2207
	dmErrInvalidLogin        = -2501
	dmErrInvalidHaving       = -4028
	dmErrNotGroupBy          = -4080
	dmErrValueAsCondition    = -4104
	dmErrForUpdateNotAllowed = -4596
	dmErrTypeMismatch        = -6105
	dmErrDeadlock            = -6403
	dmErrLockWait            = -6409
	dmErrReadOnly            = -6506
	dmErrUniqueViolation     = -6602
	dmErrCheckViolation      = -6604
	dmErrNotNullViolation    = -6609
)

// assertDmError asserts that `err` is an error of the DM server with the error code `code`.
func assertDmError(t *gtest.T, err error, code int32) {
	t.AssertNE(err, nil)
	if !isDmError(err, code) {
		t.Error(fmt.Sprintf("expected DM error %d, got: %v", code, err))
	}
}
