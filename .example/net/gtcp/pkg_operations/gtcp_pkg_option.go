package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/net/gtcp"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/util/gconv"
)

func main() {
	// Server
	go gtcp.NewServer("127.0.0.1:8999", func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.RecvPkg(gtcp.PkgOption{MaxSize: 1})
			if err != nil {
				fmt.Println(err)
				break
			}
			fmt.Println("RecvPkg:", string(data))
		}
	}).Run()

	time.Sleep(time.Second)

	// Client
	conn, err := gtcp.NewConn("127.0.0.1:8999")
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for i := 0; i < 10000; i++ {
		if err := conn.SendPkg([]byte(gconv.String(i))); err != nil {
			glog.Error(err)
		}
		time.Sleep(1 * time.Second)
	}
}
