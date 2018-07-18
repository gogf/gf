package main

import (
    "gitee.com/johng/gf/g/util/gutil"
    "fmt"
)

type User struct {
    Name string
    Age  int
    Adds string
}

func main() {
    m := map[string]interface{} {
        "name" : "john",
        "age"  : 16,
        "adds" : "test",
    }
    o := User{}
    e := gutil.MapToStruct(m, &o)
    fmt.Println(e)
    fmt.Println(o)
}
