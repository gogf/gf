package ghttp

import (
	"context"
	"net/http"
	"reflect"
	"strings"
	"testing"
)

type HelloReq struct {
	//g.Meta  `path:"/hello" tags:"Hello" method:"post" summary:"You first hello api"`
	Id      int32  `p:"id"  dc:"Id" json:"id"`
	Name    string `p:"name" dc:"Id" json:"name"`
	Address string `p:"address" dc:"Id" json:"address"`
}

type HelloRes struct {
	//g.Meta  `mime:"text/html" example:"string"`
	Id      int32  `json:"id"`
	Name    string `json:"name"`
	Address string `json:"address"`
	Ip      string `json:"ip"`
}

// 原始的 createRouterFunc 性能测试
// go test -bench='Benchmark(Original|Optimized)' -run=none -benchmem -benchtime=2s -count=3
func BenchmarkOriginal(b *testing.B) {
	req, funcInfo := SetupTest()
	routerFunc := createRouterFuncOriginal(funcInfo)

	// 执行基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		routerFunc(req)
	}
}

// 优化后的 createRouterFunc 性能测试
// go test -bench='Benchmark(Original|Optimized)' -run=none -benchmem -benchtime=2s -count=3
func BenchmarkOptimized(b *testing.B) {
	req, funcInfo := SetupTest()
	routerFunc := createRouterFunc(funcInfo)

	// 执行基准测试
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		routerFunc(req)
	}
}

func SetupTest() (*Request, handlerFuncInfo) {
	// 创建一个模拟的 HandlerFunc
	mockHandler := func(r *Request) {
		r.Host = "127.0.0.1"
	}

	mockControllerHandler := func(ctx context.Context, r *HelloReq) (res *HelloRes, err error) {
		// 实现具体的处理逻辑
		return res, nil
	}
	// 创建 handlerFuncInfo 结构体
	funcInfo := handlerFuncInfo{
		Func:            mockHandler,
		Type:            reflect.TypeOf(mockControllerHandler),
		Value:           reflect.ValueOf(mockControllerHandler),
		IsStrictRoute:   false,
		ReqStructFields: nil,
	}

	// 创建一个模拟的请求
	req := &Request{}
	request, err := http.NewRequest("GET", "/hello", strings.NewReader("name=bar"))
	if err != nil {
		return nil, handlerFuncInfo{}
	}
	req.Request = request
	return req, funcInfo
}
