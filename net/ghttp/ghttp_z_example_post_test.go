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

func ExampleClient_Post() {
	url := "http://127.0.0.1:8999"
	// Send with string parameter in request body.
	r1, err := g.Client().Post(url, "id=10000&name=john")
	if err != nil {
		panic(err)
	}
	defer r1.Close()
	fmt.Println(r1.ReadAllString())

	// Send with map parameter.
	r2, err := g.Client().Post(url, g.Map{
		"id":   10000,
		"name": "john",
	})
	if err != nil {
		panic(err)
	}
	defer r2.Close()
	fmt.Println(r2.ReadAllString())

	// Output:
	// POST: form: 10000, john
	// POST: form: 10000, john
}

func ExampleClient_PostBytes() {
	url := "http://127.0.0.1:8999"
	fmt.Println(string(g.Client().PostBytes(url, g.Map{
		"id":   10000,
		"name": "john",
	})))

	// Output:
	// POST: form: 10000, john
}

func ExampleClient_PostContent() {
	url := "http://127.0.0.1:8999"
	fmt.Println(g.Client().PostContent(url, g.Map{
		"id":   10000,
		"name": "john",
	}))

	// Output:
	// POST: form: 10000, john
}

func ExampleClient_PostVar() {
	type User struct {
		Id   int
		Name string
	}
	var (
		users []User
		url   = "http://127.0.0.1:8999/var/jsons"
	)
	err := g.Client().PostVar(url).Scan(&users)
	if err != nil {
		panic(err)
	}
	fmt.Println(users)

	// Output:
	// [{1 john} {2 smith}]
}
