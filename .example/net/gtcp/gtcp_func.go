package main

import (
	"fmt"
	"os"

	"github.com/gogf/gf/v2/net/gtcp"
)

func main() {
	dstConn, err := gtcp.NewPoolConn("www.medlinker.com:80")
	_, err = dstConn.Write([]byte("HEAD / HTTP/1.1\n\n"))
	if err != nil {
		fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
	}
	fmt.Println(dstConn.RecvLine())
}
