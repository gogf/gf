// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

// Unit tests for the builders and selectors in package gsel.
package gsel_test

import (
	"context"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/net/gsel"
	"github.com/gogf/gf/v2/net/gsvc"
	"github.com/gogf/gf/v2/test/gtest"
)

// mockNode implements gsel.Node for testing purposes.
type mockNode struct {
	address string
	service gsvc.Service
}

func (n *mockNode) Service() gsvc.Service { return n.service }
func (n *mockNode) Address() string       { return n.address }

// newNode creates a plain test node without weight metadata.
func newNode(address string) gsel.Node {
	svc := gsvc.NewServiceWithName("test-svc")
	return &mockNode{address: address, service: svc}
}

// newWeightedNode creates a test node with the given weight in its metadata.
func newWeightedNode(address string, weight int) gsel.Node {
	svc := &gsvc.LocalService{
		Name:      "test-svc",
		Endpoints: gsvc.NewEndpoints(address),
		Metadata:  gsvc.Metadata{gsvc.MDWeight: weight},
	}
	return &mockNode{address: address, service: svc}
}

// ───────────────────────────── Builder ─────────────────────────────

func TestBuilder_SetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		previous := gsel.GetBuilder()
		t.Cleanup(func() { gsel.SetBuilder(previous) })

		// default is round-robin
		t.Assert(gsel.GetBuilder().Name(), "BalancerRoundRobin")

		// switch to random
		gsel.SetBuilder(gsel.NewBuilderRandom())
		t.Assert(gsel.GetBuilder().Name(), "BalancerRandom")

		// switch to least-connection
		gsel.SetBuilder(gsel.NewBuilderLeastConnection())
		t.Assert(gsel.GetBuilder().Name(), "BalancerLeastConnection")

		// switch to weight
		gsel.SetBuilder(gsel.NewBuilderWeight())
		t.Assert(gsel.GetBuilder().Name(), "BalancerWeight")

		// restore default
		gsel.SetBuilder(gsel.NewBuilderRoundRobin())
		t.Assert(gsel.GetBuilder().Name(), "BalancerRoundRobin")
	})
}

func TestBuilder_Build(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		builders := []gsel.Builder{
			gsel.NewBuilderRoundRobin(),
			gsel.NewBuilderRandom(),
			gsel.NewBuilderLeastConnection(),
			gsel.NewBuilderWeight(),
		}
		for _, b := range builders {
			sel := b.Build()
			t.AssertNE(sel, nil)
		}
	})
}

// ───────────────────────────── Nodes.String ────────────────────────

func TestNodes_String(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		nodes := gsel.Nodes{newNode("1.1.1.1:8000"), newNode("2.2.2.2:9000")}
		s := nodes.String()
		t.Assert(s, "1.1.1.1:8000,2.2.2.2:9000")
	})
	gtest.C(t, func(t *gtest.T) {
		var nodes gsel.Nodes
		t.Assert(nodes.String(), "")
	})
}

// ─────────────────────────── RoundRobin ───────────────────────────

func TestSelectorRoundRobin_StrictOrder(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorRoundRobin()
		nodes := gsel.Nodes{newNode("a:1"), newNode("b:2"), newNode("c:3")}
		t.AssertNil(sel.Update(ctx, nodes))

		got := make([]string, 9)
		for i := range got {
			node, _, err := sel.Pick(ctx)
			t.AssertNil(err)
			got[i] = node.Address()
		}
		// must cycle strictly: a b c a b c a b c
		t.Assert(got, []string{"a:1", "b:2", "c:3", "a:1", "b:2", "c:3", "a:1", "b:2", "c:3"})
	})
}

func TestSelectorRoundRobin_UpdateResetsCounter(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorRoundRobin()
		nodes := gsel.Nodes{newNode("a:1"), newNode("b:2")}
		t.AssertNil(sel.Update(ctx, nodes))

		// advance the cursor once only, so that it is left at 1
		_, _, err := sel.Pick(ctx) // a:1
		t.AssertNil(err)

		// update with new nodes — counter must reset to 0
		nodes2 := gsel.Nodes{newNode("x:1"), newNode("y:2")}
		t.AssertNil(sel.Update(ctx, nodes2))
		node, _, err := sel.Pick(ctx)
		t.AssertNil(err)
		// without the reset the cursor stays at 1 and "y:2" would be picked
		t.Assert(node.Address(), "x:1")
	})
}

func TestSelectorRoundRobin_EmptyNodes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorRoundRobin()
		t.AssertNil(sel.Update(ctx, gsel.Nodes{}))
		node, done, err := sel.Pick(ctx)
		t.AssertNil(err)
		t.AssertNil(node)
		t.AssertNil(done)
	})
}

func TestSelectorRoundRobin_SingleNode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorRoundRobin()
		t.AssertNil(sel.Update(ctx, gsel.Nodes{newNode("only:1")}))
		for i := 0; i < 5; i++ {
			node, _, err := sel.Pick(ctx)
			t.AssertNil(err)
			t.Assert(node.Address(), "only:1")
		}
	})
}

// ─────────────────────────── Random ───────────────────────────────

func TestSelectorRandom_AllNodesReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorRandom()
		nodes := gsel.Nodes{newNode("a:1"), newNode("b:2"), newNode("c:3")}
		t.AssertNil(sel.Update(ctx, nodes))

		seen := make(map[string]bool)
		for i := 0; i < 300; i++ {
			node, _, err := sel.Pick(ctx)
			t.AssertNil(err)
			seen[node.Address()] = true
		}
		t.Assert(len(seen), 3)
	})
}

func TestSelectorRandom_EmptyNodes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorRandom()
		t.AssertNil(sel.Update(ctx, gsel.Nodes{}))
		node, done, err := sel.Pick(ctx)
		t.AssertNil(err)
		t.AssertNil(node)
		t.AssertNil(done)
	})
}

func TestSelectorRandom_SingleNode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorRandom()
		t.AssertNil(sel.Update(ctx, gsel.Nodes{newNode("only:1")}))
		for i := 0; i < 10; i++ {
			node, _, err := sel.Pick(ctx)
			t.AssertNil(err)
			t.Assert(node.Address(), "only:1")
		}
	})
}

// ─────────────────────── LeastConnection ──────────────────────────

func TestSelectorLeastConnection_PicksLeastInflight(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorLeastConnection()
		nodes := gsel.Nodes{newNode("a:1"), newNode("b:2"), newNode("c:3")}
		t.AssertNil(sel.Update(ctx, nodes))

		// pick a:1 and hold (don't call done)
		n1, _, err := sel.Pick(ctx) // a:1 inflight=1
		t.AssertNil(err)
		t.Assert(n1.Address(), "a:1")

		// next two picks should choose b:2 and c:3 (inflight=0)
		n2, _, err := sel.Pick(ctx)
		t.AssertNil(err)
		n3, _, err := sel.Pick(ctx)
		t.AssertNil(err)
		addrs := []string{n2.Address(), n3.Address()}
		t.AssertIN("b:2", addrs)
		t.AssertIN("c:3", addrs)
	})
}

func TestSelectorLeastConnection_DoneDecrementsInflight(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorLeastConnection()
		t.AssertNil(sel.Update(ctx, gsel.Nodes{newNode("a:1"), newNode("b:2")}))

		// two nodes: the first pick takes a:1 (inflight 1); the second finds b:2
		// less busy (inflight 0); the third is a tie at 1 and falls back to a:1.
		n1, done1, err := sel.Pick(ctx) // a:1 inflight->1
		t.AssertNil(err)
		t.Assert(n1.Address(), "a:1")
		n2, _, err := sel.Pick(ctx) // b:2 inflight->1
		t.AssertNil(err)
		t.Assert(n2.Address(), "b:2")
		n3, _, err := sel.Pick(ctx) // tie 1:1, first node wins: a:1 inflight->2
		t.AssertNil(err)
		t.Assert(n3.Address(), "a:1")

		// releasing one pick of a:1 brings it back to inflight 1, a tie with b:2.
		// Only if done really decrements is a:1 (the first node) picked again;
		// otherwise a:1 stays at 2 and b:2 (inflight 1) is picked.
		done1(ctx, gsel.DoneInfo{})
		node, done4, err := sel.Pick(ctx)
		t.AssertNil(err)
		t.Assert(node.Address(), "a:1")
		done4(ctx, gsel.DoneInfo{})
	})
}

func TestSelectorLeastConnection_ConcurrentSafe(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorLeastConnection()
		var ns gsel.Nodes
		for i := 0; i < 5; i++ {
			ns = append(ns, newNode("node:"+string(rune('0'+i))))
		}
		t.AssertNil(sel.Update(ctx, ns))

		var wg sync.WaitGroup
		for i := 0; i < 100; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()
				_, done, err := sel.Pick(ctx)
				if err == nil && done != nil {
					done(ctx, gsel.DoneInfo{})
				}
			}()
		}
		wg.Wait()
		// if we reach here without a race or panic, concurrent safety holds
	})
}

func TestSelectorLeastConnection_InflightZeroAfterAllDone(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorLeastConnection()
		// two nodes, with a:1 busier than b:2 so that a non-decrementing done
		// would leave a:1 at inflight 2 and keep preferring b:2.
		t.AssertNil(sel.Update(ctx, gsel.Nodes{newNode("a:1"), newNode("b:2")}))

		dones := make([]gsel.DoneFunc, 0, 3)
		for _, expect := range []string{"a:1", "b:2", "a:1"} {
			node, done, err := sel.Pick(ctx)
			t.AssertNil(err)
			t.Assert(node.Address(), expect)
			dones = append(dones, done)
		}
		for _, d := range dones {
			d(ctx, gsel.DoneInfo{})
		}
		// both inflight back to 0 → tie → the first node is selected again.
		node, done, err := sel.Pick(ctx)
		t.AssertNil(err)
		t.Assert(node.Address(), "a:1")
		done(ctx, gsel.DoneInfo{})
	})
}

// ─────────────────────────── Weight ───────────────────────────────

func TestSelectorWeight_ProportionalDistribution(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorWeight()
		nodes := gsel.Nodes{
			newWeightedNode("heavy:1", 3),
			newWeightedNode("light:2", 1),
		}
		t.AssertNil(sel.Update(ctx, nodes))

		counts := make(map[string]int)
		const total = 4000
		for i := 0; i < total; i++ {
			node, _, err := sel.Pick(ctx)
			t.AssertNil(err)
			counts[node.Address()]++
		}
		// heavy should be ~75% ± 5%
		heavyRatio := float64(counts["heavy:1"]) / float64(total)
		t.AssertGT(heavyRatio, 0.70)
		t.AssertLT(heavyRatio, 0.80)
	})
}

func TestSelectorWeight_AllNodesReachable(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorWeight()
		nodes := gsel.Nodes{
			newWeightedNode("a:1", 1),
			newWeightedNode("b:2", 2),
			newWeightedNode("c:3", 3),
		}
		t.AssertNil(sel.Update(ctx, nodes))

		seen := make(map[string]bool)
		for i := 0; i < 600; i++ {
			node, _, err := sel.Pick(ctx)
			t.AssertNil(err)
			t.AssertNE(node, nil)
			seen[node.Address()] = true
		}
		t.Assert(len(seen), 3)
	})
}

func TestSelectorWeight_EmptyNodes(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorWeight()
		t.AssertNil(sel.Update(ctx, gsel.Nodes{}))
		node, done, err := sel.Pick(ctx)
		t.AssertNil(err)
		t.AssertNil(node)
		t.AssertNil(done)
	})
}

func TestSelectorWeight_ZeroWeightExcludesNode(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		ctx := context.Background()
		sel := gsel.NewSelectorWeight()
		nodes := gsel.Nodes{
			newWeightedNode("good:1", 2),
			newWeightedNode("zero:2", 0),
		}
		t.AssertNil(sel.Update(ctx, nodes))

		// weight=0 is never expanded, so only "good:1" should appear
		for i := 0; i < 50; i++ {
			node, _, err := sel.Pick(ctx)
			t.AssertNil(err)
			t.Assert(node.Address(), "good:1")
		}
	})
}
