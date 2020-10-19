// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/frame/g"

	"github.com/gogf/gf/test/gtest"
)

func Test_Types(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		if _, err := db.Exec(fmt.Sprintf(`
    CREATE TABLE IF NOT EXISTS types (
        id int(10) unsigned NOT NULL AUTO_INCREMENT,
        %s blob NOT NULL,
        %s binary(8) NOT NULL,
        %s date NOT NULL,
        %s time NOT NULL,
        %s decimal(5,2) NOT NULL,
        %s double NOT NULL,
        %s bit(2) NOT NULL,
        %s tinyint(1) NOT NULL,
        %s bool NOT NULL,
        PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `,
			"`blob`",
			"`binary`",
			"`date`",
			"`time`",
			"`decimal`",
			"`double`",
			"`bit`",
			"`tinyint`",
			"`bool`")); err != nil {
			gtest.Error(err)
		}
		defer dropTable("types")
		data := g.Map{
			"id":      1,
			"blob":    "i love gf",
			"binary":  []byte("abcdefgh"),
			"date":    "1880-10-24",
			"time":    "10:00:01",
			"decimal": -123.456,
			"double":  -123.456,
			"bit":     2,
			"tinyint": true,
			"bool":    false,
		}
		r, err := db.Table("types").Data(data).Insert()
		t.Assert(err, nil)
		n, _ := r.RowsAffected()
		t.Assert(n, 1)

		one, err := db.Table("types").One()
		t.Assert(err, nil)
		t.Assert(one["id"].Int(), 1)
		t.Assert(one["blob"].String(), data["blob"])
		t.Assert(one["binary"].String(), data["binary"])
		t.Assert(one["date"].String(), data["date"])
		t.Assert(one["time"].String(), data["time"])
		t.Assert(one["decimal"].String(), -123.46)
		t.Assert(one["double"].String(), data["double"])
		t.Assert(one["bit"].Int(), data["bit"])
		t.Assert(one["tinyint"].Bool(), data["tinyint"])
		t.Assert(one["tinyint"].Bool(), data["tinyint"])
	})
}
