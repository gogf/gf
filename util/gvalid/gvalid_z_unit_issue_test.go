// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"fmt"
	"runtime"
	"testing"
	"time"

	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/guid"
	"github.com/gogf/gf/v2/util/gvalid"
)

type Foo struct {
	Bar *Bar `p:"bar" v:"required-without:Baz"`
	Baz *Baz `p:"baz" v:"required-without:Bar"`
}
type Bar struct {
	BarKey string `p:"bar_key" v:"required"`
}
type Baz struct {
	BazKey string `p:"baz_key" v:"required"`
}

// https://github.com/gogf/gf/issues/2503
func Test_Issue2503(t *testing.T) {
	foo := &Foo{
		Bar: &Bar{BarKey: "value"},
	}
	err := gvalid.New().Data(foo).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}

type Issue3636SliceV struct{}

func init() {
	rule := Issue3636SliceV{}
	gvalid.RegisterRule(rule.Name(), rule.Run)
}

func (r Issue3636SliceV) Name() string {
	return "slice-v"
}

func (r Issue3636SliceV) Message() string {
	return "not a slice"
}

func (r Issue3636SliceV) Run(_ context.Context, in gvalid.RuleFuncInput) error {
	for _, v := range in.Value.Slice() {
		if v == "" {
			return gerror.New("empty value")
		}
	}
	if !in.Value.IsSlice() {
		return gerror.New("not a slice")
	}
	return nil
}

type Issue3636HelloReq struct {
	g.Meta `path:"/hello" method:"POST"`

	Name string   `json:"name" v:"required" dc:"Your name"`
	S    []string `json:"s" v:"slice-v" dc:"S"`
}
type Issue3636HelloRes struct {
	Name string   `json:"name" v:"required" dc:"Your name"`
	S    []string `json:"s" v:"slice-v" dc:"S"`
}

type Issue3636Hello struct{}

func (Issue3636Hello) Say(ctx context.Context, req *Issue3636HelloReq) (res *Issue3636HelloRes, err error) {
	res = &Issue3636HelloRes{
		Name: req.Name,
		S:    req.S,
	}
	return
}

// https://github.com/gogf/gf/issues/3636
func Test_Issue3636(t *testing.T) {
	s := g.Server(guid.S())
	s.Use(ghttp.MiddlewareHandlerResponse)
	s.Group("/", func(group *ghttp.RouterGroup) {
		group.Bind(
			new(Issue3636Hello),
		)
	})
	s.SetDumpRouterMap(false)
	s.Start()
	defer s.Shutdown()

	time.Sleep(100 * time.Millisecond)

	gtest.C(t, func(t *gtest.T) {
		c := g.Client()
		c.SetPrefix(fmt.Sprintf("http://127.0.0.1:%d", s.GetListenedPort()))
		t.Assert(
			c.PostContent(ctx, "/hello", `{"name": "t", "s" : []}`),
			`{"code":0,"message":"OK","data":{"name":"t","s":[]}}`,
		)
	})
}

// https://github.com/gogf/gf/issues/4092
func Test_Issue4092(t *testing.T) {
	type Model struct {
		Raw  []byte `v:"required"`
		Test []byte `v:"foreach|in:1,2,3"`
	}
	gtest.C(t, func(t *gtest.T) {
		const kb = 1024
		const mb = 1024 * kb
		raw := make([]byte, 50*mb)
		in := &Model{
			Raw:  raw,
			Test: []byte{40, 5, 6},
		}
		err := g.Validator().
			Data(in).
			Run(context.Background())
		t.Assert(err, "The Test value `6` is not in acceptable range: 1,2,3")
		allocMb := getMemAlloc()
		t.AssertLE(allocMb, 110)
	})
}

func getMemAlloc() uint64 {
	byteToMb := func(b uint64) uint64 {
		return b / 1024 / 1024
	}
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	// For info on each, see: https://golang.org/pkg/runtime/#MemStats
	alloc := byteToMb(m.Alloc)
	return alloc
}
