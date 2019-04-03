package main

import (
	"fmt"
	"github.com/gogf/gf/g/net/gtcp"
	"time"
)

func main() {
	gtcp.NewServer("127.0.0.1:8999", func(conn *gtcp.Conn) {
		defer conn.Close()
		conn.SetRecvDeadline(time.Now().Add(10 * time.Second))
		for {
			data, err := conn.Recv(-1)
			fmt.Println(err)
			if len(data) > 0 {
				fmt.Println(string(data))
			}
			if err != nil {
				break
			}
		}
	}).Run()
}
