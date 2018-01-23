package main

import (
    "fmt"
    "gitee.com/johng/gf/g/encoding/gjson"
    "gitee.com/johng/gf/g/os/glog"
)

func main() {
    j, _ := gjson.Load("/home/john/Workspace/Go/GOPATH/src/gitee.com/johng/gf/geg/frame/config.yml")
    //c, _ := j.ToToml()
    //fmt.Println(j.Get("database.default").([]interface{})[0])
    fmt.Println(j.Get("database.default.0"))
    return
    data :=
        `{
            "users" : {
                "count" : 100
            }
        }`
    j, err := gjson.DecodeToJson([]byte(data))
    if err != nil {
        glog.Error(err)
    } else {
        //j.Set("users.count",  1)
        j.Set("users.list",  []string{"John", "小明"})
        fmt.Println(j.Set("users.list.10",  []string{"John", "小明10"}))
        fmt.Println(j.Set("users.list.9",  []string{"John", "小明9"}))
        j.Set("users",     "a")
        //fmt.Println(j.Get("users.count"))
        //fmt.Println(j.Get("users.count"))
        fmt.Println(j.Get("users.list.10"))
        c, _ := j.ToJson()
        fmt.Println(string(c))
    }
}