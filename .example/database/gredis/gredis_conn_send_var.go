package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
)

func main() {
	conn := g.Redis().Conn()
	defer conn.Close()
	conn.Send("SET", "foo", "bar")
	conn.Send("GET", "foo")
	conn.Flush()
	// reply from SET
	conn.Receive()
	// reply from GET
	v, _ := conn.ReceiveVar()
	fmt.Println(v.String())
}
