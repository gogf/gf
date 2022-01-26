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

const SelectorNameRandom = "BalancerRandom"

type selectorNameRandom struct {
	mu    sync.RWMutex
	nodes []Node
}

func NewSelectorRandom() Selector {
	return &selectorNameRandom{
		nodes: make([]Node, 0),
	}
}

// Update modify Service.
func (p *selectorNameRandom) Update(nodes []Node) error {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.nodes = append(p.nodes, nodes...)
	return nil
}

func (p *selectorNameRandom) Pick(ctx context.Context) (node Node, done DoneFunc, err error) {
	p.mu.RLock()
	defer p.mu.RUnlock()
	return p.nodes[grand.Intn(len(p.nodes))], nil, nil
}
