package main

import (
	"fmt"
	"github.com/gogf/gf/os/gfile"
)

func main() {
	s := `/Users/john/Workspace/Go/GOPATH/pkg/mod/github.com/nats-io/nats-server/v2@v2.1.4`
	d := `/Users/john/Workspace/Go/GOPATH/src/github.com/nats-io/nats-server/v2`
	fmt.Println(gfile.Copy(s, d))
}
