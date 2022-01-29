// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

import (
	"context"
	"sync"

	"github.com/gogf/gf/v2/util/grand"
)

const SelectorRandom = "BalancerRandom"

type selectorRandom struct {
	mu    sync.RWMutex
	nodes []Node
}

func NewSelectorRandom() Selector {
	return &selectorRandom{
		nodes: make([]Node, 0),
	}
}

func (s *selectorRandom) Update(nodes []Node) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes = nodes
	return nil
}

func (s *selectorRandom) Pick(ctx context.Context) (node Node, done DoneFunc, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if len(s.nodes) == 0 {
		return nil, nil, nil
	}
	return s.nodes[grand.Intn(len(s.nodes))], nil, nil
}
