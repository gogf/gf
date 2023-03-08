package main

import (
	"context"
	"fmt"
	"net"

	"google.golang.org/grpc"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	pb "github.com/gogf/gf/example/rpc/grpcx/rawgrpc/helloworld"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gctx"
)

type GreetingServer struct {
	pb.UnimplementedGreeterServer
}

// SayHello implements helloworld.GreeterServer
func (s *GreetingServer) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	g.Log().Printf(ctx, "Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	gsvc.SetRegistry(etcd.New("127.0.0.1:2379"))

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
	pb.RegisterGreeterServer(s, &GreetingServer{})
	g.Log().Printf(ctx, "server listening at %v", listen.Addr())
	if err = s.Serve(listen); err != nil {
		g.Log().Fatalf(ctx, "failed to serve: %v", err)
	}
}
