package main

import (
    "fmt"
    "gitee.com/johng/gf/g/util/gregex"
)



func main() {
    query := "select * from user"
    q, err := gregex.ReplaceString(`(?i)(SELECT)\s+(.+)\s+(FROM)`, `$1 COUNT($2) $3`, query)
    fmt.Println(err)
    fmt.Println(q)
}