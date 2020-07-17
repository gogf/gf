package main

import (
	"fmt"
	"time"

	"github.com/jin502437344/gf/net/gtcp"
	"github.com/jin502437344/gf/os/glog"
	"github.com/jin502437344/gf/util/gconv"
)

func main() {
	address := "127.0.0.1:8999"
	crtFile := "server.crt"
	keyFile := "server.key"
	// TLS Server
	go gtcp.NewServerKeyCrt(address, crtFile, keyFile, func(conn *gtcp.Conn) {
		defer conn.Close()
		for {
			data, err := conn.Recv(-1)
			if len(data) > 0 {
				fmt.Println(string(data))
			}
			if err != nil {
				// if client closes, err will be: EOF
				glog.Error(err)
				break
			}
		}
	}).Run()

	time.Sleep(time.Second)

	// Client
	tlsConfig, err := gtcp.LoadKeyCrt(crtFile, keyFile)
	if err != nil {
		panic(err)
	}
	tlsConfig.InsecureSkipVerify = true

	conn, err := gtcp.NewConnTLS(address, tlsConfig)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	for i := 0; i < 10; i++ {
		if err := conn.Send([]byte(gconv.String(i))); err != nil {
			glog.Error(err)
		}
		time.Sleep(time.Second)
		if i == 5 {
			conn.Close()
			break
		}
	}

	// exit after 5 seconds
	time.Sleep(5 * time.Second)
}
