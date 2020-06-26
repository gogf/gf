// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/net/ghttp"
	"time"
)

func ExampleGetServer() {
	s := g.Server()
	s.BindHandler("/", func(r *ghttp.Request) {
		r.Response.Write("hello world")
	})
	s.SetPort(8999)
	s.Run()
}

func ExampleClientResponse_RawDump() {
	response, err := g.Client().Get("https://goframe.org")
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
	client := ghttp.NewClient()
	client.SetProxy("http://127.0.0.1:1081")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err := client.Get("https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	response.RawDump()
	// connect to a http proxy server which needs auth
	client.SetProxy("http://user:password:127.0.0.1:1081")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err = client.Get("https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	response.RawDump()

	// connect to a socks5 proxy server
	client.SetProxy("socks5://127.0.0.1:1080")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err = client.Get("https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": socks connect tcp 127.0.0.1:1087->api.ip.sb:443: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	fmt.Println(response.RawResponse())

	// connect to a socks5 proxy server which needs auth
	client.SetProxy("socks5://user:password@127.0.0.1:1080")
	client.SetTimeout(5 * time.Second) // it's suggested to set http client timeout
	response, err = client.Get("https://api.ip.sb/ip")
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
	client := ghttp.NewClient()
	response, err := client.Proxy("http://127.0.0.1:1081").Get("https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": proxyconnect tcp: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	fmt.Println(response.RawResponse())

	client2 := ghttp.NewClient()
	response, err = client2.Proxy("socks5://127.0.0.1:1080").Get("https://api.ip.sb/ip")
	if err != nil {
		// err is not nil when your proxy server is down.
		// eg. Get "https://api.ip.sb/ip": socks connect tcp 127.0.0.1:1087->api.ip.sb:443: dial tcp 127.0.0.1:1087: connect: connection refused
		fmt.Println(err)
	}
	fmt.Println(response.RawResponse())
}
