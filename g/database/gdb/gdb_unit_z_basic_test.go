// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"testing"

	"github.com/gogf/gf/g/database/gdb"
	"github.com/gogf/gf/g/test/gtest"
)

func Test_Instance(t *testing.T) {
	gtest.Case(t, func() {
		_, err := gdb.Instance("none")
		gtest.AssertNE(err, nil)

		db, err := gdb.Instance()
		gtest.Assert(err, nil)

		err1 := db.PingMaster()
		err2 := db.PingSlave()
		gtest.Assert(err1, nil)
		gtest.Assert(err2, nil)
	})
}

/*
func Test_Config(t *testing.T) {

	gtest.Case(t, func() {
		//gdb.SetDefaultGroup(gdb.DEFAULT_GROUP_NAME)
		gtest.Assert(gdb.GetDefaultGroup(), gdb.DEFAULT_GROUP_NAME)
	})

	gtest.Case(t, func() {
		nodeConfig := gdb.Config{
			"mysqltest1": gdb.ConfigGroup{
				gdb.ConfigNode{
					Host:     "127.0.0.1",
					Port:     "3306",
					User:     "root",
					Pass:     "",
					Name:     "",
					Type:     "mysql",
					Role:     "master",
					Charset:  "utf8",
					Weight:	1,
					MaxIdleConnCount:	10,
					MaxOpenConnCount:	10,
					MaxConnLifetime:	600,
				},
			},
		}

		gdb.SetConfig(nodeConfig)

		groupConfig := gdb.ConfigGroup{
			gdb.ConfigNode{
				Host:     "127.0.0.1",
				Port:     "3306",
				User:     "root",
				Pass:     "",
				Name:     "",
				Type:     "mysql",
				Role:     "master",
				Charset:  "utf8",
				LinkInfo: "root:@tcp(127.0.0.1:3306)/test",
				Weight:	1,
				MaxIdleConnCount:	10,
				MaxOpenConnCount:	10,
				MaxConnLifetime:	600,
			},
		}

		gdb.AddConfigGroup("mysqltest2", groupConfig)



		res := gdb.GetConfig("mysqltest2")
		gtest.Assert(res[0].Host, groupConfig[0].Host)
		gtest.Assert(res[0].Port, groupConfig[0].Port)

		gtest.Assert(res[0].String(), groupConfig[0].LinkInfo)
	})
}*/
