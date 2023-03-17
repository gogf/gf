package main

import (
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gogf/gf/contrib/registry/file/v2"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2"
	"github.com/gogf/gf/v2/os/gfile"

	pb "github.com/gogf/gf/example/rpc/grpcx/rawgrpc/helloworld"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	grpcx.Resolver.Register(file.New(gfile.Temp("gsvc")))

	var (
		ctx     = gctx.GetInitCtx()
		service = gsvc.NewServiceWithName(`hello`)
	)
	// Set up a connection to the server.
	conn, err := grpc.Dial(
		fmt.Sprintf(`%s://%s`, gsvc.Schema, service.GetKey()),
		grpcx.Balancer.WithRandom(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		g.Log().Fatalf(ctx, "did not connect: %v", err)
	}
	defer conn.Close()

	// Send requests.
	client := pb.NewGreeterClient(conn)
	for i := 0; i < 10; i++ {
		res, err := client.SayHello(ctx, &pb.HelloRequest{Name: `GoFrame`})
		if err != nil {
			g.Log().Fatalf(ctx, "could not greet: %+v", err)
		}
		g.Log().Printf(ctx, "Greeting: %s", res.Message)
		time.Sleep(time.Second)
	}
}
