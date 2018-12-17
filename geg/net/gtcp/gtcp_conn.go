package main

import (
    "fmt"
    "bytes"
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/util/gconv"
    "os"
)

func main() {
    conn, err := gtcp.NewConn("www.baidu.com:80")
    if err != nil {
        panic(err)
    }
    defer conn.Close()

    if err := conn.Send([]byte("GET / HTTP/1.1\n\n")); err != nil {
        panic(err)
    }

    header        := make([]byte, 0)
    content       := make([]byte, 0)
    contentLength := 0
    for {
        data, err := conn.RecvLine()
        // header读取，解析文本长度
        if len(data) > 0 {
            array := bytes.Split(data, []byte(": "))
            // 获得页面内容长度
            if contentLength == 0 && len(array) == 2 && bytes.EqualFold([]byte("Content-Length"), array[0]) {
                contentLength = gconv.Int(array[1])
            }
            header = append(header, data...)
            header = append(header, '\n')
        }
        // header读取完毕，读取文本内容
        if contentLength > 0 && len(data) == 0 {
            content, _ = conn.Recv(contentLength)
            break
        }
        if err != nil {
            fmt.Fprintf(os.Stderr, "ERROR: %s\n", err.Error())
            break
        }
    }

    if len(header) > 0 {
        fmt.Println(string(header))
    }
    if len(content) > 0 {
        fmt.Println(string(content))
    }
}