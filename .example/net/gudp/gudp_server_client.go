package main

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/net/gudp"
	"github.com/gogf/gf/v2/os/glog"
	"github.com/gogf/gf/v2/os/gtime"
)

func main() {
	// Server
	go gudp.NewServer("127.0.0.1:8999", func(conn *gudp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				if err := conn.Send(append([]byte("> "), data...)); err != nil {
					glog.Error(err)
				}
			}
			if err != nil {
				glog.Error(err)
			}
		}
	}).Run()

	time.Sleep(time.Second)

	// Client
	for {
		if conn, err := gudp.NewConn("127.0.0.1:8999"); err == nil {
			if b, err := conn.SendRecv([]byte(gtime.Datetime()), -1); err == nil {
				fmt.Println(string(b), conn.LocalAddr(), conn.RemoteAddr())
			} else {
				glog.Error(err)
			}
			conn.Close()
		} else {
			glog.Error(err)
		}
		time.Sleep(time.Second)
	}
}
