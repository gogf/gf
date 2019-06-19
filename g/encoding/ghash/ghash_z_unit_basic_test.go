package ghash_test

import (
	"testing"

	"github.com/gogf/gf/g/encoding/ghash"
	"github.com/gogf/gf/g/test/gtest"
)

var (
	strBasic = []byte("This is the test string for hash.")
)

func Test_BKDRHash(t *testing.T) {
	var x uint32 = 200645773
	gtest.Case(t, func() {
		j := ghash.BKDRHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_BKDRHash64(t *testing.T) {
	var x uint64 = 4214762819217104013
	gtest.Case(t, func() {
		j := ghash.BKDRHash64(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_SDBMHash(t *testing.T) {
	var x uint32 = 1069170245
	gtest.Case(t, func() {
		j := ghash.SDBMHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_SDBMHash64(t *testing.T) {
	var x uint64 = 9881052176572890693
	gtest.Case(t, func() {
		j := ghash.SDBMHash64(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_RSHash(t *testing.T) {
	var x uint32 = 1944033799
	gtest.Case(t, func() {
		j := ghash.RSHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_RSHash64(t *testing.T) {
	var x uint64 = 13439708950444349959
	gtest.Case(t, func() {
		j := ghash.RSHash64(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_JSHash(t *testing.T) {
	var x uint32 = 498688898
	gtest.Case(t, func() {
		j := ghash.JSHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_JSHash64(t *testing.T) {
	var x uint64 = 13410163655098759877
	gtest.Case(t, func() {
		j := ghash.JSHash64(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_PJWHash(t *testing.T) {
	var x uint32 = 7244206
	gtest.Case(t, func() {
		j := ghash.PJWHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_PJWHash64(t *testing.T) {
	var x uint64 = 31150
	gtest.Case(t, func() {
		j := ghash.PJWHash64(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_ELFHash(t *testing.T) {
	var x uint32 = 7244206
	gtest.Case(t, func() {
		j := ghash.ELFHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_ELFHash64(t *testing.T) {
	var x uint64 = 31150
	gtest.Case(t, func() {
		j := ghash.ELFHash64(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_DJBHash(t *testing.T) {
	var x uint32 = 959862602
	gtest.Case(t, func() {
		j := ghash.DJBHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_DJBHash64(t *testing.T) {
	var x uint64 = 2519720351310960458
	gtest.Case(t, func() {
		j := ghash.DJBHash64(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_APHash(t *testing.T) {
	var x uint32 = 3998202516
	gtest.Case(t, func() {
		j := ghash.APHash(strBasic)
		gtest.Assert(j, x)
	})
}

func Test_APHash64(t *testing.T) {
	var x uint64 = 2531023058543352243
	gtest.Case(t, func() {
		j := ghash.APHash64(strBasic)
		gtest.Assert(j, x)
	})
}
