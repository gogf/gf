package main

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"path/filepath"

	"github.com/gogf/gf/net/ghttp"
	"github.com/gogf/gf/os/glog"
)

func SendXmlFile(gameId int, areaName string, filePath string) error {
	path := filepath.FromSlash(filePath)
	fmt.Println(path)
	data := g.Map{
		"gameName":    gameId,
		"area":        areaName,
		"file":        "@file:" + path,
		"contentType": "json",
	}
	if r, err := ghttp.Post("http://127.0.0.1:8199/upload", data); err != nil {
		panic(err)
	} else {
		defer r.Close()
		fmt.Println("ok")
	}
	return nil
}

func main() {
	SendXmlFile(1, "xxx", "/Users/john/Workspace/Go/GOPATH/src/github.com/gogf/gf/.example/net/ghttp/server/session.go")
	return
	path := "/home/john/Workspace/Go/github.com/gogf/gf/version.go"
	r, e := ghttp.Post("http://127.0.0.1:8199/upload", "upload-file=@file:"+path)
	if e != nil {
		glog.Error(e)
	} else {
		fmt.Println(string(r.ReadAll()))
		r.Close()
	}
}
