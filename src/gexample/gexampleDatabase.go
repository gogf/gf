package main

import (
    "fmt"
    "gf"
    "g"
)

func main() {
    db := g.Db.New(g.GstDbConfig{
        Host : "192.168.2.124",
        Port : "3306",
        User : "root",
        Pass : "123456",
        Name : "hhzl_gdg",
    })
    list, err := db.GetAll("SELECT * FROM henghe_geoinfo_province;")
    if err != nil {
        panic(err.Error())
    }

    for i := range list {
        fmt.Printf("index:%d\n", i)
        for k, v := range list[i] {
            fmt.Printf("%s:%s\n", k, v)
        }
    }

    fmt.Println(db.GetOne("SELECT * FROM henghe_geoinfo_province;"))
    fmt.Println(db.LastInsertId())
}