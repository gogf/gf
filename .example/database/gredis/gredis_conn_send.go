package main

import (
	"fmt"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
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
	v, _ := conn.Receive()
	fmt.Println(gconv.String(v))
}
