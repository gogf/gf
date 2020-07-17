// Copyright 2018 gf Author(https://github.com/jin502437344/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/jin502437344/gf.

package gcron_test

import (
	"testing"

	"github.com/jin502437344/gf/os/gcron"
)

func Benchmark_Add(b *testing.B) {
	for i := 0; i < b.N; i++ {
		gcron.Add("1 1 1 1 1 1", func() {

		})
	}
}
