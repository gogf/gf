package main

import (
	"context"
	"fmt"
	"net"

	"github.com/gogf/gf/contrib/registry/etcd/v2"
	"github.com/gogf/gf/contrib/resolver/v2"
	pb "github.com/gogf/gf/example/registry/rawgrpc/helloworld"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gipv4"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/net/gtcp"
	"github.com/gogf/gf/v2/os/gctx"
	"google.golang.org/grpc"
)

type server struct {
	pb.UnimplementedGreeterServer
}

func (s *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
	g.Log().Printf(ctx, "Received: %v", in.GetName())
	return &pb.HelloReply{Message: "Hello " + in.GetName()}, nil
}

func main() {
	resolver.SetRegistry(etcd.New("127.0.0.1:2379"))

	var (
		err     error
		ctx     = gctx.New()
		address = fmt.Sprintf("%s:%d", gipv4.MustGetIntranetIp(), gtcp.MustGetFreePort())
		service = &gsvc.Service{
			Name:      "hello",
			Endpoints: []string{address},
		}
	)
	err = gsvc.Register(ctx, service)
	if err != nil {
		panic(err)
	}
	defer func() {
		_ = gsvc.Deregister(ctx, service)
	}()

	lis, err := net.Listen("tcp", address)
	if err != nil {
		g.Log().Fatalf(ctx, "failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGreeterServer(s, &server{})
	g.Log().Printf(ctx, "server listening at %v", lis.Addr())
	if err = s.Serve(lis); err != nil {
		g.Log().Fatalf(ctx, "failed to serve: %v", err)
	}
}
