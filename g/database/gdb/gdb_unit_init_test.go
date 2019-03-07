package gdb_test

import (
    "github.com/gogf/gf/g/database/gdb"
    "github.com/gogf/gf/g/test/gtest"
	"os"
)

var (
	// 数据库对象/接口
	db gdb.DB
)

// 初始化连接参数。
// 测试前需要修改连接参数。
func init() {
	node := gdb.ConfigNode{
		Host:     "127.0.0.1",
		Port:     "3306",
		User:     "root",
		Pass:     "",
		Name:     "",
		Type:     "mysql",
		Role:     "master",
		Charset:  "utf8",
		Priority: 1,
	}
	hostname, _ := os.Hostname()
	if hostname == "ijohn" {
		node.Pass = "12345678"
	}
	gdb.AddDefaultConfigNode(node)
	if r, err := gdb.New(); err != nil {
        gtest.Fatal(err)
	} else {
		db = r
	}
	// 准备测试数据结构
    if _, err := db.Exec("CREATE DATABASE IF NOT EXISTS `test` CHARACTER SET UTF8"); err != nil {
        gtest.Fatal(err)
    }
    db.SetSchema("test")
    if _, err := db.Exec("DROP TABLE IF EXISTS `user`"); err != nil {
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
