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

var (
	Random          = gsel.NewBuilderRandom()
	Weight          = gsel.NewBuilderWeight()
	RoundRobin      = gsel.NewBuilderRoundRobin()
	LeastConnection = gsel.NewBuilderLeastConnection()
)

func init() {
	Register(Random, Weight, RoundRobin, LeastConnection)
}

// Register registers the given balancer builder with the given name.
func Register(builders ...gsel.Builder) {
	for _, builder := range builders {
		balancer.Register(base.NewBalancerBuilder(
			builder.Name(),
			&Builder{builder: builder},
			base.Config{HealthCheck: true},
		))
	}
}

// WithRandom returns a grpc.DialOption which enables random load balancing.
func WithRandom() grpc.DialOption {
	return WithName(Random.Name())
}

// WithWeight returns a grpc.DialOption which enables weight load balancing.
func WithWeight() grpc.DialOption {
	return WithName(Weight.Name())
}

// WithRoundRobin returns a grpc.DialOption which enables round-robin load balancing.
func WithRoundRobin() grpc.DialOption {
	return WithName(RoundRobin.Name())
}

// WithLeastConnection returns a grpc.DialOption which enables the least connection load balancing.
func WithLeastConnection() grpc.DialOption {
	return WithName(LeastConnection.Name())
}

// WithName returns a grpc.DialOption which enables the load balancing by name.
func WithName(name string) grpc.DialOption {
	return grpc.WithDefaultServiceConfig(fmt.Sprintf(
		`{"loadBalancingPolicy": "%s"}`,
		name,
	))
}
