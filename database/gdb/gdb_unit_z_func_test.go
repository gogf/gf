// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"github.com/gogf/gf/test/gtest"
	"testing"
)

func Test_Func_bindArgsToQuery(t *testing.T) {
	// mysql
	gtest.Case(t, func() {
		var s string
		s = bindArgsToQuery("select * from table where id>=? and sex=?", []interface{}{100, 1})
		gtest.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// mssql
	gtest.Case(t, func() {
		var s string
		s = bindArgsToQuery("select * from table where id>=@p1 and sex=@p2", []interface{}{100, 1})
		gtest.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// pgsql
	gtest.Case(t, func() {
		var s string
		s = bindArgsToQuery("select * from table where id>=$1 and sex=$2", []interface{}{100, 1})
		gtest.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// oracle
	gtest.Case(t, func() {
		var s string
		s = bindArgsToQuery("select * from table where id>=:1 and sex=:2", []interface{}{100, 1})
		gtest.Assert(s, "select * from table where id>=100 and sex=1")
	})
}

func Test_Func_doQuoteWord(t *testing.T) {
	gtest.Case(t, func() {
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
			gtest.Assert(doQuoteWord(k, "`", "`"), v)
		}
	})
}

func Test_Func_doQuoteString(t *testing.T) {
	gtest.Case(t, func() {
		// "user", "user u", "user,user_detail", "user u, user_detail ut", "u.id asc".
		array := map[string]string{
			"user":                             "`user`",
			"user u":                           "`user` u",
			"user,user_detail":                 "`user`,`user_detail`",
			"user u, user_detail ut":           "`user` u,`user_detail` ut",
			"u.id asc":                         "`u`.`id` asc",
			"u.id asc, ut.uid desc":            "`u`.`id` asc,`ut`.`uid` desc",
			"user.user u, user.user_detail ut": "`user`.`user` u,`user`.`user_detail` ut",
			// mssql global schema access with double dots.
			"user..user u, user.user_detail ut": "`user`..`user` u,`user`.`user_detail` ut",
		}
		for k, v := range array {
			gtest.Assert(doQuoteString(k, "`", "`"), v)
		}
	})
}

func Test_Func_addTablePrefix(t *testing.T) {
	gtest.Case(t, func() {
		prefix := ""
		array := map[string]string{
			"user":                         "`user`",
			"user u":                       "`user` u",
			"user as u":                    "`user` as u",
			"user,user_detail":             "`user`,`user_detail`",
			"user u, user_detail ut":       "`user` u,`user_detail` ut",
			"user as u, user_detail as ut": "`user` as u,`user_detail` as ut",
			"UserCenter.user as u, UserCenter.user_detail as ut": "`UserCenter`.`user` as u,`UserCenter`.`user_detail` as ut",
			// mssql global schema access with double dots.
			"UserCenter..user as u, user_detail as ut": "`UserCenter`..`user` as u,`user_detail` as ut",
		}
		for k, v := range array {
			gtest.Assert(doHandleTableName(k, prefix, "`", "`"), v)
		}
	})
	gtest.Case(t, func() {
		prefix := "gf_"
		array := map[string]string{
			"user":                         "`gf_user`",
			"user u":                       "`gf_user` u",
			"user as u":                    "`gf_user` as u",
			"user,user_detail":             "`gf_user`,`gf_user_detail`",
			"user u, user_detail ut":       "`gf_user` u,`gf_user_detail` ut",
			"user as u, user_detail as ut": "`gf_user` as u,`gf_user_detail` as ut",
			"UserCenter.user as u, UserCenter.user_detail as ut": "`UserCenter`.`gf_user` as u,`UserCenter`.`gf_user_detail` as ut",
			// mssql global schema access with double dots.
			"UserCenter..user as u, user_detail as ut": "`UserCenter`..`gf_user` as u,`gf_user_detail` as ut",
		}
		for k, v := range array {
			gtest.Assert(doHandleTableName(k, prefix, "`", "`"), v)
		}
	})
}
