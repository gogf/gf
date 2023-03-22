package main

import (
	"context"
	"fmt"
	"time"

	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

type cTrace struct{}

const (
	ServiceName       = "http-server-with-db"
	JaegerUdpEndpoint = "localhost:6831"
)

func main() {
	var ctx = gctx.New()
	tp, err := jaeger.Init(ServiceName, JaegerUdpEndpoint)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer tp.Shutdown(ctx)

	// Set ORM cache adapter with redis.
	g.DB().GetCache().SetAdapter(gcache.NewAdapterRedis(g.Redis()))

	// Start HTTP server.
	s := g.Server()
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.ALL("/user", new(cTrace))
	})
	s.SetPort(8199)
	s.Run()
}

type InsertReq struct {
	Name string `v:"required#Please input user name."`
}
type InsertRes struct {
	Id int64
}

// Insert is a route handler for inserting user info into database.
func (c *cTrace) Insert(ctx context.Context, req *InsertReq) (res *InsertRes, err error) {
	result, err := g.Model("user").Ctx(ctx).Insert(req)
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	res = &InsertRes{
		Id: id,
	}
	return
}

type QueryReq struct {
	Id int `v:"min:1#User id is required for querying"`
}
type QueryRes struct {
	User gdb.Record
}

// Query is a route handler for querying user info. It firstly retrieves the info from redis,
// if there's nothing in the redis, it then does db select.
func (c *cTrace) Query(ctx context.Context, req *QueryReq) (res *QueryRes, err error) {
	one, err := g.Model("user").Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 5 * time.Second,
		Name:     c.userCacheKey(req.Id),
		Force:    false,
	}).WherePri(req.Id).One()
	if err != nil {
		return nil, err
	}
	res = &QueryRes{
		User: one,
	}
	return
}

type DeleteReq struct {
	Id int `v:"min:1#User id is required for deleting."`
}
type DeleteRes struct{}

// Delete is a route handler for deleting specified user info.
func (c *cTrace) Delete(ctx context.Context, req *DeleteReq) (res *DeleteRes, err error) {
	_, err = g.Model("user").Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     c.userCacheKey(req.Id),
		Force:    false,
	}).WherePri(req.Id).Delete()
	if err != nil {
		return nil, err
	}
	return
}

func (c *cTrace) userCacheKey(id int) string {
	return fmt.Sprintf(`userInfo:%d`, id)
}
