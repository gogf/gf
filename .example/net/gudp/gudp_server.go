package main

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gudp"
)

func main() {
	gudp.NewServer("127.0.0.1:8999", func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			fmt.Println(err, string(data))
		}
	}).Run()
}
