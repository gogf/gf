// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"testing"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/test/gtest"
	"github.com/gogf/gf/v2/util/grand"
	"github.com/gogf/gf/v2/util/guid"
)

func Test_GroupSortedSet_ZADD(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			maxn int = 100000000
			k1       = guid.S()
			k1m1     = guid.S()
			k1m2     = guid.S()

			option  gredis.ZAddOption
			member1 gredis.ZAddMember
			member2 gredis.ZAddMember
		)

		member1 = gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: k1m1,
		}

		_, err := redis.GroupSortedSet().ZAdd(ctx, k1, &option, member1)
		t.AssertNil(err)

		member2 = gredis.ZAddMember{
			Score:  float64(grand.Intn(1000000)),
			Member: k1m2,
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, &option, member2)
		t.AssertNil(err)

		_, err = redis.GroupSortedSet().ZScore(ctx, k1, k1m1)
		t.AssertNil(err)

		_, err = redis.GroupSortedSet().ZScore(ctx, k1, k1m2)
		t.AssertNil(err)

		var (
			k2   string = guid.S()
			k2m1 string = guid.S()
			k2m2 string = guid.S()
			k2m3 int    = grand.Intn(maxn)
		)

		member3 := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: k2m1,
		}

		member4 := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: k2m2,
		}

		member5 := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: k2m3,
		}

		_, err = redis.GroupSortedSet().ZAdd(ctx, k2, &option, member3, member4, member5)
	})

	// with option
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			maxn int = 100000000
			k1       = guid.S()
			k1m1     = guid.S()
		)

		member1 := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: k1m1,
		}

		option := gredis.ZAddOption{}
		_, err := redis.GroupSortedSet().ZAdd(ctx, k1, &option, member1)
		t.AssertNil(err)

		// option XX
		optionXX := &gredis.ZAddOption{
			XX: true,
		}
		memberXX := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: k1m1,
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionXX, memberXX)
		t.AssertNil(err)

		scoreXX, err := redis.GroupSortedSet().ZScore(ctx, k1, memberXX.Member)
		t.AssertNil(err)
		t.AssertEQ(scoreXX, memberXX.Score)

		// option NX
		optionNX := &gredis.ZAddOption{
			NX: true,
		}
		memberNX := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: guid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionNX, memberNX)
		t.AssertNil(err)

		scoreNXOrigin := memberNX.Score
		memberNX.Score = float64(grand.Intn(maxn))
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionNX, memberNX)
		t.AssertNil(err)

		score, err := redis.GroupSortedSet().ZScore(ctx, k1, memberNX.Member)
		t.AssertNil(err)
		t.AssertEQ(score, scoreNXOrigin)

		// option LT
		optionLT := &gredis.ZAddOption{
			LT: true,
		}
		memberLT := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: guid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionLT, memberLT)
		t.AssertNil(err)

		memberLT.Score += 1
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionLT, memberLT)
		t.AssertNil(err)
		scoreLT, err := redis.GroupSortedSet().ZScore(ctx, k1, memberLT.Member)
		t.AssertLT(scoreLT, memberLT.Score)

		memberLT.Score -= 3
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionLT, memberLT)
		t.AssertNil(err)
		scoreLT, err = redis.GroupSortedSet().ZScore(ctx, k1, memberLT.Member)
		t.AssertEQ(scoreLT, memberLT.Score)

		// option GT
		optionGT := &gredis.ZAddOption{
			GT: true,
		}
		memberGT := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: guid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionGT, memberGT)
		t.AssertNil(err)

		memberLT.Score -= 1
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionGT, memberGT)
		t.AssertNil(err)
		scoreGT, err := redis.GroupSortedSet().ZScore(ctx, k1, memberLT.Member)
		t.AssertGT(scoreGT, memberLT.Score)

		memberLT.Score += 3
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionGT, memberGT)
		t.AssertNil(err)
		scoreGT, err = redis.GroupSortedSet().ZScore(ctx, k1, memberGT.Member)
		t.AssertEQ(scoreGT, memberGT.Score)

		// option CH
		optionCH := &gredis.ZAddOption{
			CH: true,
		}
		memberCH := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: guid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)

		changed, err := redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)
		t.AssertEQ(changed.Val(), int64(0))

		memberCH.Score += 1
		changed, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)
		t.AssertEQ(changed.Val(), int64(1))

		memberCH.Member = guid.S()
		changed, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionCH, memberCH)
		t.AssertNil(err)
		t.AssertEQ(changed.Val(), int64(1))

		// option INCR
		optionINCR := &gredis.ZAddOption{
			INCR: true,
		}
		memberINCR := gredis.ZAddMember{
			Score:  float64(grand.Intn(maxn)),
			Member: guid.S(),
		}
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionINCR, memberINCR)
		t.AssertNil(err)

		scoreIncrOrigin := memberINCR.Score
		memberINCR.Score += 1
		_, err = redis.GroupSortedSet().ZAdd(ctx, k1, optionINCR, memberINCR)
		t.AssertNil(err)

		scoreINCR, err := redis.GroupSortedSet().ZScore(ctx, k1, memberINCR.Member)
		t.AssertNil(err)
		t.AssertEQ(scoreINCR, memberINCR.Score+scoreIncrOrigin)
	})
}

func Test_GroupSortedSet_ZScore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k1   string = guid.S()
			m1   string = guid.S()
			maxn int    = 1000000

			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		member := gredis.ZAddMember{
			Member: m1,
			Score:  float64(grand.Intn(maxn)),
		}

		_, err := redis.GroupSortedSet().ZAdd(ctx, k1, option, member)
		t.AssertNil(err)

		score, err := redis.GroupSortedSet().ZScore(ctx, k1, m1)
		t.AssertNil(err)
		t.AssertEQ(score, member.Score)

		m2 := guid.S()
		score, err = redis.GroupSortedSet().ZScore(ctx, k1, m2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		k2 := guid.S()
		score, err = redis.GroupSortedSet().ZScore(ctx, k2, m2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))
	})
}

func Test_GroupSortedSet_ZIncrBy(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		k := guid.S()
		m := guid.S()

		var incr float64 = 6
		_, err := redis.GroupSortedSet().ZIncrBy(ctx, k, incr, m)
		t.AssertNil(err)

		incr2 := float64(3)
		incredScore, err := redis.GroupSortedSet().ZIncrBy(ctx, k, incr2, m)
		t.AssertNil(err)
		t.AssertEQ(incredScore, incr+incr2)
	})
}

func Test_GroupSortedSet_ZCard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)
		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		rand := grand.N(10, 20)
		for i := 1; i <= rand; i++ {
			_, err := redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(i),
			})
			t.AssertNil(err)

			cnt, err := redis.GroupSortedSet().ZCard(ctx, k)
			t.AssertNil(err)
			t.AssertEQ(cnt, int64(i))
		}
	})
}

func Test_GroupSortedSet_ZCount(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k      string             = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		min, max := "5", "378"
		memSlice := []int{-6, 3, 7, 9, 100, 500, 666}

		for i := 0; i < len(memSlice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: memSlice[i],
				Score:  float64(memSlice[i]),
			})
		}

		cnt, err := redis.GroupSortedSet().ZCount(ctx, k, min, max)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(3))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "-inf", max)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(5))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "-inf", "+inf")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(len(memSlice)))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "(500", "(567")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "(500", "+inf")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(1))

		cnt, err = redis.GroupSortedSet().ZCount(ctx, k, "(3", "(567")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(4))
	})
}

func Test_GroupSortedSet_ZRange(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []string{"one", "two", "three"}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: slice[i],
				Score:  float64(i + 1),
			})
		}

		ret, err := redis.GroupSortedSet().ZRange(ctx, k, 0, -1)
		t.AssertNil(err)
		expected := []string{"one", "two", "three"}
		for i := 0; i < len(ret); i++ {
			t.AssertEQ(ret[i].String(), expected[i])
		}

		ret, err = redis.GroupSortedSet().ZRange(ctx, k, 2, 3)
		t.AssertNil(err)
		expected = []string{"three"}
		for i := 0; i < len(ret); i++ {
			t.AssertEQ(ret[i].String(), expected[i])
		}

		// ret, err = redis.GroupSortedSet().ZRange(ctx, k, 0, -1,
		// 	gredis.ZRangeOption{WithScores: true})
		// t.AssertNil(err)
		// expectedScore := []interface{}{1, "one", 2, "two", 3, "three"}
		// for i := 0; i < len(ret); i++ {
		// 	t.AssertEQ(ret[i].String(), expectedScore[i])
		// }
	})
}

func Test_GroupSortedSet_ZRevRange(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []string{"one", "two", "three"}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: slice[i],
				Score:  float64(i + 1),
			})
		}

		ret, err := redis.GroupSortedSet().ZRevRange(ctx, k, 0, -1)
		t.AssertNil(err)
		expected := []interface{}{"three", "two", "one"}
		t.AssertEQ(ret.Slice(), expected)

		ret, err = redis.GroupSortedSet().ZRevRange(ctx, k, 0, 1)
		t.AssertNil(err)
		expected = []interface{}{"three", "two"}
		t.AssertEQ(ret.Slice(), expected)
	})
}

func Test_GroupSortedSet_ZRank(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rank, err := redis.ZRank(ctx, k, 0)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))

		rank, err = redis.ZRank(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(3))

		rank, err = redis.ZRank(ctx, k, 6)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))
	})
}

func Test_GroupSortedSet_ZRevRank(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rank, err := redis.ZRevRank(ctx, k, 0)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(4))

		rank, err = redis.ZRevRank(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(1))

		rank, err = redis.ZRevRank(ctx, k, 9)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))

		rank, err = redis.ZRevRank(ctx, k, 6)
		t.AssertNil(err)
		t.AssertEQ(rank, int64(0))
	})
}

func Test_GroupSortedSet_ZRem(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k      = guid.S()
			m      = guid.S()
			option = new(gredis.ZAddOption)
		)

		cnt, err := redis.ZRem(ctx, k, m)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))

		member := gredis.ZAddMember{
			Member: m,
			Score:  123,
		}
		redis.ZAdd(ctx, k, option, member)

		cnt, err = redis.ZRem(ctx, k, m)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(1))

		member2 := gredis.ZAddMember{
			Member: guid.S(),
			Score:  456,
		}
		_, err = redis.ZAdd(ctx, k, option, member, member2)
		t.AssertNil(err)

		cnt, err = redis.ZRem(ctx, k, m, "non_exists")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(1))

		_, err = redis.ZAdd(ctx, k, option, member, member2)
		t.AssertNil(err)

		cnt, err = redis.ZRem(ctx, k, m, member2.Member)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(2))
	})
}

func Test_GroupSortedSet_ZRemRangeByRank(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByRank(ctx, k, 0, 2)
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(3))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 0)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(7))
	})

	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByRank(ctx, k, -3, -2)
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(2))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(9))
	})

	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByRank(ctx, k, 3, -1)
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(2))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(3))
	})
}
func Test_GroupSortedSet_ZRemRangeByScore(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "(3", "9")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(3))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 5)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(3))
	})

	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "3", "(9")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(3))

		score, err := redis.GroupSortedSet().ZScore(ctx, k, 1)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 2)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 3)
		t.AssertNil(err)
		t.AssertEQ(score, float64(0))

		score, err = redis.GroupSortedSet().ZScore(ctx, k, 4)
		t.AssertNil(err)
		t.AssertEQ(score, float64(9))
	})

	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "-inf", "9")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(5))

		cnt, err := redis.GroupSortedSet().ZCard(ctx, k)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))
	})

	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []int64{1, 3, 5, 7, 9}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: i,
				Score:  float64(slice[i]),
			})
		}

		rmd, err := redis.GroupSortedSet().ZRemRangeByScore(ctx, k, "-inf", "+inf")
		t.AssertNil(err)
		t.AssertEQ(rmd, int64(5))

		cnt, err := redis.GroupSortedSet().ZCard(ctx, k)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(0))
	})
}
func Test_GroupSortedSet_ZRemRangeByLex(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []string{"aaaa", "b", "c", "d", "e", "foo", "zap", "zip", "ALPHA", "alpha"}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: slice[i],
				Score:  float64(0),
			})
		}

		cnt, err := redis.GroupSortedSet().ZRemRangeByLex(ctx, k, "[alpha", "[omega")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(6))

		cnt, err = redis.GroupSortedSet().ZCard(ctx, k)
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(4))
	})
}
func Test_GroupSortedSet_ZLexCount(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		defer redis.FlushDB(ctx)

		var (
			k                         = guid.S()
			option *gredis.ZAddOption = new(gredis.ZAddOption)
		)

		slice := []string{"a", "b", "c", "d", "e", "f", "g"}
		for i := 0; i < len(slice); i++ {
			redis.GroupSortedSet().ZAdd(ctx, k, option, gredis.ZAddMember{
				Member: slice[i],
				Score:  float64(0),
			})
		}

		cnt, err := redis.GroupSortedSet().ZLexCount(ctx, k, "-", "+")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(7))

		cnt, err = redis.GroupSortedSet().ZLexCount(ctx, k, "[b", "[f")
		t.AssertNil(err)
		t.AssertEQ(cnt, int64(5))

	})
}
