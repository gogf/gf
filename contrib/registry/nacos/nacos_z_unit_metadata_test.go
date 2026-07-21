// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package nacos

import (
	"fmt"
	"sync"
	"testing"

	"github.com/gogf/gf/v2/test/gtest"
)

// Test_DefaultMetadata_Race covers #4649: concurrent SetDefaultMetadata and
// reads (as Register does) must not race under -race.
func Test_DefaultMetadata_Race(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		reg := NewWithClient(nil, WithDefaultMetadata(map[string]string{"init": "1"}))
		var wg sync.WaitGroup
		wg.Add(2)
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				reg.SetDefaultMetadata(map[string]string{"k": fmt.Sprint(i)})
			}
		}()
		go func() {
			defer wg.Done()
			for i := 0; i < 500; i++ {
				m := reg.loadDefaultMetadata()
				_ = m["k"]
				_ = m["init"]
			}
		}()
		wg.Wait()

		// Options apply at construction.
		reg2 := NewWithClient(nil, WithClusterName("C1"), WithGroupName("G1"), WithDefaultEndpoint("1.2.3.4:80"))
		t.Assert(reg2.clusterName, "C1")
		t.Assert(reg2.groupName, "G1")
		t.Assert(reg2.defaultEndpoint, "1.2.3.4:80")
		md := reg2.loadDefaultMetadata()
		t.Assert(len(md), 0)
	})
}
