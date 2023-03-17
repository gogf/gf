package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/gogf/gf/contrib/registry/file/v2"
	"github.com/gogf/gf/example/rpc/grpcx/rawgrpc/helloworld"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gfile"
)

type GreetingServer struct {
	helloworld.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *GreetingServer) SayHello(ctx context.Context, in *helloworld.HelloRequest) (*helloworld.HelloReply, error) {
	return &helloworld.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	gsvc.SetRegistry(file.New(gfile.Temp("gsvc")))

	var (
		err     error
		ctx     = gctx.GetInitCtx()
		address = fmt.Sprintf("%s:%d", gipv4.MustGetIntranetIp(), gtcp.MustGetFreePort())
		service = &gsvc.LocalService{
			Name:      "hello",
			Endpoints: gsvc.NewEndpoints(address),
		}
	)

	// Service registry.
	_, err = gsvc.Register(ctx, service)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = gsvc.Deregister(ctx, service)
	}()

	// Server listening.
	listen, err := net.Listen("tcp", address)
	if err != nil {
		g.Log().Fatalf(ctx, "failed to listen: %v", err)
	}

	s := grpc.NewServer()
	helloworld.RegisterGreeterServer(s, &GreetingServer{})
	g.Log().Printf(ctx, "server listening at %v", listen.Addr())
	if err = s.Serve(listen); err != nil {
		g.Log().Fatalf(ctx, "failed to serve: %v", err)
	}
}
