package gdb_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
)

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
}

func TestModel_Batch(t *testing.T) {
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
    }).Batch(10).Insert()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.RowsAffected()
    gtest.Assert(n, 2)
}

func TestModel_Replace(t *testing.T) {
    result, err := db.Table("user").Data(g.Map{
        "id"          : 1,
        "passport"    : "t11",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T11",
        "create_time" : gtime.Now().String(),
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
        "create_time" : gtime.Now().String(),
    }).Save()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.RowsAffected()
    gtest.Assert(n, 2)
}

func TestModel_Update(t *testing.T) {
    result, err := db.Table("user").Data("passport", "t22").Where("passport=?", "t2").Update()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.RowsAffected()
    gtest.Assert(n, 1)
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
    gtest.Assert(count, 2)
    gtest.Assert(record["id"].Int(), 3)
    gtest.Assert(len(result), 2)
    gtest.Assert(result[0]["id"].Int(), 1)
    gtest.Assert(result[1]["id"].Int(), 3)
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
    gtest.Assert(user.NickName, "T111")
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

func TestModel_Where1(t *testing.T) {
    result, err := db.Table("user").Where("id IN(?)", g.Slice{1,3}).OrderBy("id ASC").All()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(len(result), 2)
    gtest.Assert(result[0]["id"].Int(), 1)
    gtest.Assert(result[1]["id"].Int(), 3)
}

func TestModel_Where2(t *testing.T) {
    result, err := db.Table("user").Where("nickname=? AND id IN(?)", "T3", g.Slice{1,3}).OrderBy("id ASC").All()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(len(result), 1)
    gtest.Assert(result[0]["id"].Int(), 3)
}

func TestModel_Delete(t *testing.T) {
    result, err := db.Table("user").Delete()
    if err != nil {
        gtest.Fatal(err)
    }
    n, _ := result.RowsAffected()
    gtest.Assert(n, 3)
}


