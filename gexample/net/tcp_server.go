package main

import (
    "net"
    "fmt"
    "g/net/gtcp"
    "io"
    "log"
    "time"
    "g/util/gutil"
)

func main() {
    gtcp.NewServer(":8999", func(conn net.Conn) {


        try        := 0
        buffersize := 5
        data       := make([]byte, 0)
        for {
            buffer      := make([]byte, buffersize)
            length, err := conn.Read(buffer)
            if err != nil {
                log.Println(err)
                if err != io.EOF {
                    log.Println("node recieve:", err, "try:", try)
                }
                if try > 2 {
                    break;
                }
                try ++
                time.Sleep(100 * time.Millisecond)
            } else {
                if length == buffersize {
                    data = gutil.MergeSlice(data, buffer)
                } else {
                    data = gutil.MergeSlice(data, buffer[0:length])
                    break;
                }
            }
        }
        fmt.Println(string(data))
    }).Run()
    select {

    }
}