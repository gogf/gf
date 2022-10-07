// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package grpcctx

import (
	"context"

	"google.golang.org/grpc/metadata"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gconv"
)

type (
	Ctx struct{}
)

func (c Ctx) NewIncoming(ctx context.Context, data ...g.Map) context.Context {
	if len(data) > 0 {
		incomingMd := make(metadata.MD)
		for key, value := range data[0] {
			incomingMd.Set(key, gconv.String(value))
		}
		return metadata.NewIncomingContext(ctx, incomingMd)
	}
	return metadata.NewIncomingContext(ctx, nil)
}

func (c Ctx) NewOutgoing(ctx context.Context, data ...g.Map) context.Context {
	if len(data) > 0 {
		outgoingMd := make(metadata.MD)
		for key, value := range data[0] {
			outgoingMd.Set(key, gconv.String(value))
		}
		return metadata.NewOutgoingContext(ctx, outgoingMd)
	}
	return metadata.NewOutgoingContext(ctx, nil)
}

func (c Ctx) IncomingToOutgoing(ctx context.Context, keys ...string) context.Context {
	incomingMd, _ := metadata.FromIncomingContext(ctx)
	if incomingMd == nil {
		return ctx
	}
	outgoingMd, _ := metadata.FromOutgoingContext(ctx)
	if outgoingMd == nil {
		outgoingMd = make(metadata.MD)
	}
	if len(keys) > 0 {
		for _, key := range keys {
			outgoingMd[key] = append(outgoingMd[key], incomingMd.Get(key)...)
		}
	} else {
		for key, values := range incomingMd {
			outgoingMd[key] = append(outgoingMd[key], values...)
		}
	}
	return metadata.NewOutgoingContext(ctx, outgoingMd)
}

func (c Ctx) IncomingMap(ctx context.Context) *gmap.Map {
	var (
		data          = gmap.New()
		incomingMd, _ = metadata.FromIncomingContext(ctx)
	)
	for key, values := range incomingMd {
		if len(values) == 1 {
			data.Set(key, values[0])
		} else {
			data.Set(key, values)
		}
	}
	return data
}

func (c Ctx) OutgoingMap(ctx context.Context) *gmap.Map {
	var (
		data          = gmap.New()
		outgoingMd, _ = metadata.FromOutgoingContext(ctx)
	)
	for key, values := range outgoingMd {
		if len(values) == 1 {
			data.Set(key, values[0])
		} else {
			data.Set(key, values)
		}
	}
	return data
}

func (c Ctx) SetIncoming(ctx context.Context, data g.Map) context.Context {
	incomingMd, _ := metadata.FromIncomingContext(ctx)
	if incomingMd == nil {
		incomingMd = make(metadata.MD)
	}
	for key, value := range data {
		incomingMd.Set(key, gconv.String(value))
	}
	return metadata.NewIncomingContext(ctx, incomingMd)
}

func (c Ctx) SetOutgoing(ctx context.Context, data g.Map) context.Context {
	outgoingMd, _ := metadata.FromOutgoingContext(ctx)
	if outgoingMd == nil {
		outgoingMd = make(metadata.MD)
	}
	for key, value := range data {
		outgoingMd.Set(key, gconv.String(value))
	}
	return metadata.NewOutgoingContext(ctx, outgoingMd)
}
