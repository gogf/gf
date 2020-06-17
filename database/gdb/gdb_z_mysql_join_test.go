// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/os/gtime"
	"github.com/gogf/gf/test/gtest"
)

func Test_Table_Join(t *testing.T) {
	var (
		tableUser       = "user_" + gtime.TimestampMicroStr()
		tableUserDetail = "user_detail_" + gtime.TimestampMicroStr()
	)
	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  uid int(10) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(45) NOT NULL COMMENT,
  PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, tableUser)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUser)

	if _, err := db.Exec(fmt.Sprintf(`
CREATE TABLE %s (
  uid int(10) unsigned NOT NULL AUTO_INCREMENT,
  name varchar(45) NOT NULL COMMENT,
  PRIMARY KEY (uid)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, tableUserDetail)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(tableUserDetail)

	gtest.C(t, func(t *gtest.T) {
		type EntityUser struct {
			Uid  int
			Name string
		}
		type EntityUserDetail struct {
			Uid      int
			TrueName string
		}
		type Entity struct {
			User       *EntityUser
			UserDetail *EntityUserDetail
		}

	})
}
