// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/trace/otlpgrpc/v2"
	"github.com/gogf/gf/example/trace/grpc-with-db/protobuf/user"

	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

// Controller is the gRPC controller for user management.
type Controller struct {
	user.UnimplementedUserServer
}

const (
	serviceName = "otlp-grpc-server"
	endpoint    = "tracing-analysis-dc-bj.aliyuncs.com:8090"
	traceToken  = "******_******"
)

func main() {
	grpcx.Resolver.Register(etcd.New("127.0.0.1:2379"))

	var ctx = gctx.New()
	shutdown, err := otlpgrpc.Init(serviceName, endpoint, traceToken)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer shutdown()

	// Set ORM cache adapter with redis.
	g.DB().GetCache().SetAdapter(gcache.NewAdapterRedis(g.Redis()))

	s := grpcx.Server.New()
	user.RegisterUserServer(s.Server, &Controller{})
	s.Run()
}

// Insert is a route handler for inserting user info into database.
func (s *Controller) Insert(ctx context.Context, req *user.InsertReq) (res *user.InsertRes, err error) {
	result, err := g.Model("user").Ctx(ctx).Insert(g.Map{
		"name": req.Name,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	res = &user.InsertRes{
		Id: int32(id),
	}
	return
}

// Query is a route handler for querying user info. It firstly retrieves the info from redis,
// if there's nothing in the redis, it then does db select.
func (s *Controller) Query(ctx context.Context, req *user.QueryReq) (res *user.QueryRes, err error) {
	if err = g.Model("user").Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 5 * time.Second,
		Name:     s.userCacheKey(req.Id),
		Force:    false,
	}).WherePri(req.Id).Scan(&res); err != nil {
		return nil, err
	}
	return
}

// Delete is a route handler for deleting specified user info.
func (s *Controller) Delete(ctx context.Context, req *user.DeleteReq) (res *user.DeleteRes, err error) {
	err = g.Model("user").Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     s.userCacheKey(req.Id),
		Force:    false,
	}).WherePri(req.Id).Scan(&res)
	return
}

func (s *Controller) userCacheKey(id int32) string {
	return fmt.Sprintf(`userInfo:%d`, id)
}
