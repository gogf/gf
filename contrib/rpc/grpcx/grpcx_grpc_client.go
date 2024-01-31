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

	"github.com/gogf/gf/v2/errors/gcode"
	"github.com/gogf/gf/v2/errors/gerror"
	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/text/gstr"
)

// DefaultGrpcDialOptions returns the default options for creating grpc client connection.
func (c modClient) DefaultGrpcDialOptions() []grpc.DialOption {
	return []grpc.DialOption{
		Balancer.WithName(gsel.GetBuilder().Name()),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
}

// NewGrpcClientConn creates and returns a client connection for given service `appId`.
func (c modClient) NewGrpcClientConn(serviceNameOrAddress string, opts ...grpc.DialOption) (*grpc.ClientConn, error) {
	autoLoadAndRegisterFileRegistry()

	var (
		dialAddress       = serviceNameOrAddress
		grpcClientOptions = make([]grpc.DialOption, 0)
	)
	if isServiceName(serviceNameOrAddress) {
		dialAddress = fmt.Sprintf(
			`%s://%s`,
			gsvc.Schema, gsvc.NewServiceWithName(serviceNameOrAddress).GetKey(),
		)
	} else {
		addressParts := gstr.Split(serviceNameOrAddress, gsvc.EndpointHostPortDelimiter)
		switch len(addressParts) {
		case 2:
			if addressParts[0] == "" {
				return nil, gerror.NewCodef(
					gcode.CodeInvalidParameter,
					`invalid address "%s" for client, missing host`,
					serviceNameOrAddress,
				)
			}
		default:
			return nil, gerror.NewCodef(
				gcode.CodeInvalidParameter,
				`invalid address "%s" for client`,
				serviceNameOrAddress,
			)
		}
	}
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
	conn, err := grpc.Dial(dialAddress, grpcClientOptions...)
	if err != nil {
		return nil, err
	}
	return conn, nil
}

// MustNewGrpcClientConn creates and returns a client connection for given service `appId`.
// It panics if any error occurs.
func (c modClient) MustNewGrpcClientConn(serviceNameOrAddress string, opts ...grpc.DialOption) *grpc.ClientConn {
	conn, err := c.NewGrpcClientConn(serviceNameOrAddress, opts...)
	if err != nil {
		panic(err)
	}
	return conn
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

// isServiceName checks and returns whether given input parameter is service name or not.
// It checks by whether the parameter is address by containing port delimiter character ':'.
//
// It does not contain any port number if using service discovery.
func isServiceName(serviceNameOrAddress string) bool {
	return !gstr.Contains(serviceNameOrAddress, gsvc.EndpointHostPortDelimiter)
}
