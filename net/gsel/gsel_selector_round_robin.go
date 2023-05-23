// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/internal/intlog"
)

type selectorRoundRobin struct {
	mu    sync.RWMutex
	nodes Nodes
	next  int
}

func NewSelectorRoundRobin() Selector {
	return &selectorRoundRobin{
		nodes: make(Nodes, 0),
	}
}

func (s *selectorRoundRobin) Update(ctx context.Context, nodes Nodes) error {
	intlog.Printf(ctx, `Update nodes: %s`, nodes.String())
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
	intlog.Printf(ctx, `Picked node: %s`, node.Address())
	return
}
