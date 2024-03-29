// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	"github.com/gogf/gf/contrib/trace/otlphttp/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
)

const (
	serviceName = "otlp-http-client-with-db"
	endpoint    = "tracing-analysis-dc-hz.aliyuncs.com"
	path        = "adapt_******_******/api/otlp/traces"
)

func main() {
	var ctx = gctx.New()
	shutdown, err := otlphttp.Init(serviceName, endpoint, path)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer shutdown()

	StartRequests()
}

// StartRequests starts requests.
func StartRequests() {
	ctx, span := gtrace.NewSpan(gctx.New(), "StartRequests")
	defer span.End()

	var (
		err    error
		client = g.Client()
	)
	// Add user info.
	var insertRes = struct {
		ghttp.DefaultHandlerResponse
		Data struct{ ID int64 } `json:"data"`
	}{}
	err = client.PostVar(ctx, "http://127.0.0.1:8199/user/insert", g.Map{
		"name": "john",
	}).Scan(&insertRes)
	if err != nil {
		panic(err)
	}
	g.Log().Info(ctx, "insert result:", insertRes)
	if insertRes.Data.ID == 0 {
		g.Log().Error(ctx, "retrieve empty id string")
		return
	}

	// Query user info.
	var queryRes = struct {
		ghttp.DefaultHandlerResponse
		Data struct{ User gdb.Record } `json:"data"`
	}{}
	err = client.GetVar(ctx, "http://127.0.0.1:8199/user/query", g.Map{
		"id": insertRes.Data.ID,
	}).Scan(&queryRes)
	if err != nil {
		panic(err)
	}
	g.Log().Info(ctx, "query result:", queryRes)

	// Delete user info.
	var deleteRes = struct {
		ghttp.DefaultHandlerResponse
	}{}
	err = client.PostVar(ctx, "http://127.0.0.1:8199/user/delete", g.Map{
		"id": insertRes.Data.ID,
	}).Scan(&deleteRes)
	if err != nil {
		panic(err)
	}
	g.Log().Info(ctx, "delete result:", deleteRes)
}
