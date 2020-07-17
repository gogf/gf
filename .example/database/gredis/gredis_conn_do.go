package main

import (
	"fmt"

	"github.com/jin502437344/gf/frame/g"
	"github.com/jin502437344/gf/util/gconv"
)

func main() {
	conn := g.Redis().Conn()
	defer conn.Close()
	conn.Do("SET", "k", "v")
	v, _ := conn.Do("GET", "k")
	fmt.Println(gconv.String(v))
}
