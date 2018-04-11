package main

import (
    "strings"
    "fmt"
    "gitee.com/johng/gf/g/util/gregx"
    "strconv"
)

func ruleToRegx(rule string) (regrule string, querystr string) {
    regrule = "/"
    array  := strings.Split(rule[1:], "/")
    index  := 1
    for _, v := range array {
        switch v[0] {
            case ':':
                regrule += `/([\w\.\-]+)`
                if len(querystr) > 0 {
                    querystr += "&"
                }
                querystr += v[1:] + "=$" + strconv.Itoa(index)
                index++
            case '*':
                regrule += `/(.*)`
                if len(querystr) > 0 {
                    querystr += "&"
                }
                querystr += v[1:] + "=$" + strconv.Itoa(index)
                return
            default:
                regrule += v
        }
    }
    return
}

func main() {
    fmt.Println(strings.Compare("A", "A"))
    fmt.Println(gregx.MatchString(`(\w+):(.+)@([\w\.\-]+)`, "get/users/name/*action@www.johng-cn.com"))
    fmt.Println(ruleToRegx("/users/:name/*action"))
    fmt.Println(gregx.IsMatch(`/users/([\w\.\-]+)/(.*)`, []byte("/users/john/aa/aa/aa")))
}