// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcx

import (
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
)

// DefaultGrpcDialOptions returns the default options for creating grpc client connection.
func (c modClient) DefaultGrpcDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		Balancer.WithName(gsel.GetBuilder().Name()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

// NewGrpcClientConn NewGrpcConn creates and returns a client connection for given service `appId`.
func (c modClient) NewGrpcClientConn(name string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	autoLoadAndRegisterEtcdRegistry()

	var (
		service           = gsvc.NewServiceWithName(name)
		grpcClientOptions = make([]grpc.DialOption, 0)
	)
	grpcClientOptions = append(grpcClientOptions, c.DefaultGrpcDialOptions()...)
	if len(opts) > 0 {
		grpcClientOptions = append(grpcClientOptions, opts...)
	}
	grpcClientOptions = append(grpcClientOptions, c.ChainUnary(
		c.UnaryTracing,
		c.UnaryError,
	))
	grpcClientOptions = append(grpcClientOptions, c.ChainStream(
		c.StreamTracing,
	))
	conn, err := grpc.Dial(fmt.Sprintf(`%s://%s`, gsvc.Schema, service.GetKey()), grpcClientOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// ChainUnary creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainUnaryClient(one, two, three) will execute one before two before three.
func (c modClient) ChainUnary(interceptors ...grpc.UnaryClientInterceptor) grpc.DialOption {
	return grpc.WithChainUnaryInterceptor(interceptors...)
}

// ChainStream creates a single interceptor out of a chain of many interceptors.
//
// Execution is done in left-to-right order, including passing of context.
// For example ChainStreamClient(one, two, three) will execute one before two before three.
func (c modClient) ChainStream(interceptors ...grpc.StreamClientInterceptor) grpc.DialOption {
	return grpc.WithChainStreamInterceptor(interceptors...)
}
