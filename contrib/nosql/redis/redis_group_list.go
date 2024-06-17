// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis

import (
	"context"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/util/gconv"
)

// GroupList is the redis group list object.
type GroupList struct {
	Operation gredis.AdapterOperation
}

// GroupList creates and returns a redis group object for list operations.
func (r *Redis) GroupList() gredis.IGroupList {
	return GroupList{
		Operation: r.AdapterOperation,
	}
}

// LPush inserts all the specified values at the head of the list stored at key
// Insert all the specified values at the head of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operations.
// When key holds a value that is not a list, an error is returned.
//
// It returns the length of the list after the push operations.
//
// https://redis.io/commands/lpush/
func (r GroupList) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.Operation.Do(ctx, "LPush", append([]interface{}{key}, values...)...)
	return v.Int64(), err
}

// LPushX insert value at the head of the list stored at key, only if key exists and holds a list.
// Inserts specified values at the head of the list stored at key, only if key already exists and holds a list.
// In contrary to LPush, no operation will be performed when key does not yet exist.
// Return Integer reply: the length of the list after the push operation.
//
// It returns the length of the list after the push operations.
//
// https://redis.io/commands/lpushx
func (r GroupList) LPushX(ctx context.Context, key string, element interface{}, elements ...interface{}) (int64, error) {
	v, err := r.Operation.Do(ctx, "LPushX", append([]interface{}{key, element}, elements...)...)
	return v.Int64(), err
}

// RPush inserts all the specified values at the tail of the list stored at key.
// Insert all the specified values at the tail of the list stored at key.
// If key does not exist, it is created as empty list before performing the push operation.
//
// When key holds a value that is not a list, an error is returned.
// It is possible to push multiple elements using a single command call just specifying multiple
// arguments at  the end of the command. Elements are inserted one after the other to the tail of the
// list, from the leftmost element to the rightmost element.
// So for instance the command RPush mylist a b c will result into a list containing a as first element,
// b as second element and c as third element.
//
// It returns the length of the list after the push operation.
//
// https://redis.io/commands/rpush
func (r GroupList) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	v, err := r.Operation.Do(ctx, "RPush", append([]interface{}{key}, values...)...)
	return v.Int64(), err
}

// RPushX inserts value at the tail of the list stored at key, only if key exists and holds a list.
// Inserts specified values at the tail of the list stored at key, only if key already exists and
// holds a list.
//
// In contrary to RPush, no operation will be performed when key does not yet exist.
//
// It returns the length of the list after the push operation.
//
// https://redis.io/commands/rpushx
func (r GroupList) RPushX(ctx context.Context, key string, value interface{}) (int64, error) {
	v, err := r.Operation.Do(ctx, "RPushX", key, value)
	return v.Int64(), err
}

// LPop remove and returns the first element of the list stored at key.
// Removes and returns the first elements of the list stored at key.
//
// Starting with Redis version 6.2.0: Added the count argument.
//
// By default, the command pops a single element from the beginning of the list.
// When provided with the optional count argument, the reply will consist of up to count elements,
// depending on the list's length.
//
// Return When called without the count argument:
// Bulk string reply: the value of the first element, or nil when key does not exist.
//
// When called with the count argument:
// Array reply: list of popped elements, or nil when key does not exist.
//
// https://redis.io/commands/lpop
func (r GroupList) LPop(ctx context.Context, key string, count ...int) (*gvar.Var, error) {
	if len(count) > 0 {
		return r.Operation.Do(ctx, "LPop", key, count[0])
	}
	return r.Operation.Do(ctx, "LPop", key)
}

// RPop remove and returns the last element of the list stored at key.
// Removes and returns the last elements of the list stored at key.
//
// Starting with Redis version 6.2.0: Added the count argument.
//
// By default, the command pops a single element from the end of the list.
// When provided with the optional count argument, the reply will consist of up to count elements,
// depending on the list's length.
//
// It returns:
//   - When called without the count argument:
//     the value of the last element, or nil when key does not exist.
//   - When called with the count argument:
//     list of popped elements, or nil when key does not exist.
//
// https://redis.io/commands/rpop
func (r GroupList) RPop(ctx context.Context, key string, count ...int) (*gvar.Var, error) {
	if len(count) > 0 {
		return r.Operation.Do(ctx, "RPop", key, count[0])
	}
	return r.Operation.Do(ctx, "RPop", key)
}

// LRem removes the first count occurrences of elements equal to value from the list stored at key.
//
// It returns the number of removed elements.
//
// https://redis.io/commands/lrem/
func (r GroupList) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	v, err := r.Operation.Do(ctx, "LRem", key, count, value)
	return v.Int64(), err
}

// LLen returns the length of the list stored at key.
// Returns the length of the list stored at key.
// If key does not exist, it is interpreted as an empty list and 0 is returned.
// An error is returned when the value stored at key is not a list.
//
// https://redis.io/commands/llen
func (r GroupList) LLen(ctx context.Context, key string) (int64, error) {
	v, err := r.Operation.Do(ctx, "LLen", key)
	return v.Int64(), err
}

// LIndex return the element at index in the list stored at key.
// Returns the element at index in the list stored at key.
// The index is zero-based, so 0 means the first element, 1 the second element and so on.
// Negative indices can be used to designate elements starting at the tail of the list.
// Here, -1 means the last element, -2 means the penultimate and so forth.
// When the value at key is not a list, an error is returned.
//
// It returns the requested element, or nil when index is out of range.
//
// https://redis.io/commands/lindex
func (r GroupList) LIndex(ctx context.Context, key string, index int64) (*gvar.Var, error) {
	return r.Operation.Do(ctx, "LIndex", key, index)
}

// LInsert inserts element in the list stored at key either before or after the reference value pivot.
// When key does not exist, it is considered an empty list and no operation is performed.
// An error is returned when key exists but does not hold a list value.
//
// It returns the length of the list after the insert operation, or -1 when the value pivot was not found.
//
// https://redis.io/commands/linsert/
func (r GroupList) LInsert(ctx context.Context, key string, op gredis.LInsertOp, pivot, value interface{}) (int64, error) {
	v, err := r.Operation.Do(ctx, "LInsert", key, string(op), pivot, value)
	return v.Int64(), err
}

// LSet sets the list element at index to element.
// For more information on the index argument, see LIndex.
// An error is returned for out of range indexes.
//
// https://redis.io/commands/lset/
func (r GroupList) LSet(ctx context.Context, key string, index int64, value interface{}) (*gvar.Var, error) {
	return r.Operation.Do(ctx, "LSet", key, index, value)
}

// LRange returns the specified elements of the list stored at key.
// The offsets start and stop are zero-based indexes, with 0 being the first element of the list (the
// head of the list), 1 being the next element and so on.
//
// These offsets can also be negative numbers indicating offsets starting at the end of the list.
// For example, -1 is the last element of the list, -2 the penultimate, and so on.
//
// https://redis.io/commands/lrange/
func (r GroupList) LRange(ctx context.Context, key string, start, stop int64) (gvar.Vars, error) {
	v, err := r.Operation.Do(ctx, "LRange", key, start, stop)
	return v.Vars(), err
}

// LTrim trims an existing list so that it will contain only the specified range of elements
// specified. Both start and stop are zero-based indexes, where 0 is the first element of the list
// (the head), 1 the next element and so on.
//
// https://redis.io/commands/ltrim/
func (r GroupList) LTrim(ctx context.Context, key string, start, stop int64) error {
	_, err := r.Operation.Do(ctx, "LTrim", key, start, stop)
	return err
}

// BLPop is a blocking list pop primitive.
// It is the blocking version of LPop because it blocks the connection when there are no elements to
// pop from any of the given lists.
// An element is popped from the head of the first list that is non-empty, with the given keys being
// checked in the order that they are given.
//
// The timeout argument is interpreted as a double value specifying the maximum number of seconds to
// block. A timeout of zero can be used to block indefinitely.
//
// https://redis.io/commands/blpop/
func (r GroupList) BLPop(ctx context.Context, timeout int64, keys ...string) (gvar.Vars, error) {
	v, err := r.Operation.Do(ctx, "BLPop", append(gconv.Interfaces(keys), timeout)...)
	return v.Vars(), err
}

// BRPop is a blocking list pop primitive.
// It is the blocking version of RPop because it blocks the connection when there are no elements to
// pop from any of the given lists. An element is popped from the tail of the first list that is
// non-empty, with the given keys being checked in the order that they are given.
//
// The timeout argument is interpreted as a double value specifying the maximum number of seconds to
// block. A timeout of zero can be used to block indefinitely.
//
// https://redis.io/commands/brpop/
func (r GroupList) BRPop(ctx context.Context, timeout int64, keys ...string) (gvar.Vars, error) {
	v, err := r.Operation.Do(ctx, "BRPop", append(gconv.Interfaces(keys), timeout)...)
	return v.Vars(), err
}

// RPopLPush remove the last element in list source, appends it to the front of list destination and
// returns it.
//
// https://redis.io/commands/rpoplpush/
func (r GroupList) RPopLPush(ctx context.Context, source, destination string) (*gvar.Var, error) {
	return r.Operation.Do(ctx, "RPopLPush", source, destination)
}

// BRPopLPush is the blocking variant of RPopLPush.
// When source contains elements, this command behaves exactly like RPopLPush. When used inside a
// MULTI/EXEC block,
// this command behaves exactly like RPopLPush. When source is empty, Redis will block the connection
// until another // client pushes to it or until timeout is reached.
// A timeout of zero can be used to block indefinitely.
//
// It returns the element being popped from source and pushed to destination.
// If timeout is reached, a Null reply is returned.
//
// https://redis.io/commands/brpoplpush/
func (r GroupList) BRPopLPush(ctx context.Context, source, destination string, timeout int64) (*gvar.Var, error) {
	return r.Operation.Do(ctx, "BRPopLPush", source, destination, timeout)
}
