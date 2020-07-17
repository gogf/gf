package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/net/gtcp"
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
