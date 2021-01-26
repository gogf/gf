package main

import (
	"context"
	"fmt"
	"github.com/gogf/gcache-adapter/adapter"
	"github.com/gogf/gf/.example/net/gtrace/5.grpc+db+redis+log/protobuf/user"
	"github.com/gogf/gf/errors/gerror"
	"github.com/gogf/gf/frame/g"
	"github.com/gogf/gf/util/gvalid"
	"go.opentelemetry.io/otel/exporters/trace/jaeger"
	sdkTrace "go.opentelemetry.io/otel/sdk/trace"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"net"
	"time"
)

type server struct{}

const (
	JaegerEndpoint = "http://localhost:14268/s/traces"
	ServiceName    = "tracing-grpc-server"
)

// initTracer creates a new trace provider instance and registers it as global trace provider.
func initTracer() func() {
	// Create and install Jaeger export pipeline.
	flush, err := jaeger.InstallNewPipeline(
		jaeger.WithCollectorEndpoint(JaegerEndpoint),
		jaeger.WithProcess(jaeger.Process{
			ServiceName: ServiceName,
		}),
		jaeger.WithSDK(&sdkTrace.Config{DefaultSampler: sdkTrace.AlwaysSample()}),
	)
	if err != nil {
		g.Log().Fatal(err)
	}
	return flush
}

// Common validation unary interpreter.
func UnaryValidate(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	// It does nothing if there's no validation tag in the struct definition.
	if err := gvalid.CheckStruct(req, nil); err != nil {
		return nil, gerror.NewCode(
			int(codes.InvalidArgument),
			gerror.Current(err).Error(),
		)
	}
	return handler(ctx, req)
}

// Insert is a route handler for inserting user info into dtabase.
func (s *server) Insert(ctx context.Context, req *user.InsertReq) (*user.InsertRes, error) {
	res := user.InsertRes{}
	result, err := g.Table("user").Ctx(ctx).Insert(g.Map{
		"name": req.Name,
	})
	if err != nil {
		return nil, err
	}
	id, _ := result.LastInsertId()
	res.Id = int32(id)
	return &res, nil
}

// Query is a route handler for querying user info. It firstly retrieves the info from redis,
// if there's nothing in the redis, it then does db select.
func (s *server) Query(ctx context.Context, req *user.QueryReq) (*user.QueryRes, error) {
	res := user.QueryRes{}
	err := g.Table("user").
		Ctx(ctx).
		Cache(5*time.Second, s.userCacheKey(req.Id)).
		WherePri(req.Id).
		Scan(&res)
	if err != nil {
		return nil, err
	}
	return &res, nil
}

// Delete is a route handler for deleting specified user info.
func (s *server) Delete(ctx context.Context, req *user.DeleteReq) (*user.DeleteRes, error) {
	res := user.DeleteRes{}
	_, err := g.Table("user").
		Ctx(ctx).
		Cache(-1, s.userCacheKey(req.Id)).
		WherePri(req.Id).
		Delete()
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (s *server) userCacheKey(id int32) string {
	return fmt.Sprintf(`userInfo:%d`, id)
}

func main() {
	flush := initTracer()
	defer flush()

	g.DB().GetCache().SetAdapter(adapter.NewRedis(g.Redis()))

	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		g.Log().Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(UnaryValidate),
	)
	user.RegisterUserServer(s, &server{})
	if err := s.Serve(listen); err != nil {
		g.Log().Fatalf("failed to serve: %v", err)
	}
}
