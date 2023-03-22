// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsel

type builderRoundRobin struct{}

func NewBuilderRoundRobin() Builder {
	return &builderRoundRobin{}
}

func (*builderRoundRobin) Name() string {
	return "BalancerRoundRobin"
}

func (*builderRoundRobin) Build() Selector {
	return NewSelectorRoundRobin()
}
