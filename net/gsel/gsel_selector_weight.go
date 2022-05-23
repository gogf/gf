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
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/grand"
)

const SelectorWeight = "BalancerWeight"

type selectorWeight struct {
	mu    sync.RWMutex
	nodes Nodes
}

func NewSelectorWeight() Selector {
	return &selectorWeight{
		nodes: make(Nodes, 0),
	}
}

func (s *selectorWeight) Update(ctx context.Context, nodes Nodes) error {
	intlog.Printf(ctx, `Update nodes: %s`, nodes.String())
	var newNodes []Node
	for _, v := range nodes {
		node := v
		for i := 0; i < s.getWeight(node); i++ {
			newNodes = append(newNodes, node)
		}
	}
	s.mu.Lock()
	s.nodes = newNodes
	s.mu.Unlock()
	return nil
}

func (s *selectorWeight) Pick(ctx context.Context) (node Node, done DoneFunc, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.nodes) == 0 {
		return nil, nil, nil
	}
	node = s.nodes[grand.Intn(len(s.nodes))]
	intlog.Printf(ctx, `Picked node: %s`, node.Address())
	return node, nil, nil
}

func (s *selectorWeight) getWeight(node Node) int {
	return node.Service().GetMetadata().Get(gsvc.MDWeight).Int()
}
