package main

import (
    "fmt"
    "gitee.com/johng/gf/g/net/gtcp"
)

func main() {
    data, err := gtcp.SendReceive("www.baidu.com:80", []byte("GET / HTTP/1.1\n\n"), -1)
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
}