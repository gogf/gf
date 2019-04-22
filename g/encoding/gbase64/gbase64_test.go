// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.
package gbase64_test

import (
	"github.com/gogf/gf/g/encoding/gbase64"
	"github.com/gogf/gf/g/test/gtest"
	"testing"
)

type testpair struct {
	decoded, encoded string
}

var pairs = []testpair{
	// RFC 3548 examples
	{"\x14\xfb\x9c\x03\xd9\x7e", "FPucA9l+"},
	{"\x14\xfb\x9c\x03\xd9", "FPucA9k="},
	{"\x14\xfb\x9c\x03", "FPucAw=="},

	// RFC 4648 examples
	{"", ""},
	{"f", "Zg=="},
	{"fo", "Zm8="},
	{"foo", "Zm9v"},
	{"foob", "Zm9vYg=="},
	{"fooba", "Zm9vYmE="},
	{"foobar", "Zm9vYmFy"},

	// Wikipedia examples
	{"sure.", "c3VyZS4="},
	{"sure", "c3VyZQ=="},
	{"sur", "c3Vy"},
	{"su", "c3U="},
	{"leasure.", "bGVhc3VyZS4="},
	{"easure.", "ZWFzdXJlLg=="},
	{"asure.", "YXN1cmUu"},
	{"sure.", "c3VyZS4="},
}

func TestBase64(t *testing.T) {
	for k := range pairs{
		gtest.Assert(gbase64.Encode(pairs[k].decoded), pairs[k].encoded)

		e, _ := gbase64.Decode(pairs[k].encoded)
		gtest.Assert(e, pairs[k].decoded)
	}
}
