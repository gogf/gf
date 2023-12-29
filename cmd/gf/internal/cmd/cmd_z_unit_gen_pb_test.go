// Copyright GoFrame gf Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package cmd

import (
	"github.com/gogf/gf/cmd/gf/v2/internal/cmd/genpb"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/gutil"
	"testing"
)

func Test_Gen_Pb_Default(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			//path      = gfile.Temp(guid.S())
			in = genpb.CGenPbInput{
				Path:       gtest.DataPath("genpb", "protobuf"),
				OutputApi:  gtest.DataPath("genpb", "api"),
				OutputCtrl: gtest.DataPath("genpb", "controller"),
			}
		)
		err := gutil.FillStructWithDefault(&in)
		t.AssertNil(err)

		_, err = genpb.CGenPb{}.Pb(ctx, in)
		if err != nil {
			panic(err)
		}
	})
}
