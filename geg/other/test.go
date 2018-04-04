package main

import (
    "gitee.com/johng/gf/g/database/gdb"
    "fmt"
)

type Model struct {
    TableName string
}

var Db *gdb.Db

func init() {
    gdb.AddDefaultConfigNode(gdb.ConfigNode {
        Host    : "127.0.0.1",
        Port    : "3306",
        User    : "root",
        Pass    : "123456",
        Name    : "test",
        Type    : "mysql",
        Role    : "master",
        Charset : "utf8",
    })
    var err error
    Db, err = gdb.Instance()
    checkErr(err)
}

func checkErr(err error) {
    if err != nil {
        panic(err)
    }
}

type UserModel struct {
    Model
}

func (u *UserModel) Get() (user gdb.Map){
    user, _ = Db.Table("user").Fields("uid, nickname, email").Where("uid = ?", 15).One()
    return
}

func (u *UserModel) Insert(data gdb.Map) (id int64) {
    ret, _ := Db.Table("user").Data(data).Insert()
    id, _   = ret.LastInsertId()
    return
}

func main() {
    u    := &UserModel{}
    user := u.Get()
    fmt.Println(user)
    u.Insert(gdb.Map{"uid": 100, "name": "jack"})
}