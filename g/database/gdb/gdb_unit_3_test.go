package gdb_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
)

func TestTX_Query(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if rows, err := tx.Query("SELECT ?", 1); err != nil {
        gtest.Fatal(err)
    } else {
        rows.Close()
    }
    if _, err := tx.Query("ERROR"); err == nil {
        gtest.Fatal("FAIL")
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_Exec(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.Exec("SELECT ?", 1); err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.Exec("ERROR"); err == nil {
        gtest.Fatal("FAIL")
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_Commit(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_Rollback(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Rollback(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_Prepare(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    st, err := tx.Prepare("SELECT 100")
    if err != nil {
        gtest.Fatal(err)
    }
    rows, err := st.Query()
    if err != nil {
        gtest.Fatal(err)
    }
    array, err := rows.Columns()
    if err != nil {
        gtest.Fatal(err)
    }
    gtest.Assert(array[0], "100")
    if err := rows.Close(); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_Insert(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.Insert("user", g.Map {
        "id"          : 1,
        "passport"    : "t1",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T1",
        "create_time" : gtime.Now().String(),
    }); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
    if n, err := db.Table("user").Count(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(n, 1)
    }
}

func TestTX_BatchInsert(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.BatchInsert("user", g.List {
        {
            "id"          : 2,
            "passport"    : "t",
            "password"    : "25d55ad283aa400af464c76d713c07ad",
            "nickname"    : "T2",
            "create_time" : gtime.Now().String(),
        },
        {
            "id"          : 3,
            "passport"    : "t3",
            "password"    : "25d55ad283aa400af464c76d713c07ad",
            "nickname"    : "T3",
            "create_time" : gtime.Now().String(),
        },
    }, 10); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
    if n, err := db.Table("user").Count(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(n, 3)
    }
}

func TestTX_BatchReplace(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.BatchReplace("user", g.List {
        {
            "id"          : 2,
            "passport"    : "t2",
            "password"    : "p2",
            "nickname"    : "T2",
            "create_time" : gtime.Now().String(),
        },
        {
            "id"          : 4,
            "passport"    : "t4",
            "password"    : "25d55ad283aa400af464c76d713c07ad",
            "nickname"    : "T4",
            "create_time" : gtime.Now().String(),
        },
    }, 10); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
    // 数据数量
    if n, err := db.Table("user").Count(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(n, 4)
    }
    // 检查replace后的数值
    if value, err := db.Table("user").Fields("password").Where("id", 2).Value(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(value.String(), "p2")
    }
}

func TestTX_BatchSave(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.BatchSave("user", g.List {
        {
            "id"          : 4,
            "passport"    : "t4",
            "password"    : "p4",
            "nickname"    : "T4",
            "create_time" : gtime.Now().String(),
        },
    }, 10); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
    // 数据数量
    if n, err := db.Table("user").Count(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(n, 4)
    }
    // 检查replace后的数值
    if value, err := db.Table("user").Fields("password").Where("id", 4).Value(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(value.String(), "p4")
    }
}

func TestTX_Replace(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.Replace("user", g.Map {
        "id"          : 1,
        "passport"    : "t11",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T11",
        "create_time" : gtime.Now().String(),
    }); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Rollback(); err != nil {
        gtest.Fatal(err)
    }
    if value, err := db.Table("user").Fields("nickname").Where("id", 1).Value(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(value.String(), "T1")
    }
}

func TestTX_Save(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.Save("user", g.Map {
        "id"          : 1,
        "passport"    : "t11",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T11",
        "create_time" : gtime.Now().String(),
    }); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
    if value, err := db.Table("user").Fields("nickname").Where("id", 1).Value(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(value.String(), "T11")
    }
}

func TestTX_GetAll(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if result, err := tx.GetAll("SELECT * FROM user WHERE id=?", 1); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(len(result), 1)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_GetOne(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if record, err := tx.GetOne("SELECT * FROM user WHERE passport=?", "t2"); err != nil {
        gtest.Fatal(err)
    } else {
        if record == nil {
            gtest.Fatal("FAIL")
        }
        gtest.Assert(record["nickname"].String(), "T2")
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_GetValue(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if value, err := tx.GetValue("SELECT id FROM user WHERE passport=?", "t3"); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(value.Int(), 3)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_GetCount(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if count, err := tx.GetCount("SELECT * FROM user"); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(count, 4)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_GetStruct(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    type User struct {
        Id         int
        Passport   string
        Password   string
        NickName   string
        CreateTime gtime.Time
    }
    user := new(User)
    if err := tx.GetStruct(user, "SELECT * FROM user WHERE id=?", 1); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(user.NickName, "T11")
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
}

func TestTX_Delete(t *testing.T) {
    tx, err := db.Begin()
    if err != nil {
        gtest.Fatal(err)
    }
    if _, err := tx.Delete("user", nil); err != nil {
        gtest.Fatal(err)
    }
    if err := tx.Commit(); err != nil {
        gtest.Fatal(err)
    }
    if n, err := db.Table("user").Count(); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(n, 0)
    }
}


