// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis_test

import (
	"context"
	"fmt"
	"github.com/gogf/gf/container/gvar"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gutil"
)

func Example_autoMarshalUnmarshalMap() {
	var (
		err    error
		result *gvar.Var
		ctx    = context.Background()
		key    = "user"
		data   = g.Map{
			"id":   10000,
			"name": "john",
		}
	)
	_, err = g.Redis().Do(ctx, "SET", key, data)
	if err != nil {
		panic(err)
	}
	result, err = g.Redis().Do(ctx, "GET", key)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Map())
}

func Example_autoMarshalUnmarshalStruct() {
	type User struct {
		Id   int
		Name string
	}
	var (
		err    error
		result *gvar.Var
		ctx    = context.Background()
		key    = "user"
		user   = &User{
			Id:   10000,
			Name: "john",
		}
	)

	_, err = g.Redis().Do(ctx, "SET", key, user)
	if err != nil {
		panic(err)
	}
	result, err = g.Redis().Do(ctx, "GET", key)
	if err != nil {
		panic(err)
	}

	var user2 *User
	if err = result.Struct(&user2); err != nil {
		panic(err)
	}
	fmt.Println(user2.Id, user2.Name)
}

func Example_autoMarshalUnmarshalStructSlice() {
	type User struct {
		Id   int
		Name string
	}
	var (
		err    error
		result *gvar.Var
		ctx    = context.Background()
		key    = "user-slice"
		users1 = []User{
			{
				Id:   1,
				Name: "john1",
			},
			{
				Id:   2,
				Name: "john2",
			},
		}
	)

	_, err = g.Redis().Do(ctx, "SET", key, users1)
	if err != nil {
		panic(err)
	}
	result, err = g.Redis().Do(ctx, "GET", key)
	if err != nil {
		panic(err)
	}

	var users2 []User
	if err = result.Structs(&users2); err != nil {
		panic(err)
	}
	fmt.Println(users2)
}

func Example_hSet() {
	var (
		err    error
		result *gvar.Var
		ctx    = context.Background()
		key    = "user"
	)
	_, err = g.Redis().Do(ctx, "HSET", key, "id", 10000)
	if err != nil {
		panic(err)
	}
	_, err = g.Redis().Do(ctx, "HSET", key, "name", "john")
	if err != nil {
		panic(err)
	}
	result, err = g.Redis().Do(ctx, "HGETALL", key)
	if err != nil {
		panic(err)
	}
	fmt.Println(result.Map())

	// May Output:
	// map[id:10000 name:john]
}

func Example_hMSet_Map() {
	var (
		ctx  = context.Background()
		key  = "user_100"
		data = g.Map{
			"name":  "gf",
			"sex":   0,
			"score": 100,
		}
	)
	_, err := g.Redis().Do(ctx, "HMSET", append(g.Slice{key}, gutil.MapToSlice(data)...)...)
	if err != nil {
		g.Log().Fatal(err)
	}
	v, err := g.Redis().Do(ctx, "HMGET", key, "name")
	if err != nil {
		g.Log().Fatal(err)
	}
	fmt.Println(v.Slice())

	// May Output:
	// [gf]
}

func Example_hMSet_Struct() {
	type User struct {
		Name  string `json:"name"`
		Sex   int    `json:"sex"`
		Score int    `json:"score"`
	}
	var (
		ctx  = context.Background()
		key  = "user_100"
		data = &User{
			Name:  "gf",
			Sex:   0,
			Score: 100,
		}
	)
	_, err := g.Redis().Do(ctx, "HMSET", append(g.Slice{key}, gutil.StructToSlice(data)...)...)
	if err != nil {
		g.Log().Fatal(err)
	}
	v, err := g.Redis().Do(ctx, "HMGET", key, "name")
	if err != nil {
		g.Log().Fatal(err)
	}
	fmt.Println(v.Slice())

	// May Output:
	// ["gf"]
}
