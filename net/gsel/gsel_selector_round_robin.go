// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/encoding/gjson"
	"github.com/gogf/gf/v2/internal/intlog"
)

const SelectorRoundRobin = "BalancerRoundRobin"

type selectorRoundRobin struct {
	mu    sync.RWMutex
	nodes []Node
	next  int
}

func NewSelectorRoundRobin() Selector {
	return &selectorRoundRobin{
		nodes: make([]Node, 0),
	}
}

func (s *selectorRoundRobin) Update(nodes []Node) error {
	intlog.Printf(context.Background(), `Update nodes: %s`, gjson.MustEncode(nodes))
	s.mu.Lock()
	s.nodes = nodes
	s.mu.Unlock()
	return nil
}

func (s *selectorRoundRobin) Pick(ctx context.Context) (node Node, done DoneFunc, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	node = s.nodes[s.next]
	s.next = (s.next + 1) % len(s.nodes)
	intlog.Printf(ctx, `Pick node: %s`, node.Address())
	return
}
