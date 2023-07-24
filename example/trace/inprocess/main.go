package main

import (
	"context"

	"github.com/gogf/gf/contrib/trace/jaeger/v2"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/gtrace"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/util/gutil"
)

const (
	ServiceName       = "inprocess"
	JaegerUdpEndpoint = "localhost:6831"
)

func main() {
	var ctx = gctx.New()
	tp, err := jaeger.Init(ServiceName, JaegerUdpEndpoint)
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	defer tp.Shutdown(ctx)

	ctx, span := gtrace.NewSpan(ctx, "main")
	defer span.End()

	// Trace 1.
	user1 := GetUser(ctx, 1)
	g.Dump(user1)

	// Trace 2.
	user100 := GetUser(ctx, 100)
	g.Dump(user100)
}

// GetUser retrieves and returns hard coded user data for demonstration.
func GetUser(ctx context.Context, id int) g.Map {
	ctx, span := gtrace.NewSpan(ctx, "GetUser")
	defer span.End()
	m := g.Map{}
	gutil.MapMerge(
		m,
		GetInfo(ctx, id),
		GetDetail(ctx, id),
		GetScores(ctx, id),
	)
	return m
}

// GetInfo retrieves and returns hard coded user info for demonstration.
func GetInfo(ctx context.Context, id int) g.Map {
	ctx, span := gtrace.NewSpan(ctx, "GetInfo")
	defer span.End()
	if id == 100 {
		return g.Map{
			"id":     100,
			"name":   "john",
			"gender": 1,
		}
	}
	return nil
}

// GetDetail retrieves and returns hard coded user detail for demonstration.
func GetDetail(ctx context.Context, id int) g.Map {
	ctx, span := gtrace.NewSpan(ctx, "GetDetail")
	defer span.End()
	if id == 100 {
		return g.Map{
			"site":  "https://goframe.org",
			"email": "john@goframe.org",
		}
	}
	return nil
}

// GetScores retrieves and returns hard coded user scores for demonstration.
func GetScores(ctx context.Context, id int) g.Map {
	ctx, span := gtrace.NewSpan(ctx, "GetScores")
	defer span.End()
	if id == 100 {
		return g.Map{
			"math":    100,
			"english": 60,
			"chinese": 50,
		}
	}
	return nil
}
