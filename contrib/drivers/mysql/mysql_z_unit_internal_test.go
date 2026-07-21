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
	// loc values with special characters must be query-escaped.
	gtest.C(t, func(t *gtest.T) {
		configNode := &gdb.ConfigNode{
			Host:     "127.0.0.1",
			Port:     "3306",
			User:     "u",
			Pass:     "p",
			Name:     "db",
			Type:     "mysql",
			Protocol: "tcp",
			Charset:  "utf8",
			Timezone: "Asia/Shanghai",
		}
		source := configNodeToSource(configNode)
		t.Assert(source, "u:p@tcp(127.0.0.1:3306)/db?charset=utf8&loc=Asia%2FShanghai")
	})
	// Extra keeps MySQL system variables (e.g. time_zone) untouched.
	gtest.C(t, func(t *gtest.T) {
		configNode := &gdb.ConfigNode{
			Host:     "127.0.0.1",
			Port:     "3306",
			User:     "u",
			Pass:     "p",
			Name:     "db",
			Type:     "mysql",
			Protocol: "tcp",
			Charset:  "utf8",
			Timezone: "UTC",
			Extra:    "time_zone=%27%2B00%3A00%27&parseTime=true",
		}
		source := configNodeToSource(configNode)
		t.Assert(source, "u:p@tcp(127.0.0.1:3306)/db?charset=utf8&loc=UTC&time_zone=%27%2B00%3A00%27&parseTime=true")
	})
}
