// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gfile"
	"time"
)

func ExampleHelloWorld() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello world")
	})
	s.SetPort(8999)
	s.Run()
}

// Custom saving file name.
func ExampleUploadFile_Save() {
	s := g.Server()
	s.BindHandler("/upload", func(r *ghttp.Request) {
		file := r.GetUploadFile("TestFile")
		if file == nil {
			r.Response.Write("empty file")
			return
		}
		file.Filename = "MyCustomFileName.txt"
		fileName, err := file.Save(gfile.TempDir())
		if err != nil {
			r.Response.Write(err)
			return
		}
		r.Response.Write(fileName)
	})
	s.SetPort(8999)
	s.Run()
}

func ExampleClientResponse_RawDump() {
	var (
		ctx = context.Background()
	)
	response, err := g.Client().Get(ctx, "https://goframe.org")
	if err != nil {
		panic(err)
	}
	response.RawDump()
}

// ExampleClient_SetProxy a example for `ghttp.Client.SetProxy` method.
// please prepare two proxy server before running this example.
// http proxy server listening on `127.0.0.1:1081`
// socks5 proxy server listening on `127.0.0.1:1080`
func ExampleClient_SetProxy() {
	// connect to a http proxy server
	client := g.Client()
	client.SetProxy("http://127.0.0.1:1081")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err := client.Get(ctx, "https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	response.RawDump()
	// connect to a http proxy server which needs auth
	client.SetProxy("http://user:password:127.0.0.1:1081")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err = client.Get(ctx, "https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	response.RawDump()

	// connect to a socks5 proxy server
	client.SetProxy("socks5://127.0.0.1:1080")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err = client.Get(ctx, "https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": socks connect tcp 127.0.0.1:1087->api.ip.sb:443: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	fmt.Println(response.RawResponse())

	// connect to a socks5 proxy server which needs auth
	client.SetProxy("socks5://user:password@127.0.0.1:1080")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err = client.Get(ctx, "https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": socks connect tcp 127.0.0.1:1087->api.ip.sb:443: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	fmt.Println(response.RawResponse())
}

// ExampleClientChain_Proxy a chain version of example for `ghttp.Client.Proxy` method.
// please prepare two proxy server before running this example.
// http proxy server listening on `127.0.0.1:1081`
// socks5 proxy server listening on `127.0.0.1:1080`
// for more details, please refer to ExampleClient_SetProxy
func ExampleClientChain_Proxy() {
	var (
		ctx = context.Background()
	)
	client := g.Client()
	response, err := client.Proxy("http://127.0.0.1:1081").Get(ctx, "https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	fmt.Println(response.RawResponse())

	client2 := g.Client()
	response, err = client2.Proxy("socks5://127.0.0.1:1080").Get(ctx, "https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": socks connect tcp 127.0.0.1:1087->api.ip.sb:443: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	fmt.Println(response.RawResponse())
}
