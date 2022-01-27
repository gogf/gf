// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package balancer defines APIs for load balancing in gRPC.
package balancer

import (
	"fmt"

	"github.com/gogf/gf/v2/net/gsel"
	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"
)

const (
	rawSvcKeyInSubConnInfo = `RawService`
)

func init() {
	Register(gsel.SelectorRandom, gsel.NewSelectorRandom())
	Register(gsel.SelectorWeight, gsel.NewSelectorWeight())
	Register(gsel.SelectorRoundRobin, gsel.NewSelectorRoundRobin())
	Register(gsel.SelectorLeastConnection, gsel.NewSelectorLeastConnection())
}

func Register(name string, selector gsel.Selector) {
	balancer.Register(base.NewBalancerBuilder(
		name,
		&Builder{selector: selector},
		base.Config{HealthCheck: true},
	))
}

func WithRandom() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorRandom)
}

func WithWeight() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorWeight)
}

func WithRoundRobin() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorRoundRobin)
}

func WithLeastConnection() grpc.DialOption {
	return doWithSelectorName(gsel.SelectorLeastConnection)
}

func doWithSelectorName(name string) grpc.DialOption {
	return grpc.WithDefaultServiceConfig(fmt.Sprintf(
		`{"loadBalancingPolicy": "%s"}`,
		name,
	))
}
