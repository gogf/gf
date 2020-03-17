// Copyright 2017 gf Author(https://github.com/gogf/gf). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gbase64_test

import (
	"github.com/gogf/gf/debug/gdebug"
	"github.com/gogf/gf/os/gfile"
	"testing"

	"github.com/gogf/gf/encoding/gbase64"
	"github.com/gogf/gf/test/gtest"
)

type testPair struct {
	decoded, encoded string
}

var pairs = []testPair{
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

func Test_Basic(t *testing.T) {
	gtest.Case(t, func() {
		for k := range pairs {
			// Encode
			gtest.Assert(gbase64.Encode([]byte(pairs[k].decoded)), []byte(pairs[k].encoded))
			gtest.Assert(gbase64.EncodeToString([]byte(pairs[k].decoded)), pairs[k].encoded)
			gtest.Assert(gbase64.EncodeString(pairs[k].decoded), pairs[k].encoded)

			// Decode
			r1, _ := gbase64.Decode([]byte(pairs[k].encoded))
			gtest.Assert(r1, []byte(pairs[k].decoded))

			r2, _ := gbase64.DecodeString(pairs[k].encoded)
			gtest.Assert(r2, []byte(pairs[k].decoded))

			r3, _ := gbase64.DecodeToString(pairs[k].encoded)
			gtest.Assert(r3, pairs[k].decoded)
		}
	})
}

func Test_File(t *testing.T) {
	path := gfile.Join(gdebug.TestDataPath(), "test")
	expect := "dGVzdA=="
	gtest.Case(t, func() {
		b, err := gbase64.EncodeFile(path)
		gtest.Assert(err, nil)
		gtest.Assert(string(b), expect)
	})
	gtest.Case(t, func() {
		s, err := gbase64.EncodeFileToString(path)
		gtest.Assert(err, nil)
		gtest.Assert(s, expect)
	})
}

func Test_File_Error(t *testing.T) {
	path := "none-exist-file"
	expect := ""
	gtest.Case(t, func() {
		b, err := gbase64.EncodeFile(path)
		gtest.AssertNE(err, nil)
		gtest.Assert(string(b), expect)
	})
	gtest.Case(t, func() {
		s, err := gbase64.EncodeFileToString(path)
		gtest.AssertNE(err, nil)
		gtest.Assert(s, expect)
	})
}
