package apollo

import (
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/test/gtest"
	"testing"
)

func TestApollo(t *testing.T) {
	ctx := gctx.New()

	//g.Dump(g.Cfg().Data(ctx))

	// 测试1
	timeoutVar, err := g.Cfg().Get(ctx, "timeout")
	gtest.AssertNil(err)
	gtest.AssertEQ(timeoutVar.Int(), 100)

	// 测试2，中间有点
	serverAddressVar, err := g.Cfg().Get(ctx, "server.address")
	gtest.AssertNil(err)
	gtest.AssertEQ(serverAddressVar.String(), ":8000")

	// 启动服务
	//s := g.Server()
	//s.Group("/", func(group *ghttp.RouterGroup) {
	//	group.Middleware(ghttp.MiddlewareHandlerResponse)
	//	group.Bind()
	//})
	//s.Run()
}
