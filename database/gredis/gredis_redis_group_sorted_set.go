// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gredis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
)

// IGroupSortedSet manages redis sorted set operations.
// Implements see redis.GroupSortedSet.
type IGroupSortedSet interface {
	ZAdd(ctx context.Context, key string, option *ZAddOption, member ZAddMember, members ...ZAddMember) (*gvar.Var, error)
	ZScore(ctx context.Context, key string, member interface{}) (float64, error)
	ZIncrBy(ctx context.Context, key string, increment float64, member interface{}) (float64, error)
	ZCard(ctx context.Context, key string) (int64, error)
	ZCount(ctx context.Context, key string, min, max string) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64, option ...ZRangeOption) (gvar.Vars, error)
	ZRevRange(ctx context.Context, key string, start, stop int64, option ...ZRevRangeOption) (*gvar.Var, error)
	ZRank(ctx context.Context, key string, member interface{}) (int64, error)
	ZRevRank(ctx context.Context, key string, member interface{}) (int64, error)
	ZRem(ctx context.Context, key string, member interface{}, members ...interface{}) (int64, error)
	ZRemRangeByRank(ctx context.Context, key string, start, stop int64) (int64, error)
	ZRemRangeByScore(ctx context.Context, key string, min, max string) (int64, error)
	ZRemRangeByLex(ctx context.Context, key string, min, max string) (int64, error)
	ZLexCount(ctx context.Context, key, min, max string) (int64, error)
}

// ZAddOption provides options for function ZAdd.
type ZAddOption struct {
	XX bool // Only update elements that already exist. Don't add new elements.
	NX bool // Only add new elements. Don't update already existing elements.
	// Only update existing elements if the new score is less than the current score.
	// This flag doesn't prevent adding new elements.
	LT bool

	// Only update existing elements if the new score is greater than the current score.
	// This flag doesn't prevent adding new elements.
	GT bool

	// Modify the return value from the number of new elements added, to the total number of elements changed (CH is an abbreviation of changed).
	// Changed elements are new elements added and elements already existing for which the score was updated.
	// So elements specified in the command line having the same score as they had in the past are not counted.
	// Note: normally the return value of ZAdd only counts the number of new elements added.
	CH bool

	// When this option is specified ZAdd acts like ZIncrBy. Only one score-element pair can be specified in this mode.
	INCR bool
}

// ZAddMember is element struct for set.
type ZAddMember struct {
	Score  float64
	Member interface{}
}

// ZRangeOption provides extra option for ZRange function.
type ZRangeOption struct {
	ByScore bool
	ByLex   bool
	// The optional REV argument reverses the ordering, so elements are ordered from highest to lowest score,
	// and score ties are resolved by reverse lexicographical ordering.
	Rev   bool
	Limit *ZRangeOptionLimit
	// The optional WithScores argument supplements the command's reply with the scores of elements returned.
	WithScores bool
}

// ZRangeOptionLimit provides LIMIT argument for ZRange function.
// The optional LIMIT argument can be used to obtain a sub-range from the matching elements
// (similar to SELECT LIMIT offset, count in SQL). A negative `Count` returns all elements from the `Offset`.
type ZRangeOptionLimit struct {
	Offset *int
	Count  *int
}

// ZRevRangeOption provides options for function ZRevRange.
type ZRevRangeOption struct {
	WithScores bool
}
