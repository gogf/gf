package main

import (
    "gitee.com/johng/gf/g"
    "gitee.com/johng/gf/g/net/ghttp"
)

func main() {
    s := g.Server()
    s.BindHandler("/ws", func(r *ghttp.Request) {
        conn, _ := r.WebSocket()
        for {
            msgType, msg, err := conn.ReadMessage()
            if err != nil {
                return
            }
            if err = conn.WriteMessage(msgType, msg); err != nil {
                return
            }
        }
    })
    s.SetPort(8199)
    s.Run()
}

