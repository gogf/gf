// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gvalid_test

import (
	"context"
	"testing"

	"github.com/gogf/gf/v2/util/gvalid"
)

type Foo struct {
	Bar *Bar `p:"bar" v:"required-without:Baz"`
	Baz *Baz `p:"baz" v:"required-without:Bar"`
}
type Bar struct {
	BarKey string `p:"bar_key" v:"required"`
}
type Baz struct {
	BazKey string `p:"baz_key" v:"required"`
}

// https://github.com/gogf/gf/issues/2503
func Test_Issue2503(t *testing.T) {
	foo := &Foo{
		Bar: &Bar{BarKey: "value"},
	}
	err := gvalid.New().Data(foo).Run(context.Background())
	if err != nil {
		t.Fatal(err)
	}
}
