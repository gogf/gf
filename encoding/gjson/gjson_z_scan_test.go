package gjson_test

import (
	"database/sql"
	"fmt"
	//_ "github.com/go-sql-driver/mysql" //测试时使用库测试后注释掉
	"github.com/gogf/gf/v2/encoding/gjson"
	"testing"
)

func Test_Scanner(t *testing.T) {
	dns := "root:root@tcp(127.0.0.1:3306)/demo" //填写自己的数据库地址
	db, err := sql.Open("mysql", dns)
	if err != nil {
		panic(err)
	}
	scan := gjson.JsonScanner{}
	err = db.QueryRow("select info from person").Scan(&scan) //按照实际情况填写SQL
	if err != nil {
		panic(err)
	}
	fmt.Print(scan.String())
}
