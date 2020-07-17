package main

import (
	"fmt"
	"time"

	//_ "github.com/denisenkom/go-mssqldb"
	"github.com/jin502437344/gf/database/gdb"
	"github.com/jin502437344/gf/frame/g"
)

// 本文件用于gf框架的mssql数据库操作示例，不作为单元测试使用

var db gdb.DB

// 初始化配置及创建数据库
func init() {
	gdb.AddDefaultConfigNode(gdb.ConfigNode{
		Host:    "127.0.0.1",
		Port:    "1433",
		User:    "sa",
		Pass:    "123456",
		Name:    "test",
		Type:    "mssql",
		Role:    "master",
		Charset: "utf8",
	})
	db, _ = gdb.New()

	//gins.Config().SetPath("/home/john/Workspace/Go/GOPATH/src/github.com/jin502437344/gf/.example/frame")
	//db = g.Database()

	//gdb.SetConfig(gdb.ConfigNode {
	//    Host : "127.0.0.1",
	//    Port : 3306,
	//    User : "root",
	//    Pass : "123456",
	//    Name : "test",
	//    Type : "mysql",
	//})
	//db, _ = gdb.Instance()

	//gdb.SetConfig(gdb.Config {
	//    "default" : gdb.ConfigGroup {
	//        gdb.ConfigNode {
	//            Host     : "127.0.0.1",
	//            Port     : "3306",
	//            User     : "root",
	//            Pass     : "123456",
	//            Name     : "test",
	//            Type     : "mysql",
	//            Role     : "master",
	//            Weight : 100,
	//        },
	//        gdb.ConfigNode {
	//            Host     : "127.0.0.2",
	//            Port     : "3306",
	//            User     : "root",
	//            Pass     : "123456",
	//            Name     : "test",
	//            Type     : "mysql",
	//            Role     : "master",
	//            Weight : 100,
	//        },
	//        gdb.ConfigNode {
	//            Host     : "127.0.0.3",
	//            Port     : "3306",
	//            User     : "root",
	//            Pass     : "123456",
	//            Name     : "test",
	//            Type     : "mysql",
	//            Role     : "master",
	//            Weight : 100,
	//        },
	//        gdb.ConfigNode {
	//            Host     : "127.0.0.4",
	//            Port     : "3306",
	//            User     : "root",
	//            Pass     : "123456",
	//            Name     : "test",
	//            Type     : "mysql",
	//            Role     : "master",
	//            Weight : 100,
	//        },
	//    },
	//})
	//db, _ = gdb.Instance()
}

// 创建测试数据库
func create() error {
	fmt.Println("drop table aa_user:")
	_, err := db.Exec("drop table aa_user")
	if err != nil {
		fmt.Println("drop table aa_user error.", err)
	}

	s := `
        CREATE TABLE aa_user (
            id  int not null,
            name VARCHAR(60),
            age  int,
            addr varchar(60),
            PRIMARY KEY (id)
        )
    `
	fmt.Println("create table aa_user:")
	_, err = db.Exec(s)
	if err != nil {
		fmt.Println("create table error.", err)
		return err
	}

	/*_, err = db.Exec("drop sequence id_seq")
	  if err != nil {
	      fmt.Println("drop sequence id_seq", err)
	  }

	  fmt.Println("create sequence id_seq")
	  _, err = db.Exec("create sequence id_seq increment by 1 start with 1 maxvalue 9999999999 cycle cache 10")
	  if err != nil {
	      fmt.Println("create sequence id_seq error.", err)
	      return err
	  }

	  s = `
	  CREATE TRIGGER id_trigger before insert on aa_user  for each row
	  begin
	  select id_seq.nextval into :new.id from dual;
	  end;
	  `
	  _, err = db.Exec(s)
	  if err != nil {
	      fmt.Println("create trigger error.", err)
	      return err
	  }*/

	_, err = db.Exec("drop table user_detail")
	if err != nil {
		fmt.Println("drop table user_detail", err)
	}

	s = `
        CREATE TABLE user_detail (
            id   int not null,
            site  VARCHAR(255),
            PRIMARY KEY (id)
        )
    `
	fmt.Println("create table user_detail:")
	_, err = db.Exec(s)
	if err != nil {
		fmt.Println("create table user_detail error.", err)
		return err
	}
	fmt.Println("create table success.")
	return nil
}

// 数据写入
func insert(id int) {
	fmt.Println("insert:")
	r, err := db.Insert("aa_user", gdb.Map{
		"id":   id,
		"name": "john",
		"age":  id,
	})
	fmt.Println(r.LastInsertId())
	fmt.Println(r.RowsAffected())
	if err == nil {
		r, err = db.Insert("user_detail", gdb.Map{
			"id":   id,
			"site": "http://johng.cn",
		})
		if err == nil {
			fmt.Printf("id: %d\n", id)
		} else {
			fmt.Println(err)
		}

	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 基本sql查询
func query() {
	fmt.Println("query:")
	list, err := db.GetAll("select * from aa_user where 1=1")
	if err == nil {
		fmt.Println(list)
	} else {
		fmt.Println(err)
	}

	list, err = db.Table("aa_user").OrderBy("id").Limit(0, 5).Select()
	if err == nil {
		fmt.Println(list)
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// replace into
func replace() {
	fmt.Println("replace:")
	r, err := db.Save("aa_user", gdb.Map{
		"id":   1,
		"name": "john",
	})
	if err == nil {
		fmt.Println(r.LastInsertId())
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 数据保存
func save() {
	fmt.Println("save:")
	r, err := db.Save("aa_user", gdb.Map{
		"id":   1,
		"name": "john",
	})
	if err == nil {
		fmt.Println(r.LastInsertId())
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 批量写入
func batchInsert() {
	fmt.Println("batchInsert:")
	_, err := db.BatchInsert("aa_user", gdb.List{
		{"id": 11, "name": "batchInsert_john_1", "age": 11},
		{"id": 12, "name": "batchInsert_john_2", "age": 12},
		{"id": 13, "name": "batchInsert_john_3", "age": 13},
		{"id": 14, "name": "batchInsert_john_4", "age": 14},
	}, 10)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println()
}

// 数据更新
func update1() {
	fmt.Println("update1:")
	r, err := db.Update("aa_user", gdb.Map{"name": "john1", "age": 1}, "id=?", 1)
	if err == nil {
		fmt.Println(r.LastInsertId())
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 数据更新
func update2() {
	fmt.Println("update2:")
	r, err := db.Update("aa_user", gdb.Map{"name": "john6", "age": 6}, "id=?", 2)
	if err == nil {
		fmt.Println(r.LastInsertId())
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 数据更新
func update3() {
	fmt.Println("update3:")
	r, err := db.Update("aa_user", "name=?", "id=?", "john2", 3)
	if err == nil {
		fmt.Println(r.LastInsertId())
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 链式查询操作1
func linkopSelect1() {
	fmt.Println("linkopSelect1:")
	r, err := db.Table("aa_user u").LeftJoin("user_detail ud", "u.id=ud.id").Fields("u.*, ud.site").Where("u.id > ?", 1).Limit(3, 5).Select()
	if err == nil {
		fmt.Println(r)
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 链式查询操作2
func linkopSelect2() {
	fmt.Println("linkopSelect2:")
	r, err := db.Table("aa_user u").LeftJoin("user_detail ud", "u.id=ud.id").Fields("u.*,ud.site").Where("u.id=?", 1).One()
	if err == nil {
		fmt.Println(r)
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 链式查询操作3
func linkopSelect3() {
	fmt.Println("linkopSelect3:")
	r, err := db.Table("aa_user u").LeftJoin("user_detail ud", "u.id=ud.id").Fields("ud.site").Where("u.id=?", 1).Value()
	if err == nil {
		fmt.Println(r.String())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 链式查询数量1
func linkopCount1() {
	fmt.Println("linkopCount1:")
	r, err := db.Table("aa_user u").LeftJoin("user_detail ud", "u.id=ud.id").Where("name like ?", "john").Count()
	if err == nil {
		fmt.Println(r)
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 错误操作
func linkopUpdate1() {
	fmt.Println("linkopUpdate1:")
	r, err := db.Table("henghe_setting").Update()
	if err == nil {
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println("error", err)
	}
	fmt.Println()
}

// 通过Map指针方式传参方式
func linkopUpdate2() {
	fmt.Println("linkopUpdate2:")
	r, err := db.Table("aa_user").Data(gdb.Map{"name": "john2"}).Where("name=?", "john").Update()
	if err == nil {
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 通过字符串方式传参
func linkopUpdate3() {
	fmt.Println("linkopUpdate3:")
	r, err := db.Table("aa_user").Data("name='john3'").Where("name=?", "john2").Update()
	if err == nil {
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// Where条件使用Map
func linkopUpdate4() {
	fmt.Println("linkopUpdate4:")
	r, err := db.Table("aa_user").Data(gdb.Map{"name": "john11111"}).Where(g.Map{"id": 1}).Update()
	if err == nil {
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 链式批量写入
func linkopBatchInsert1() {
	fmt.Println("linkopBatchInsert1:")
	r, err := db.Table("aa_user").Filter().Data(gdb.List{
		{"id": 21, "name": "linkopBatchInsert1_john_1", "amt": 21.21, "tt": "haha"},
		{"id": 22, "name": "linkopBatchInsert1_john_2", "amt": 22.22, "cc": "hahacc"},
		{"id": 23, "name": "linkopBatchInsert1_john_3", "amt": 23.23, "bb": "hahabb"},
		{"id": 24, "name": "linkopBatchInsert1_john_4", "amt": 24.24, "aa": "hahaaa"},
	}).Insert()
	if err == nil {
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 链式批量写入，指定每批次写入的条数
func linkopBatchInsert2() {
	fmt.Println("linkopBatchInsert2:")
	r, err := db.Table("aa_user").Data(gdb.List{
		{"id": 25, "name": "linkopBatchInsert2john_1"},
		{"id": 26, "name": "linkopBatchInsert2john_2"},
		{"id": 27, "name": "linkopBatchInsert2john_3"},
		{"id": 28, "name": "linkopBatchInsert2john_4"},
	}).Batch(2).Insert()
	if err == nil {
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 链式批量保存
func linkopBatchSave() {
	fmt.Println("linkopBatchSave:")
	r, err := db.Table("aa_user").Data(gdb.List{
		{"id": 1, "name": "john_1"},
		{"id": 2, "name": "john_2"},
		{"id": 3, "name": "john_3"},
		{"id": 4, "name": "john_4"},
	}).Save()
	if err == nil {
		fmt.Println(r.RowsAffected())
	} else {
		fmt.Println(err)
	}
	fmt.Println()
}

// 事务操作示例1
func transaction1() {
	fmt.Println("transaction1:")
	if tx, err := db.Begin(); err == nil {
		r, err := tx.Insert("aa_user", gdb.Map{
			"id":   30,
			"name": "transaction1",
		})
		tx.Rollback()
		fmt.Println(r, err)
	}
	fmt.Println()
}

// 事务操作示例2
func transaction2() {
	fmt.Println("transaction2:")
	if tx, err := db.Begin(); err == nil {
		r, err := tx.Table("user_detail").Data(gdb.Map{"id": 6, "site": "www.baidu.com哈哈哈*?''\"~!@#$%^&*()"}).Insert()
		tx.Commit()
		fmt.Println(r, err)
	}
	fmt.Println()
}

// 主从io复用测试，在mysql中使用 show full processlist 查看链接信息
func keepPing() {
	fmt.Println("keepPing:")
	for i := 0; i < 30; i++ {
		fmt.Println("ping...", i)
		err := db.PingMaster()
		if err != nil {
			fmt.Println(err)
			return
		}
		err = db.PingSlave()
		if err != nil {
			fmt.Println(err)
			return
		}
		time.Sleep(1 * time.Second)
	}
}

// like语句查询
func likeQuery() {
	fmt.Println("likeQuery:")
	if r, err := db.Table("aa_user").Where("name like ?", "%john%").Select(); err == nil {
		fmt.Println(r)
	} else {
		fmt.Println(err)
	}
}

// mapToStruct
func mapToStruct() {
	type User struct {
		Id   int
		Name string
		Age  int
		Addr string
	}
	fmt.Println("mapToStruct:")
	if r, err := db.Table("aa_user").Where("id=?", 1).One(); err == nil {
		u := User{}
		if err := r.ToStruct(&u); err == nil {
			fmt.Println(r)
			fmt.Println(u)
		} else {
			fmt.Println(err)
		}
	} else {
		fmt.Println(err)
	}
}

// getQueriedSqls
func getQueriedSqls() {
	for k, v := range db.GetQueriedSqls() {
		fmt.Println(k, ":")
		fmt.Println("Sql  :", v.Sql)
		fmt.Println("Args :", v.Args)
		fmt.Println("Error:", v.Error)
	}
}

func main() {

	db.PingMaster()
	db.SetDebug(true)
	/*err := create()
	  if err != nil {
	      return
	  }*/

	//test1
	/*for i := 1; i < 5; i++ {
	    insert(i)
	}*/
	//insert(2)
	//query()

	//batchInsert()
	//query()

	//replace()
	//save()

	/*update1()
	  update2()
	  update3()
	*/

	/*linkopSelect1()
	  linkopSelect2()
	  linkopSelect3()
	  linkopCount1()
	*/

	/*linkopUpdate1()
	  linkopUpdate2()
	  linkopUpdate3()
	  linkopUpdate4()
	*/

	linkopBatchInsert1()
	query()
	//linkopBatchInsert2()

	//transaction1()
	//transaction2()
	//
	//keepPing()
	//likeQuery()
	//mapToStruct()
	//getQueriedSqls()
}
