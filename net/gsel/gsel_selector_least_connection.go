// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

import (
	"context"
	"sync"

	"github.com/gogf/gf/v3/container/gatomic"
	"github.com/gogf/gf/v3/internal/intlog"
)

type selectorLeastConnection struct {
	mu    sync.RWMutex
	nodes []*leastConnectionNode
}

type leastConnectionNode struct {
	Node
	inflight *gatomic.Int
}

func NewSelectorLeastConnection() Selector {
	return &selectorLeastConnection{
		nodes: make([]*leastConnectionNode, 0),
	}
}

func (s *selectorLeastConnection) Update(ctx context.Context, nodes Nodes) error {
	intlog.Printf(ctx, `Update nodes: %s`, nodes.String())
	var newNodes []*leastConnectionNode
	for _, v := range nodes {
		node := v
		newNodes = append(newNodes, &leastConnectionNode{
			Node:     node,
			inflight: gatomic.NewInt(),
		})
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.nodes = newNodes
	return nil
}

func (s *selectorLeastConnection) Pick(ctx context.Context) (node Node, done DoneFunc, err error) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	var pickedNode *leastConnectionNode
	if len(s.nodes) == 1 {
		pickedNode = s.nodes[0]
	} else {
		for _, v := range s.nodes {
			if pickedNode == nil {
				pickedNode = v
			} else if v.inflight.Val() < pickedNode.inflight.Val() {
				pickedNode = v
			}
		}
	}
	pickedNode.inflight.Add(1)
	done = func(ctx context.Context, di DoneInfo) {
		pickedNode.inflight.Add(-1)
	}
	node = pickedNode.Node
	intlog.Printf(ctx, `Picked node: %s`, node.Address())
	return node, done, nil
}
