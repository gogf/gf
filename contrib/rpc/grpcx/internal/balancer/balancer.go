// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package balancer defines APIs for load balancing in gRPC.
package balancer

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/balancer"
	"google.golang.org/grpc/balancer/base"

	"github.com/gogf/gf/v2/net/gsel"
)

type Balancer struct{}

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
	b := Balancer{}
	b.Register(Random, Weight, RoundRobin, LeastConnection)
}

// Register registers the given balancer builder with the given name.
func (Balancer) Register(builders ...gsel.Builder) {
	for _, builder := range builders {
		balancer.Register(
			base.NewBalancerBuilder(
				builder.Name(),
				&Builder{builder: builder},
				base.Config{HealthCheck: true},
			),
		)
	}
}

// WithRandom returns a grpc.DialOption which enables random load balancing.
func (b Balancer) WithRandom() grpc.DialOption {
	return b.WithName(Random.Name())
}

// WithWeight returns a grpc.DialOption which enables weight load balancing.
func (b Balancer) WithWeight() grpc.DialOption {
	return b.WithName(Weight.Name())
}

// WithRoundRobin returns a grpc.DialOption which enables round-robin load balancing.
func (b Balancer) WithRoundRobin() grpc.DialOption {
	return b.WithName(RoundRobin.Name())
}

// WithLeastConnection returns a grpc.DialOption which enables the least connection load balancing.
func (b Balancer) WithLeastConnection() grpc.DialOption {
	return b.WithName(LeastConnection.Name())
}

// WithName returns a grpc.DialOption which enables the load balancing by name.
func (b Balancer) WithName(name string) grpc.DialOption {
	return grpc.WithDefaultServiceConfig(fmt.Sprintf(
		`{"loadBalancingPolicy": "%s"}`,
		name,
	))
}
