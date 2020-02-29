// Copyright 2019 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package guuid_test

import (
	"github.com/gogf/gf/util/guuid"
	"testing"

	"github.com/gogf/gf/test/gtest"
)

func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		gtest.Assert(len(guuid.New().String()), 36)

		uuid, _ := guuid.NewUUID()
		gtest.Assert(len(uuid.String()), 36)

		uuid, _ = guuid.NewDCEGroup()
		gtest.Assert(len(uuid.String()), 36)

		uuid, _ = guuid.NewDCEPerson()
		gtest.Assert(len(uuid.String()), 36)

		uuid, _ = guuid.NewRandom()
		gtest.Assert(len(uuid.String()), 36)

		gtest.Assert(len(guuid.NewMD5(guuid.UUID{}, []byte("")).String()), 36)
		gtest.Assert(len(guuid.NewSHA1(guuid.UUID{}, []byte("")).String()), 36)
	})
}
