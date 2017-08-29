package main

import (
    "fmt"
    "g/database/gdb"
)



func main() {
    dbcfg   := gdb.ConfigNode{
        Host    : "192.168.2.102",
        Port    : "3306",
        User    : "root",
        Pass    : "123456",
        Name    : "test",
        Type    : "mysql",
    }
    db, err := gdb.NewByConfigNode(dbcfg)
    if err != nil || db == nil {
        fmt.Println("1")
    } else {
        if err :=db.PingMaster();err != nil {
            fmt.Println(err)
        } else {
            fmt.Println("3")
        }
        db.Close()
    }
}