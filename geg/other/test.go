package main

import (
    "fmt"
    "github.com/gogf/gf/g/util/gconv"
)

type User struct {
    Id int
}

type RPCResponse struct {
    ID      interface{}  `json:"id,omitempty"`
    JsonRPC string       `json:"jsonrpc"`
    Error   *User        `json:"error,omitempty"`
    Result  interface{}  `json:"result,omitempty"`
}

func main() {
    var rpc RPCResponse
    fmt.Println(rpc.Error)
    fmt.Println(gconv.Map(rpc))
}