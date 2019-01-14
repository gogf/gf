package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gconv"
    "time"
)

func main(){
    type Test struct {
        Date time.Time `json:"date"`
    }
    o := new(Test)
    m := map[string]interface{}{
        "Date" : "",
    }
    gconv.Struct(m, o)
    fmt.Println(o)
}
