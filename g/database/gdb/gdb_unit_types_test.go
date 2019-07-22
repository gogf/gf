// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/g/test/gtest"
)

func Test_Types(t *testing.T) {
	if _, err := db.Exec(fmt.Sprintf(`
    CREATE TABLE types (
        id      int(10) unsigned NOT NULL AUTO_INCREMENT,
        blob    blob NOT NULL,
        binary  binary(8) NOT NULL ,
        date    varchar(45) NOT NULL ,
        decimal decimal(5,2) NOT NULL ',
        double  double NOT NULL ',
        bit     bit(2) NOT NULL ',
        PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `)); err != nil {
		gtest.Error(err)
	}

	gtest.Case(t, func() {

	})
}
