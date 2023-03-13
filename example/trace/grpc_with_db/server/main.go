package main

import (
	_ "github.com/gogf/gf/contrib/drivers/mysql/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"context"
	"fmt"
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/example/trace/grpc_with_db/protocol/user"
	"github.com/gogf/gf/v2/database/gdb"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/os/gctx"
)

type server struct{}

const (
	ServiceName       = "grpc-server-with-db"
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

	s := grpcx.Server.New()
	user.RegisterUserServer(s.Server, &server{})
	s.Run()
}

// Insert is a route handler for inserting user info into database.
func (s *server) Insert(ctx context.Context, req *user.InsertReq) (res *user.InsertRes, err error) {
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
func (s *server) Query(ctx context.Context, req *user.QueryReq) (res *user.QueryRes, err error) {
	err = g.Model("user").Ctx(ctx).Cache(gdb.CacheOption{
		Duration: 5 * time.Second,
		Name:     s.userCacheKey(req.Id),
		Force:    false,
	}).WherePri(req.Id).Scan(&res)
	if err != nil {
		return nil, err
	}
	return
}

// Delete is a route handler for deleting specified user info.
func (s *server) Delete(ctx context.Context, req *user.DeleteReq) (res *user.DeleteRes, err error) {
	err = g.Model("user").Ctx(ctx).Cache(gdb.CacheOption{
		Duration: -1,
		Name:     s.userCacheKey(req.Id),
		Force:    false,
	}).WherePri(req.Id).Scan(&res)
	return
}

func (s *server) userCacheKey(id int32) string {
	return fmt.Sprintf(`userInfo:%d`, id)
}
