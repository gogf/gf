package main

import (
	"fmt"

	"github.com/jin502437344/gf/net/gudp"
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
