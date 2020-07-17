package main

import (
	"fmt"

	"github.com/jin502437344/gf/net/gtcp"
)

func main() {
	gtcp.NewServer("127.0.0.1:8999", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					fmt.Println(err)
				}
			}
			if err != nil {
				break
			}
		}
	}).Run()
}
