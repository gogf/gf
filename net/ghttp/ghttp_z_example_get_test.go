// Copyright 2020 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package ghttp_test

import (
	"fmt"

	"github.com/gogf/gf/frame/g"
)

func ExampleClient_Get() {
	url := "http://127.0.0.1:8999"
	// Send with string parameter along with URL.
	r1, err := g.Client().Get(url + "?id=10000&name=john")
	if err != nil {
		panic(err)
	}
	defer r1.Close()
	fmt.Println(r1.ReadAllString())

	// Send with string parameter in request body.
	r2, err := g.Client().Get(url, "id=10000&name=john")
	if err != nil {
		panic(err)
	}
	defer r2.Close()
	fmt.Println(r2.ReadAllString())

	// Send with map parameter.
	r3, err := g.Client().Get(url, g.Map{
		"id":   10000,
		"name": "john",
	})
	if err != nil {
		panic(err)
	}
	defer r3.Close()
	fmt.Println(r3.ReadAllString())

	// Output:
	// GET: query: 10000, john
	// GET: query: 10000, john
	// GET: query: 10000, john
}

func ExampleClient_GetBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().GetBytes(url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_GetContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().GetContent(url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// GET: query: 10000, john
}

func ExampleClient_GetVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		user *User
		url  = "http://127.0.0.1:8999/var/json"
	)
	err := g.Client().GetVar(url).Scan(&user)
	if err != nil {
		panic(err)
	}
	fmt.Println(user)

	// Output:
	// &{1 john}
}
