package main

import (
    "gitee.com/johng/gf/g/util/gregex"
    "fmt"
    "gitee.com/johng/gf/g/util/gutil"
)



func main() {
    s := `username    @   required|length:6,30  #  请输入用户名称|用户名称长度非法`
    match, err := gregex.MatchString(`\s*((\w+)\s*@){0,1}\s*([^#]+)\s*(#\s*(.*)){0,1}\s*`, s)
    fmt.Println(err)
    gutil.Dump(match)
}