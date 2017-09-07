package main

import (
    "fmt"
    "net"
    "time"
    "g/util/gutil"
)


// 获得TCP链接
func getConn(ip string, port int) net.Conn {
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", ip, port), 3 * time.Second)
    if err == nil {
        return conn
    }
    return nil
}

// 发送数据
func Send(conn net.Conn, data []byte) error {
    _, err := conn.Write(data)
    if err != nil {
        return err
    } else {
        return nil
    }
}

// 获取数据
func Receive(conn net.Conn) []byte {
    buffersize := 1024
    data       := make([]byte, 0)
    for {
        buffer      := make([]byte, buffersize)
        length, err := conn.Read(buffer)
        if err == nil {
            if length == buffersize {
                data = gutil.MergeSlice(data, buffer)
            } else {
                data = gutil.MergeSlice(data, buffer[0:length])
                break;
            }
        } else {
            break;
        }
    }
    return data
}

func main() {
    conn := getConn("127.0.0.1", 4168)
    for {
        err := Send(conn, []byte("GET /kv HTTP/1.0\r\nHost: 127.0.0.1:4168\r\n\r\n"))
        if err != nil {
            conn = getConn("127.0.0.1", 4168)
        } else {
            Receive(conn)
        }
    }

}