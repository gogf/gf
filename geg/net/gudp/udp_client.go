package main

import (
	"fmt"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:8999")
	defer conn.Close()
	if err != nil {
		os.Exit(1)
	}

	conn.Write([]byte("Hello world!"))

	buffer := make([]byte, 100)

	conn.Read(buffer)

	fmt.Println(string(buffer))
}
