// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/text/gregex"
)

func Test_HookSelect_Regex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			err              error
			toBeCommittedSql = `select * from "user" where 1=1`
		)
		toBeCommittedSql, err = gregex.ReplaceStringFuncMatch(
			`(?i) FROM ([\S]+)`,
			toBeCommittedSql,
			func(match []string) string {

				return fmt.Sprintf(` FROM "%s"`, "user_1")
			},
		)
		t.AssertNil(err)
		t.Assert(toBeCommittedSql, `select * FROM "user_1" where 1=1`)
	})
	gtest.C(t, func(t *gtest.T) {
		var (
			err              error
			toBeCommittedSql = `select * from user`
		)
		toBeCommittedSql, err = gregex.ReplaceStringFuncMatch(
			`(?i) FROM ([\S]+)`,
			toBeCommittedSql,
			func(match []string) string {
				return fmt.Sprintf(` FROM %s`, "user_1")
			},
		)
		t.AssertNil(err)
		t.Assert(toBeCommittedSql, `select * FROM user_1`)
	})
}

func Test_parseConfigNodeLink_WithType(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@tcp(9.135.69.119:3306)/khaos_oss?loc=Local&parseTime=true&charset=latin`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, `khaos_oss`)
		t.Assert(newNode.Extra, `loc=Local&parseTime=true&charset=latin`)
		t.Assert(newNode.Charset, `latin`)
		t.Assert(newNode.Protocol, `tcp`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@tcp(9.135.69.119:3306)/khaos_oss?`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, `khaos_oss`)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `tcp`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@tcp(9.135.69.119:3306)/khaos_oss`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, `khaos_oss`)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `tcp`)
	})
	// empty database preselect.
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@tcp(9.135.69.119:3306)/?loc=Local&parseTime=true&charset=latin`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, ``)
		t.Assert(newNode.Extra, `loc=Local&parseTime=true&charset=latin`)
		t.Assert(newNode.Charset, `latin`)
		t.Assert(newNode.Protocol, `tcp`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@tcp(9.135.69.119:3306)?loc=Local&parseTime=true&charset=latin`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, ``)
		t.Assert(newNode.Extra, `loc=Local&parseTime=true&charset=latin`)
		t.Assert(newNode.Charset, `latin`)
		t.Assert(newNode.Protocol, `tcp`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@tcp(9.135.69.119:3306)/`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, ``)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `tcp`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@tcp(9.135.69.119:3306)`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, ``)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `tcp`)
	})
	// udp.
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `mysql:root:CxzhD*624:27jh@udp(9.135.69.119:3306)`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `mysql`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, `9.135.69.119`)
		t.Assert(newNode.Port, `3306`)
		t.Assert(newNode.Name, ``)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `udp`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `sqlite:root:CxzhD*624:27jh@file(/var/data/db.sqlite3)?local=Local&parseTime=true`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `sqlite`)
		t.Assert(newNode.User, `root`)
		t.Assert(newNode.Pass, `CxzhD*624:27jh`)
		t.Assert(newNode.Host, ``)
		t.Assert(newNode.Port, ``)
		t.Assert(newNode.Name, `/var/data/db.sqlite3`)
		t.Assert(newNode.Extra, `local=Local&parseTime=true`)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `file`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `sqlite::CxzhD*624:2@7jh@file(/var/data/db.sqlite3)`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `sqlite`)
		t.Assert(newNode.User, ``)
		t.Assert(newNode.Pass, `CxzhD*624:2@7jh`)
		t.Assert(newNode.Host, ``)
		t.Assert(newNode.Port, ``)
		t.Assert(newNode.Name, `/var/data/db.sqlite3`)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `file`)
	})
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `sqlite::@file(/var/data/db.sqlite3)`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `sqlite`)
		t.Assert(newNode.User, ``)
		t.Assert(newNode.Pass, ``)
		t.Assert(newNode.Host, ``)
		t.Assert(newNode.Port, ``)
		t.Assert(newNode.Name, `/var/data/db.sqlite3`)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `file`)
	})
	// #3146
	gtest.C(t, func(t *gtest.T) {
		node := &ConfigNode{
			Link: `pgsql:BASIC$xxxx:123456@tcp(xxxx.hologres.aliyuncs.com:80)/xxx`,
		}
		newNode := parseConfigNodeLink(node)
		t.Assert(newNode.Type, `pgsql`)
		t.Assert(newNode.User, `BASIC$xxxx`)
		t.Assert(newNode.Pass, `123456`)
		t.Assert(newNode.Host, `xxxx.hologres.aliyuncs.com`)
		t.Assert(newNode.Port, `80`)
		t.Assert(newNode.Name, `xxx`)
		t.Assert(newNode.Extra, ``)
		t.Assert(newNode.Charset, defaultCharset)
		t.Assert(newNode.Protocol, `tcp`)
	})
}

func Test_Func_doQuoteWord(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := map[string]string{
			"user":                   "`user`",
			"user u":                 "user u",
			"user_detail":            "`user_detail`",
			"user,user_detail":       "user,user_detail",
			"user u, user_detail ut": "user u, user_detail ut",
			"u.id asc":               "u.id asc",
			"u.id asc, ut.uid desc":  "u.id asc, ut.uid desc",
		}
		for k, v := range array {
			t.Assert(doQuoteWord(k, "`", "`"), v)
		}
	})
}

func Test_Func_doQuoteString(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		array := map[string]string{
			"user":                             "`user`",
			"user u":                           "`user` u",
			"user,user_detail":                 "`user`,`user_detail`",
			"user u, user_detail ut":           "`user` u,`user_detail` ut",
			"u.id, u.name, u.age":              "`u`.`id`,`u`.`name`,`u`.`age`",
			"u.id asc":                         "`u`.`id` asc",
			"u.id asc, ut.uid desc":            "`u`.`id` asc,`ut`.`uid` desc",
			"user.user u, user.user_detail ut": "`user`.`user` u,`user`.`user_detail` ut",
			// mssql global schema access with double dots.
			"user..user u, user.user_detail ut": "`user`..`user` u,`user`.`user_detail` ut",
		}
		for k, v := range array {
			t.Assert(doQuoteString(k, "`", "`"), v)
		}
	})
}

func Test_Func_addTablePrefix(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		prefix := ""
		array := map[string]string{
			"user":                         "`user`",
			"user u":                       "`user` u",
			"user as u":                    "`user` as u",
			"user,user_detail":             "`user`,`user_detail`",
			"user u, user_detail ut":       "`user` u,`user_detail` ut",
			"`user`.user_detail":           "`user`.`user_detail`",
			"`user`.`user_detail`":         "`user`.`user_detail`",
			"user as u, user_detail as ut": "`user` as u,`user_detail` as ut",
			"UserCenter.user as u, UserCenter.user_detail as ut": "`UserCenter`.`user` as u,`UserCenter`.`user_detail` as ut",
			// mssql global schema access with double dots.
			"UserCenter..user as u, user_detail as ut": "`UserCenter`..`user` as u,`user_detail` as ut",
		}
		for k, v := range array {
			t.Assert(doQuoteTableName(k, prefix, "`", "`"), v)
		}
	})
	gtest.C(t, func(t *gtest.T) {
		prefix := "gf_"
		array := map[string]string{
			"user":                         "`gf_user`",
			"user u":                       "`gf_user` u",
			"user as u":                    "`gf_user` as u",
			"user,user_detail":             "`gf_user`,`gf_user_detail`",
			"user u, user_detail ut":       "`gf_user` u,`gf_user_detail` ut",
			"`user`.user_detail":           "`user`.`gf_user_detail`",
			"`user`.`user_detail`":         "`user`.`gf_user_detail`",
			"user as u, user_detail as ut": "`gf_user` as u,`gf_user_detail` as ut",
			"UserCenter.user as u, UserCenter.user_detail as ut": "`UserCenter`.`gf_user` as u,`UserCenter`.`gf_user_detail` as ut",
			// mssql global schema access with double dots.
			"UserCenter..user as u, user_detail as ut": "`UserCenter`..`gf_user` as u,`gf_user_detail` as ut",
		}
		for k, v := range array {
			t.Assert(doQuoteTableName(k, prefix, "`", "`"), v)
		}
	})
}

func Test_isSubQuery(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		t.Assert(isSubQuery("user"), false)
		t.Assert(isSubQuery("user.uid"), false)
		t.Assert(isSubQuery("u, user.uid"), false)
		t.Assert(isSubQuery("select 1"), true)
	})
}
