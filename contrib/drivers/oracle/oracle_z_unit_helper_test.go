// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package oracle_test

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
