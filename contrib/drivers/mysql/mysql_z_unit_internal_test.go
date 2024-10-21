// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package mysql

import (
	"testing"

	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_configNodeToSource(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		configNode := &gdb.ConfigNode{
			Host:     "/tmp/mysql.sock",
			Port:     "",
			User:     "username",
			Pass:     "password",
			Name:     "dbname",
			Type:     "mysql",
			Protocol: "unix",
		}
		source := configNodeToSource(configNode)
		t.Assert(source, "username:password@unix(/tmp/mysql.sock)/dbname?charset=")
	})
}
