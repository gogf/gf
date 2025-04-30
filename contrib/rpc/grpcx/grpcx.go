// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpcx provides grpc service functionalities.
package grpcx

import (
	"github.com/gogf/gf/contrib/rpc/grpcx/v3/internal/balancer"
	"github.com/gogf/gf/contrib/rpc/grpcx/v3/internal/grpcctx"
	"github.com/gogf/gf/contrib/rpc/grpcx/v3/internal/resolver"
)

type (
	modCtx      = grpcctx.Ctx
	modBalancer = balancer.Balancer
	modResolver = resolver.Manager
	modClient   struct{}
	modServer   struct{}
)

const (
	FreePortAddress      = ":0" // FreePortAddress marks the server listens using random free port.
	defaultListenAddress = ":0" // Default listening address for grpc server if no address configured.
)

const (
	defaultServerName        = `default`
	configNodeNameGrpcServer = `grpc`
)

var (
	Ctx      = modCtx{}      // Ctx is instance of module Context, which manages the context feature.
	Balancer = modBalancer{} // Balancer is instance of module Balancer, which manages the load balancer features.
	Resolver = modResolver{} // Resolver is instance of module Resolver, which manages the DNS resolving for client.
	Client   = modClient{}   // Client is instance of module Client, which manages the client features.
	Server   = modServer{}   // Server is instance of module Server, which manages the server feature.
)
