// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb

import (
	"context"
	"fmt"
	"testing"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/os/gcmd"
	"github.com/gogf/gf/v2/os/gtime"
	"github.com/gogf/gf/v2/test/gtest"
)

const (
	SCHEMA     = "test_internal"
	TestDbUser = "root"
	TestDbPass = "12345678"
)

var (
	db         DB
	ctx        = context.TODO()
	configNode ConfigNode
)

func init1() {
	parser, err := gcmd.Parse(map[string]bool{
		"name": true,
		"type": true,
	}, false)
	gtest.AssertNil(err)
	configNode = ConfigNode{
		Host:             "127.0.0.1",
		Port:             "3306",
		User:             TestDbUser,
		Pass:             TestDbPass,
		Name:             parser.GetOpt("name", "").String(),
		Type:             parser.GetOpt("type", "mysql").String(),
		Role:             "master",
		Charset:          "utf8",
		Weight:           1,
		MaxIdleConnCount: 10,
		MaxOpenConnCount: 10,
		MaxConnLifeTime:  600,
	}
	AddConfigNode(DefaultGroupName, configNode)
	// Default db.
	if r, err := New(); err != nil {
		gtest.Error(err)
	} else {
		db = r
	}
	schemaTemplate := "CREATE DATABASE IF NOT EXISTS `%s` CHARACTER SET UTF8"
	if _, err := db.Exec(ctx, fmt.Sprintf(schemaTemplate, SCHEMA)); err != nil {
		gtest.Error(err)
	}
	db.SetSchema(SCHEMA)
}

func dropTable(table string) {
	if _, err := db.Exec(ctx, fmt.Sprintf("DROP TABLE IF EXISTS `%s`", table)); err != nil {
		gtest.Error(err)
	}
}

func Test_Func_FormatSqlWithArgs(t *testing.T) {
	// mysql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = FormatSqlWithArgs("select * from table where id>=? and sex=?", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// mssql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = FormatSqlWithArgs("select * from table where id>=@p1 and sex=@p2", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// pgsql
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = FormatSqlWithArgs("select * from table where id>=$1 and sex=$2", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
	})
	// oracle
	gtest.C(t, func(t *gtest.T) {
		var s string
		s = FormatSqlWithArgs("select * from table where id>=:v1 and sex=:v2", []interface{}{100, 1})
		t.Assert(s, "select * from table where id>=100 and sex=1")
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
			t.Assert(doHandleTableName(k, prefix, "`", "`"), v)
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
			t.Assert(doHandleTableName(k, prefix, "`", "`"), v)
		}
	})
}

func Test_Model_getSoftFieldName(t *testing.T) {
	table1 := "soft_deleting_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  create_at datetime DEFAULT NULL,
  update_at datetime DEFAULT NULL,
  delete_at datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table1)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table1)

	table2 := "soft_deleting_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id        int(11) NOT NULL,
  name      varchar(45) DEFAULT NULL,
  createat datetime DEFAULT NULL,
  updateat datetime DEFAULT NULL,
  deleteat datetime DEFAULT NULL,
  PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table2)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		model := db.Model(table1)
		gtest.Assert(model.getSoftFieldNameCreated(table2), "createat")
		gtest.Assert(model.getSoftFieldNameUpdated(table2), "updateat")
		gtest.Assert(model.getSoftFieldNameDeleted(table2), "deleteat")
	})
}

func Test_Model_getConditionForSoftDeleting(t *testing.T) {
	table1 := "soft_deleting_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id1        int(11) NOT NULL,
  name1      varchar(45) DEFAULT NULL,
  create_at datetime DEFAULT NULL,
  update_at datetime DEFAULT NULL,
  delete_at datetime DEFAULT NULL,
  PRIMARY KEY (id1)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table1)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table1)

	table2 := "soft_deleting_table_" + gtime.TimestampNanoStr()
	if _, err := db.Exec(ctx, fmt.Sprintf(`
CREATE TABLE %s (
  id2        int(11) NOT NULL,
  name2      varchar(45) DEFAULT NULL,
  createat datetime DEFAULT NULL,
  updateat datetime DEFAULT NULL,
  deleteat datetime DEFAULT NULL,
  PRIMARY KEY (id2)
) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `, table2)); err != nil {
		gtest.Error(err)
	}
	defer dropTable(table2)

	gtest.C(t, func(t *gtest.T) {
		model := db.Model(table1)
		t.Assert(model.getConditionForSoftDeleting(), "`delete_at` IS NULL")
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(fmt.Sprintf(`%s as t`, table1))
		t.Assert(model.getConditionForSoftDeleting(), "`delete_at` IS NULL")
	})

	gtest.C(t, func(t *gtest.T) {
		model := db.Model(fmt.Sprintf(`%s, %s`, table1, table2))
		t.Assert(model.getConditionForSoftDeleting(), fmt.Sprintf(
			"`%s`.`delete_at` IS NULL AND `%s`.`deleteat` IS NULL",
			table1, table2,
		))
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(fmt.Sprintf(`%s t1, %s as t2`, table1, table2))
		t.Assert(model.getConditionForSoftDeleting(), "`t1`.`delete_at` IS NULL AND `t2`.`deleteat` IS NULL")
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(fmt.Sprintf(`%s as t1, %s as t2`, table1, table2))
		t.Assert(model.getConditionForSoftDeleting(), "`t1`.`delete_at` IS NULL AND `t2`.`deleteat` IS NULL")
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(fmt.Sprintf(`%s as t1`, table1)).LeftJoin(table2+" t2", "t2.id2=t1.id1")
		t.Assert(model.getConditionForSoftDeleting(), "`t1`.`delete_at` IS NULL AND `t2`.`deleteat` IS NULL")
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(fmt.Sprintf(`%s`, table1)).LeftJoin(table2, "t2.id2=t1.id1")
		t.Assert(model.getConditionForSoftDeleting(), fmt.Sprintf(
			"`%s`.`delete_at` IS NULL AND `%s`.`deleteat` IS NULL",
			table1, table2,
		))
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(fmt.Sprintf(`%s`, table1)).LeftJoin(table2, "t2.id2=t1.id1").RightJoin(table2, "t2.id2=t1.id1")
		t.Assert(model.getConditionForSoftDeleting(), fmt.Sprintf(
			"`%s`.`delete_at` IS NULL AND `%s`.`deleteat` IS NULL AND `%s`.`deleteat` IS NULL",
			table1, table2, table2,
		))
	})
	gtest.C(t, func(t *gtest.T) {
		model := db.Model(table1+" as t1").LeftJoin(table2+" as t2", "t2.id2=t1.id1").RightJoin(table2+" as t3 ", "t2.id2=t1.id1")
		t.Assert(
			model.getConditionForSoftDeleting(),
			"`t1`.`delete_at` IS NULL AND `t2`.`deleteat` IS NULL AND `t3`.`deleteat` IS NULL",
		)
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

func TestResult_Structs1(t *testing.T) {
	type A struct {
		Id int `orm:"id"`
	}
	type B struct {
		*A
		Name string
	}
	gtest.C(t, func(t *gtest.T) {
		r := Result{
			Record{"id": gvar.New(nil), "name": gvar.New("john")},
			Record{"id": gvar.New(nil), "name": gvar.New("smith")},
		}
		array := make([]*B, 2)
		err := r.Structs(&array)
		t.AssertNil(err)
		t.Assert(array[0].Id, 0)
		t.Assert(array[1].Id, 0)
		t.Assert(array[0].Name, "john")
		t.Assert(array[1].Name, "smith")
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_PtrAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One *S1
		}
		var (
			s   []*S3
			err error
		)
		r1 := Result{
			Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := Result{
			Record{
				"id":  gvar.New(1),
				"age": gvar.New(20),
			},
			Record{
				"id":  gvar.New(2),
				"age": gvar.New(21),
			},
		}
		err = r2.ScanList(&s, "One", "One", "id:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 20)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 21)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_StructAttribute(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One S1
		}
		var (
			s   []*S3
			err error
		)
		r1 := Result{
			Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := Result{
			Record{
				"id":  gvar.New(1),
				"age": gvar.New(20),
			},
			Record{
				"id":  gvar.New(2),
				"age": gvar.New(21),
			},
		}
		err = r2.ScanList(&s, "One", "One", "id:Id")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 20)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 21)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_SliceAttribute_Ptr(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S2 struct {
			Id    int
			Pid   int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One  *S1
			Many []*S2
		}
		var (
			s   []*S3
			err error
		)
		r1 := Result{
			Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := Result{
			Record{
				"id":   gvar.New(100),
				"pid":  gvar.New(1),
				"age":  gvar.New(30),
				"name": gvar.New("john"),
			},
			Record{
				"id":   gvar.New(200),
				"pid":  gvar.New(1),
				"age":  gvar.New(31),
				"name": gvar.New("smith"),
			},
		}
		err = r2.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 30)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 31)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)

		r3 := Result{
			Record{
				"id":  gvar.New(100),
				"pid": gvar.New(1),
				"age": gvar.New(40),
			},
			Record{
				"id":  gvar.New(200),
				"pid": gvar.New(1),
				"age": gvar.New(41),
			},
		}
		err = r3.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 40)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 41)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)
	})
}

// https://github.com/gogf/gf/issues/1159
func Test_ScanList_NoRecreate_SliceAttribute_Struct(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		type S1 struct {
			Id    int
			Name  string
			Age   int
			Score int
		}
		type S2 struct {
			Id    int
			Pid   int
			Name  string
			Age   int
			Score int
		}
		type S3 struct {
			One  S1
			Many []S2
		}
		var (
			s   []S3
			err error
		)
		r1 := Result{
			Record{
				"id":   gvar.New(1),
				"name": gvar.New("john"),
				"age":  gvar.New(16),
			},
			Record{
				"id":   gvar.New(2),
				"name": gvar.New("smith"),
				"age":  gvar.New(18),
			},
		}
		err = r1.ScanList(&s, "One")
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)

		r2 := Result{
			Record{
				"id":   gvar.New(100),
				"pid":  gvar.New(1),
				"age":  gvar.New(30),
				"name": gvar.New("john"),
			},
			Record{
				"id":   gvar.New(200),
				"pid":  gvar.New(1),
				"age":  gvar.New(31),
				"name": gvar.New("smith"),
			},
		}
		err = r2.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 30)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 31)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)

		r3 := Result{
			Record{
				"id":  gvar.New(100),
				"pid": gvar.New(1),
				"age": gvar.New(40),
			},
			Record{
				"id":  gvar.New(200),
				"pid": gvar.New(1),
				"age": gvar.New(41),
			},
		}
		err = r3.ScanList(&s, "Many", "One", "pid:Id")
		// fmt.Printf("%+v", err)
		t.AssertNil(err)
		t.Assert(len(s), 2)
		t.Assert(s[0].One.Name, "john")
		t.Assert(s[0].One.Age, 16)
		t.Assert(len(s[0].Many), 2)
		t.Assert(s[0].Many[0].Name, "john")
		t.Assert(s[0].Many[0].Age, 40)
		t.Assert(s[0].Many[1].Name, "smith")
		t.Assert(s[0].Many[1].Age, 41)

		t.Assert(s[1].One.Name, "smith")
		t.Assert(s[1].One.Age, 18)
		t.Assert(len(s[1].Many), 0)
	})
}
