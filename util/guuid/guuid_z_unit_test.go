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
	gtest.C(t, func(t *gtest.T) {
		t.Assert(len(guuid.New().String()), 36)

		uuid, _ := guuid.NewUUID()
		t.Assert(len(uuid.String()), 36)

		uuid, _ = guuid.NewDCEGroup()
		t.Assert(len(uuid.String()), 36)

		uuid, _ = guuid.NewDCEPerson()
		t.Assert(len(uuid.String()), 36)

		uuid, _ = guuid.NewRandom()
		t.Assert(len(uuid.String()), 36)

		t.Assert(len(guuid.NewMD5(guuid.UUID{}, []byte("")).String()), 36)
		t.Assert(len(guuid.NewSHA1(guuid.UUID{}, []byte("")).String()), 36)
	})
}
