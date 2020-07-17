package main

import (
	"time"

	"github.com/jin502437344/gf/net/gtcp"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/os/gtime"
)

func main() {
	conn, err := gtcp.NewConn("127.0.0.1:8999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := conn.Send([]byte(gtime.Now().String())); err != nil {
		glog.Error(err)
	}

	time.Sleep(time.Minute)
}
