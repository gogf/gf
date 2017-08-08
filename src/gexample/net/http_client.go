package main

import (
    "fmt"
    "g/net/ghttp"
)


func main() {
    c := ghttp.NewClient()

     //r := c.Request("post", "http://192.168.2.102:4168/kv", "{\"name3\":\"john3\"}")
    //r := c.Request("delete", "http://192.168.2.102:4168/kv", "[\"name2\"]")
    //r := c.Request("put", "http://192.168.2.102:4168/node", "[\"172.17.42.1\"]")
    r := c.Request("delete", "http://192.168.2.102:4168/node", "[\"172.17.42.1\"]")
    //fmt.Println(r)
    fmt.Println(r.ReadAll())
}
