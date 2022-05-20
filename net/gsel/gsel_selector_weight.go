// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/util/grand"
)

const SelectorWeight = "BalancerWeight"

type selectorWeight struct {
	mu    sync.RWMutex
	nodes []Node
}

func NewSelectorWeight() Selector {
	return &selectorWeight{
		nodes: make([]Node, 0),
	}
}

func (s *selectorWeight) Update(nodes []Node) error {
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
	return s.nodes[grand.Intn(len(s.nodes))], nil, nil
}

func (s *selectorWeight) getWeight(node Node) int {
	return node.Service().GetMetadata().Get(gsvc.MDWeight).Int()
}
