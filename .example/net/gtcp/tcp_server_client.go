package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	addr := "127.0.0.1:8999"

	// Server
	go func() {
		tcpaddr, err := net.ResolveTCPAddr("tcp4", addr)
		if err != nil {
			panic(err)
		}
		listen, err := net.ListenTCP("tcp", tcpaddr)
		if err != nil {
			panic(err)
		}
		for {
			if conn, err := listen.Accept(); err != nil {
				panic(err)
			} else if conn != nil {
				go func(conn net.Conn) {
					for {
						buffer := make([]byte, 1024)
						n, err := conn.Read(buffer)
						if err != nil {
							fmt.Println(err)
							break
						} else {
							fmt.Println(">", string(buffer[0:n]))
							conn.Close()
						}
					}

				}(conn)
			}
		}
	}()

	time.Sleep(time.Second)

	// Client
	if conn, err := net.Dial("tcp", addr); err == nil {
		// first write
		_, err := conn.Write([]byte("hello1"))
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		} else {
			fmt.Println("ok")
		}

		// sleep 10 seconds and re-send
		time.Sleep(10 * time.Second)

		// second write
		_, err = conn.Write([]byte("hello2"))
		if err != nil {
			fmt.Println(err)
			conn.Close()
			return
		} else {
			fmt.Println("ok")
		}
		// sleep 10 seconds and re-send
		time.Sleep(10 * time.Second)
	} else {
		panic(err)
	}

}
