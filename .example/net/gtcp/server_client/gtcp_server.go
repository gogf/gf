package main

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gtcp"
)

func main() {
	// Server
	gtcp.NewServer("127.0.0.1:8999", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				fmt.Println(string(data))
			}
			if err != nil {
				// client closed, err will be: EOF
				fmt.Println(err)
				break
			}
		}
	}).Run()
}
