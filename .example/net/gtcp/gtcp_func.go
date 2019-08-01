package main

import (
	"fmt"
	"os"

	"github.com/gogf/gf/net/gtcp"
)

func main() {
	data, err := gtcp.SendRecv("www.baidu.com:80", []byte("HEAD / HTTP/1.1\n\n"), -1)
	if len(data) > 0 {
		fmt.Println(string(data))
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
	}
}
