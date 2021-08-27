package main

import (
	"bytes"
	"fmt"
	"os"

	"github.com/gogf/gf/net/gtcp"
	"github.com/gogf/gf/util/gconv"
)

func main() {
	conn, err := gtcp.NewConn("www.baidu.com:80")
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	if err := conn.Send([]byte("GET / HTTP/1.1\r\n\r\n")); err != nil {
		panic(err)
	}

	header := make([]byte, 0)
	content := make([]byte, 0)
	contentLength := 0
	for {
		data, err := conn.RecvLine()
		// header读取，解析文本长度
		if len(data) > 0 {
			array := bytes.Split(data, []byte(": "))
			// 获得页面内容长度
			if contentLength == 0 && len(array) == 2 && bytes.EqualFold([]byte("Content-Length"), array[0]) {
				// http 以\r\n换行，需要把\r也去掉
				contentLength = gconv.Int(string(array[1][:len(array[1])-1]))
			}
			header = append(header, data...)
			header = append(header, '\n')
		}
		// header读取完毕，读取文本内容, 1为\r
		if contentLength > 0 && len(data) == 1 {
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
