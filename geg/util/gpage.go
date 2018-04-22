package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gpage"
)

func main() {
    // 基本分页示例
    page1 := gpage.New(100, 10, 1, "http://xxx.xxx.xxx/user/list?page=1&type=10#anchor")
    fmt.Println(page1.GetContent(3))

    // 基于静态链接的分页示例
    page2 := gpage.New(100, 10, 1, "http://xxx.xxx.xxx/user/list/1?type=10#anchor", "/user/list/:page")
    fmt.Println(page2.GetContent(3))
}