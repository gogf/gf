// Copyright 2017 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gjson_test

import (
	"testing"

	"github.com/jin502437344/gf/encoding/gjson"
)

var (
	jsonStr1 = `[1,2,3]`
	jsonStr2 = `{"CallbackCommand":"Group.CallbackAfterSendMsg","From_Account":"61934946","GroupId":"@TGS#2FLGX67FD","MsgBody":[{"MsgContent":{"Text":"是的"},"MsgType":"TIMTextElem"}],"MsgSeq":23,"MsgTime":1567032819,"Operator_Account":"61934946","Random":2804799576,"Type":"Public"}`
)

func Benchmark_Validate1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.Valid(jsonStr1)
	}
}

func Benchmark_Validate2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gjson.Valid(jsonStr2)
	}
}

func Benchmark_Set1(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := gjson.New(map[string]string{
			"k1": "v1",
			"k2": "v2",
		})
		p.Set("k1.k11", []int{1, 2, 3})
	}
}

func Benchmark_Set2(b *testing.B) {
	for i := 0; i < b.N; i++ {
		p := gjson.New([]string{"a"})
		p.Set("0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0.0", []int{1, 2, 3})
	}
}
