/*
// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
*/
package gtype_test

import (
	"testing"

	"github.com/gogf/gf/v2/container/gtype"
	"github.com/gogf/gf/v2/internal/json"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/shopspring/decimal"
	"golang.org/x/text/number"
)

func Test_Decimal(t *testing.T) {
	var d = number.Decimal(11111)
	t.Log(d)
}

func Test_Decimal_JSON(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		v := decimal.Zero
		i := gtype.NewDecimal(v)
		b1, err1 := json.Marshal(i)
		b2, err2 := json.Marshal(i.Val())
		t.Assert(err1, nil)
		t.Assert(err2, nil)
		t.Assert(b1, b2)

		i2 := gtype.NewDecimal()
		err := json.Unmarshal(b2, &i2)
		t.Assert(err, nil)
		t.Assert(i2.Val(), v)
	})
}
