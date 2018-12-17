package gdb_test

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/database/gdb"
    "gitee.com/johng/gf/g/os/gtime"
    "gitee.com/johng/gf/g/util/gtest"
    "testing"
)

var (
	// 数据库对象/接口
	db gdb.DB
)

// 初始化连接参数。
// 测试前需要修改连接参数。
func init() {
	gdb.AddDefaultConfigNode(gdb.ConfigNode{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Pass:     "",
		Name:     "test",
		Type:     "mysql",
		Role:     "master",
		Charset:  "utf8",
        Priority: 1,
	})
	if r, err := gdb.New(); err != nil {
        gtest.Fatal(err)
	} else {
		db = r
	}
	// 准备测试数据结构
    if _, err := db.Exec("DROP TABLE `user`"); err != nil {
        gtest.Fatal(err)
    }
    if _, err := db.Exec(`
    CREATE TABLE user (
        id int(10) unsigned NOT NULL AUTO_INCREMENT COMMENT '用户ID',
        passport varchar(45) NOT NULL COMMENT '账号',
        password char(32) NOT NULL COMMENT '密码',
        nickname varchar(45) NOT NULL COMMENT '昵称',
        create_time timestamp NOT NULL COMMENT '创建时间/注册时间',
        PRIMARY KEY (id)
    ) ENGINE=InnoDB DEFAULT CHARSET=utf8;
    `); err != nil {
        gtest.Fatal(err)
    }
}

func TestDbBase_Query(t *testing.T) {
    if _, err := db.Query("SELECT ?", 1); err != nil {
        gtest.Fatal(err)
    }
    if _, err := db.Query("ERROR"); err == nil {
        gtest.Fatal("FAIL")
    }
}

func TestDbBase_Exec(t *testing.T) {
	if _, err := db.Exec("SELECT ?", 1); err != nil {
		gtest.Fatal(err)
	}
	if _, err := db.Exec("ERROR"); err == nil {
		gtest.Fatal("FAIL")
	}
}

func TestDbBase_Prepare(t *testing.T) {
    st, err := db.Prepare("SELECT 100")
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
}

func TestDbBase_Insert(t *testing.T) {
    if _, err := db.Insert("user", g.Map{
        "id"          : 1,
        "passport"    : "t1",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T1",
        "create_time" : gtime.Now().String(),
    }); err != nil {
        gtest.Fatal(err)
    }
}

func TestDbBase_BatchInsert(t *testing.T) {
    if _, err := db.BatchInsert("user", g.List {
        {
            "id"          : 2,
            "passport"    : "t2",
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
}

func TestDbBase_Save(t *testing.T) {
    if _, err := db.Save("user", g.Map{
        "id"          : 1,
        "passport"    : "t1",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T11",
        "create_time" : gtime.Now().String(),
    }); err != nil {
        gtest.Fatal(err)
    }
}

func TestDbBase_Replace(t *testing.T) {
    if _, err := db.Save("user", g.Map{
        "id"          : 1,
        "passport"    : "t1",
        "password"    : "25d55ad283aa400af464c76d713c07ad",
        "nickname"    : "T111",
        "create_time" : gtime.Now().String(),
    }); err != nil {
        gtest.Fatal(err)
    }
}

func TestDbBase_Update(t *testing.T) {
    if result, err := db.Update("user", "create_time='2010-10-10 00:00:01'", "id=3"); err != nil {
        gtest.Fatal(err)
    } else {
        n, _ := result.RowsAffected()
        gtest.Assert(n, 1)
    }
}

func TestDbBase_GetAll(t *testing.T) {
    if result, err := db.GetAll("SELECT * FROM user WHERE id=?", 1); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(len(result), 1)
    }
}

func TestDbBase_GetOne(t *testing.T) {
    if record, err := db.GetOne("SELECT * FROM user WHERE passport=?", "t1"); err != nil {
        gtest.Fatal(err)
    } else {
        if record == nil {
            gtest.Fatal("FAIL")
        }
        gtest.Assert(record["nickname"].String(), "T111")
    }
}

func TestDbBase_GetValue(t *testing.T) {
    if value, err := db.GetValue("SELECT id FROM user WHERE passport=?", "t3"); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(value.Int(), 3)
    }
}

func TestDbBase_GetCount(t *testing.T) {
    if count, err := db.GetCount("SELECT * FROM user"); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(count, 3)
    }
}

func TestDbBase_GetStruct(t *testing.T) {
    type User struct {
        Id         int
        Passport   string
        Password   string
        NickName   string
        CreateTime gtime.Time
    }
    user := new(User)
    if err := db.GetStruct(user, "SELECT * FROM user WHERE id=?", 3); err != nil {
        gtest.Fatal(err)
    } else {
        gtest.Assert(user.CreateTime.String(), "2010-10-10 00:00:01")
    }
}

func TestDbBase_Delete(t *testing.T) {
    if result, err := db.Delete("user", nil); err != nil {
        gtest.Fatal(err)
    } else {
        n, _ := result.RowsAffected()
        gtest.Assert(n, 3)
    }
}

