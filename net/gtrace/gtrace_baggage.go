// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gtrace

import (
	"context"
	"github.com/gogf/gf/container/gmap"
	"github.com/gogf/gf/container/gvar"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
)

// Baggage holds the data through all tracing spans.
type Baggage struct {
	ctx context.Context
}

// NewBaggage creates and returns a new Baggage object from given tracing context.
func NewBaggage(ctx context.Context) *Baggage {
	if ctx == nil {
		ctx = context.Background()
	}
	return &Baggage{
		ctx: ctx,
	}
}

// Ctx returns the context that Baggage holds.
func (b *Baggage) Ctx() context.Context {
	return b.ctx
}

// SetValue is a convenient function for adding one key-value pair to baggage.
// Note that it uses attribute.Any to set the key-value pair.
func (b *Baggage) SetValue(key string, value interface{}) context.Context {
	b.ctx = baggage.ContextWithValues(b.ctx, attribute.Any(key, value))
	return b.ctx
}

// SetMap is a convenient function for adding map key-value pairs to baggage.
// Note that it uses attribute.Any to set the key-value pair.
func (b *Baggage) SetMap(data map[string]interface{}) context.Context {
	pairs := make([]attribute.KeyValue, 0)
	for k, v := range data {
		pairs = append(pairs, attribute.Any(k, v))
	}
	b.ctx = baggage.ContextWithValues(b.ctx, pairs...)
	return b.ctx
}

// GetMap retrieves and returns the baggage values as map.
func (b *Baggage) GetMap() *gmap.StrAnyMap {
	m := gmap.NewStrAnyMap()
	set := baggage.Set(b.ctx)
	if length := set.Len(); length > 0 {
		if length == 0 {
			return m
		}
		inter := set.Iter()
		for inter.Next() {
			m.Set(string(inter.Label().Key), inter.Label().Value.AsInterface())
		}
	}
	return m
}

// GetVar retrieves value and returns a *gvar.Var for specified key from baggage.
func (b *Baggage) GetVar(key string) *gvar.Var {
	value := baggage.Value(b.ctx, attribute.Key(key))
	return gvar.New(value.AsInterface())
}
