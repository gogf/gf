// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Package grpcx provides grpc service functionalities.
package grpcx

import (
	"time"

	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/balancer"
	"github.com/gogf/gf/contrib/rpc/grpcx/v2/internal/grpcctx"
)

type (
	modCtx      = grpcctx.Ctx
	modBalancer = balancer.Balancer
	modClient   struct{}
	modServer   struct{}
)

const (
	// FreePortAddress marks the server listens using random free port.
	FreePortAddress = ":0"
)

const (
	defaultServerName        = `default`
	defaultTimeout           = 5 * time.Second
	configNodeNameRegistry   = `registry`
	configNodeNameGrpcServer = `grpcserver`
)

var (
	Ctx      = modCtx{}      // Ctx is instance of module Context, which manages the context feature.
	Balancer = modBalancer{} // Balancer is instance of module Balancer, which manages the load balancer features.
	Client   = modClient{}   // Client is instance of module Client, which manages the client features.
	Server   = modServer{}   // Server is instance of module Server, which manages the server feature.
)
