// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gdb_test

import (
	"github.com/gogf/gf/g"
	"github.com/gogf/gf/g/os/gtime"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

// 基本测试
func TestModel_Insert(t *testing.T) {
    result, err := db.Table("user").Filter().Data(g.Map{
        "id"          : 1,
        "uid"         : 1,
        "passport"    : "t1",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T1",
        "create_time" : gtime.Now().String(),
    }).Insert()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.LastInsertId()
    gtest.Assert(n, 1)

    result, err = db.Table("user").Filter().Data(map[interface{}]interface{} {
        "id"          : "2",
        "uid"         : "2",
        "passport"    : "t2",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T2",
        "create_time" : gtime.Now().String(),
    }).Insert()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ = result.RowsAffected()
    gtest.Assert(n, 1)

    type User struct {
        Id         int    `gconv:"id"`
        Uid        int    `gconv:"uid"`
        Passport   string `json:"passport"`
        Password   string `gconv:"password"`
        Nickname   string `gconv:"nickname"`
        CreateTime string `json:"create_time"`
    }
    result, err = db.Table("user").Filter().Data(User{
        Id          : 3,
        Uid         : 3,
        Passport    : "t3",
        Password    : "25d55ad283aa400af464c76d713c07ad",
        Nickname    : "T3",
        CreateTime  : gtime.Now().String(),
    }).Insert()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ = result.RowsAffected()
    gtest.Assert(n, 1)
    value, err := db.Table("user").Fields("passport").Where("id=3").Value()
    gtest.Assert(err, nil)
    gtest.Assert(value.String(), "t3")

    result, err = db.Table("user").Filter().Data(&User{
        Id          : 4,
        Uid         : 4,
        Passport    : "t4",
        Password    : "25d55ad283aa400af464c76d713c07ad",
        Nickname    : "T4",
        CreateTime  : gtime.Now().String(),
    }).Insert()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ = result.RowsAffected()
    gtest.Assert(n, 1)
    value, err = db.Table("user").Fields("passport").Where("id=4").Value()
    gtest.Assert(err, nil)
    gtest.Assert(value.String(), "t4")

    result, err = db.Table("user").Where("id>?", 1).Delete()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ = result.RowsAffected()
    gtest.Assert(n, 3)
}

func TestModel_Batch(t *testing.T) {
    // batch insert
    gtest.Case(t, func() {
        result, err := db.Table("user").Filter().Data(g.List{
            {
                "id"          : 2,
                "uid"         : 2,
                "passport"    : "t2",
                "password"    : "25d55ad283aa400af464c76d713c07ad",
                "nickname"    : "T2",
                "create_time" : gtime.Now().String(),
            },
            {
                "id"          : 3,
                "uid"         : 3,
                "passport"    : "t3",
                "password"    : "25d55ad283aa400af464c76d713c07ad",
                "nickname"    : "T3",
                "create_time" : gtime.Now().String(),
            },
        }).Batch(1).Insert()
        if err != nil {
            gtest.Fatal(err)
        }
        n, _ := result.RowsAffected()
        gtest.Assert(n, 2)
    })

    // batch save
    gtest.Case(t, func() {
        table := createInitTable()
        defer dropTable(table)
        result, err  := db.Table(table).All()
        gtest.Assert(err,         nil)
        gtest.Assert(len(result), INIT_DATA_SIZE)
        for _, v := range result {
            v["nickname"].Set(v["nickname"].String() + v["id"].String())
        }
        r, e := db.Table(table).Data(result).Save()
        gtest.Assert(e, nil)
        n, e := r.RowsAffected()
        gtest.Assert(e, nil)
        gtest.Assert(n, INIT_DATA_SIZE*2)
    })

    // batch replace
    gtest.Case(t, func() {
        table := createInitTable()
        defer dropTable(table)
        result, err  := db.Table(table).All()
        gtest.Assert(err,         nil)
        gtest.Assert(len(result), INIT_DATA_SIZE)
        for _, v := range result {
            v["nickname"].Set(v["nickname"].String() + v["id"].String())
        }
        r, e := db.Table(table).Data(result).Replace()
        gtest.Assert(e, nil)
        n, e := r.RowsAffected()
        gtest.Assert(e, nil)
        gtest.Assert(n, INIT_DATA_SIZE*2)
    })
}

func TestModel_Replace(t *testing.T) {
    result, err := db.Table("user").Data(g.Map{
        "id"          : 1,
        "passport"    : "t11",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T11",
        "create_time" : "2018-10-10 00:01:10",
    }).Replace()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.RowsAffected()
    gtest.Assert(n, 2)
}

func TestModel_Save(t *testing.T) {
    result, err := db.Table("user").Data(g.Map{
        "id"          : 1,
        "passport"    : "t111",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T111",
        "create_time" : "2018-10-10 00:01:10",
    }).Save()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.RowsAffected()
    gtest.Assert(n, 2)
}

func TestModel_Update(t *testing.T) {
    gtest.Case(t, func() {
        result, err := db.Table("user").Data("passport", "t22").Where("passport=?", "t2").Update()
        if err != nil {
            gtest.Fatal(err)
        }
        n, _ := result.RowsAffected()
        gtest.Assert(n, 1)
    })

    gtest.Case(t, func() {
        result, err := db.Table("user").Data("passport", "t2").Where("passport='t22'").Update()
        if err != nil {
            gtest.Fatal(err)
        }
        n, _ := result.RowsAffected()
        gtest.Assert(n, 1)
    })
}

func TestModel_Clone(t *testing.T) {
    md := db.Table("user").Where("id IN(?)", g.Slice{1,3})
    count, err := md.Count()
    if err != nil {
        gtest.Fatal(err)
    }
    record, err := md.OrderBy("id DESC").One()
    if err != nil {
        gtest.Fatal(err)
    }
    result, err := md.OrderBy("id ASC").All()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(count,                 2)
    gtest.Assert(record["id"].Int(),    3)
    gtest.Assert(len(result),           2)
    gtest.Assert(result[0]["id"].Int(), 1)
    gtest.Assert(result[1]["id"].Int(), 3)
}

func TestModel_Safe(t *testing.T) {
    gtest.Case(t, func() {
        md := db.Table("user").Safe(false).Where("id IN(?)", g.Slice{1,3})
        count, err := md.Count()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(count, 2)
        md.And("id = ?", 1)
        count, err = md.Count()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(count, 1)
    })
    gtest.Case(t, func() {
        md := db.Table("user").Safe(true).Where("id IN(?)", g.Slice{1,3})
        count, err := md.Count()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(count, 2)
        md.And("id = ?", 1)
        count, err = md.Count()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(count, 2)
    })
}

func TestModel_All(t *testing.T) {
    result, err := db.Table("user").All()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(len(result), 3)
}

func TestModel_One(t *testing.T) {
    record, err := db.Table("user").Where("id", 1).One()
    if err != nil {
        gtest.Fatal(err)
    }
    if record == nil {
        gtest.Fatal("FAIL")
    }
    gtest.Assert(record["nickname"].String(), "T111")
}

func TestModel_Value(t *testing.T) {
    value, err := db.Table("user").Fields("nickname").Where("id", 1).Value()
    if err != nil {
        gtest.Fatal(err)
    }
    if value == nil {
        gtest.Fatal("FAIL")
    }
    gtest.Assert(value.String(), "T111")
}

func TestModel_Count(t *testing.T) {
    count, err := db.Table("user").Count()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(count, 3)
}

func TestModel_Select(t *testing.T) {
    result, err := db.Table("user").Select()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(len(result), 3)
}

func TestModel_Struct(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime gtime.Time
        }
        user := new(User)
        err := db.Table("user").Where("id=1").Struct(user)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(user.NickName,            "T111")
        gtest.Assert(user.CreateTime.String(), "2018-10-10 00:01:10")
    })
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime *gtime.Time
        }
        user := new(User)
        err := db.Table("user").Where("id=1").Struct(user)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(user.NickName,            "T111")
        gtest.Assert(user.CreateTime.String(), "2018-10-10 00:01:10")
    })
}

func TestModel_Structs(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime gtime.Time
        }
        var users []User
        err := db.Table("user").OrderBy("id asc").Structs(&users)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(users),  3)
        gtest.Assert(users[0].Id, 1)
        gtest.Assert(users[1].Id, 2)
        gtest.Assert(users[2].Id, 3)
        gtest.Assert(users[0].NickName,            "T111")
        gtest.Assert(users[1].NickName,            "T2")
        gtest.Assert(users[2].NickName,            "T3")
        gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
    })
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime *gtime.Time
        }
        var users []*User
        err := db.Table("user").OrderBy("id asc").Structs(&users)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(users),  3)
        gtest.Assert(users[0].Id, 1)
        gtest.Assert(users[1].Id, 2)
        gtest.Assert(users[2].Id, 3)
        gtest.Assert(users[0].NickName,            "T111")
        gtest.Assert(users[1].NickName,            "T2")
        gtest.Assert(users[2].NickName,            "T3")
        gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
    })
}

func TestModel_Scan(t *testing.T) {
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime gtime.Time
        }
        user := new(User)
        err := db.Table("user").Where("id=1").Scan(user)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(user.NickName,            "T111")
        gtest.Assert(user.CreateTime.String(), "2018-10-10 00:01:10")
    })
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime *gtime.Time
        }
        user := new(User)
        err := db.Table("user").Where("id=1").Scan(user)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(user.NickName,            "T111")
        gtest.Assert(user.CreateTime.String(), "2018-10-10 00:01:10")
    })
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime gtime.Time
        }
        var users []User
        err := db.Table("user").OrderBy("id asc").Scan(&users)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(users),  3)
        gtest.Assert(users[0].Id, 1)
        gtest.Assert(users[1].Id, 2)
        gtest.Assert(users[2].Id, 3)
        gtest.Assert(users[0].NickName,            "T111")
        gtest.Assert(users[1].NickName,            "T2")
        gtest.Assert(users[2].NickName,            "T3")
        gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
    })
    gtest.Case(t, func() {
        type User struct {
            Id         int
            Passport   string
            Password   string
            NickName   string
            CreateTime *gtime.Time
        }
        var users []*User
        err := db.Table("user").OrderBy("id asc").Scan(&users)
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(users),  3)
        gtest.Assert(users[0].Id, 1)
        gtest.Assert(users[1].Id, 2)
        gtest.Assert(users[2].Id, 3)
        gtest.Assert(users[0].NickName,            "T111")
        gtest.Assert(users[1].NickName,            "T2")
        gtest.Assert(users[2].NickName,            "T3")
        gtest.Assert(users[0].CreateTime.String(), "2018-10-10 00:01:10")
    })
}

func TestModel_OrderBy(t *testing.T) {
    result, err := db.Table("user").OrderBy("id DESC").Select()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(len(result), 3)
    gtest.Assert(result[0]["nickname"].String(), "T3")
}

func TestModel_GroupBy(t *testing.T) {
    result, err := db.Table("user").GroupBy("id").Select()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(len(result), 3)
    gtest.Assert(result[0]["nickname"].String(), "T111")
}

func TestModel_Where(t *testing.T) {
    // string
    gtest.Case(t, func() {
       result, err := db.Table("user").Where("id=? and nickname=?", 3, "T3").One()
       if err != nil {
           gtest.Fatal(err)
       }
       gtest.AssertGT(len(result),      0)
       gtest.Assert(result["id"].Int(), 3)
    })
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("id", 3).One()
        if err != nil {
           gtest.Fatal(err)
        }
        gtest.AssertGT(len(result),      0)
        gtest.Assert(result["id"].Int(), 3)
    })
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("id", 3).Where("nickname", "T3").One()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(result["id"].Int(), 3)
    })
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("id", 3).And("nickname", "T3").One()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(result["id"].Int(), 3)
    })
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("id", 30).Or("nickname", "T3").One()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(result["id"].Int(), 3)
    })
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("id", 30).Or("nickname", "T3").And("id>?", 1).One()
        gtest.Assert(err,                nil)
        gtest.Assert(result["id"].Int(), 3)
    })
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("id", 30).Or("nickname", "T3").And("id>", 1).One()
        gtest.Assert(err,                nil)
        gtest.Assert(result["id"].Int(), 3)
    })
    // map
    gtest.Case(t, func() {
        result, err := db.Table("user").Where(g.Map{"id" : 3, "nickname" : "T3"}).One()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(result["id"].Int(), 3)
    })
    // map key operator
    gtest.Case(t, func() {
        result, err := db.Table("user").Where(g.Map{"id>" : 1, "id<" : 3}).One()
        gtest.Assert(err,                nil)
        gtest.Assert(result["id"].Int(), 2)
    })
    // complicated where 1
    gtest.Case(t, func() {
        //db.SetDebug(true)
        conditions := g.Map{
           "nickname like ?"      : "%T%",
           "id between ? and ?"   : g.Slice{1,3},
           "id > 0"               : nil,
           "create_time > 0"      : nil,
           "id"                   : g.Slice{1,2,3},
        }
        result, err := db.Table("user").Where(conditions).OrderBy("id asc").All()
        gtest.Assert(err,                   nil)
        gtest.Assert(len(result),           3)
        gtest.Assert(result[0]["id"].Int(), 1)
    })
    // complicated where 2
    gtest.Case(t, func() {
        //db.SetDebug(true)
        conditions := g.Map{
            "nickname like ?"      : "%T%",
            "id between ? and ?"   : g.Slice{1,3},
            "id >= ?"              : 1,
            "create_time > ?"      : 0,
            "id in(?)"             : g.Slice{1,2,3},
        }
        result, err := db.Table("user").Where(conditions).OrderBy("id asc").All()
        gtest.Assert(err,                   nil)
        gtest.Assert(len(result),           3)
        gtest.Assert(result[0]["id"].Int(), 1)
    })
    // struct
    gtest.Case(t, func() {
        type User struct {
            Id       int    `json:"id"`
            Nickname string `gconv:"nickname"`
        }
        result, err := db.Table("user").Where(User{3, "T3"}).One()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(result["id"].Int(), 3)

        result, err  = db.Table("user").Where(&User{3, "T3"}).One()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(result["id"].Int(), 3)
    })
    // slice single
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("id IN(?)", g.Slice{1, 3}).OrderBy("id ASC").All()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(result), 2)
        gtest.Assert(result[0]["id"].Int(), 1)
        gtest.Assert(result[1]["id"].Int(), 3)
    })
    // slice + string
    gtest.Case(t, func() {
        result, err := db.Table("user").Where("nickname=? AND id IN(?)", "T3", g.Slice{1,3}).OrderBy("id ASC").All()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(result), 1)
        gtest.Assert(result[0]["id"].Int(), 3)
    })
    // slice + map
    gtest.Case(t, func() {
        result, err := db.Table("user").Where(g.Map{
            "id"       : g.Slice{1,3},
            "nickname" : "T3",
        }).OrderBy("id ASC").All()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(result), 1)
        gtest.Assert(result[0]["id"].Int(), 3)
    })
    // slice + struct
    gtest.Case(t, func() {
        type User struct {
            Ids      []int  `json:"id"`
            Nickname string `gconv:"nickname"`
        }
        result, err := db.Table("user").Where(User{
            Ids      : []int{1, 3},
            Nickname : "T3",
        }).OrderBy("id ASC").All()
        if err != nil {
            gtest.Fatal(err)
        }
        gtest.Assert(len(result), 1)
        gtest.Assert(result[0]["id"].Int(), 3)
    })
}

func TestModel_Delete(t *testing.T) {
    result, err := db.Table("user").Delete()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.RowsAffected()
    gtest.Assert(n, 3)
}


