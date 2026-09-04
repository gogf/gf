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
	"github.com/gogf/gf/v2/util/gconv"
)

// Test_Pipeline_HSetGet verifies that HSet and HGet commands queued via pipeline
// are executed correctly and results are populated after Exec.
func Test_Pipeline_HSetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		defer r.Do(ctx, "DEL", "test:pipe:hset")

		pipe := r.Pipeline(ctx)
		t.AssertNE(pipe, nil)

		cmd1 := pipe.PipelineGroupHash().HSet(ctx, "test:pipe:hset", map[string]any{"field1": "value1"})
		t.AssertNE(cmd1, nil)

		cmd2 := pipe.PipelineGroupHash().HGet(ctx, "test:pipe:hset", "field1")
		t.AssertNE(cmd2, nil)

		results, err := pipe.Exec(ctx)
		t.AssertNil(err)
		t.Assert(len(results), 2)

		val2, err := cmd2.Result()
		t.AssertNil(err)
		t.Assert(val2.String(), "value1")
	})
}

// Test_Pipeline_StringSetGet queues Set, Get and Incr commands via pipeline,
// executes them, and verifies all results are populated correctly.
func Test_Pipeline_StringSetGet(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		defer r.Do(ctx, "DEL", "test:pipe:str")

		pipe := r.Pipeline(ctx)
		t.AssertNE(pipe, nil)

		pipe.PipelineGroupString().Set(ctx, "test:pipe:str", "10")
		cmdGet := pipe.PipelineGroupString().Get(ctx, "test:pipe:str")
		cmdIncr := pipe.PipelineGroupString().Incr(ctx, "test:pipe:str")

		results, err := pipe.Exec(ctx)
		t.AssertNil(err)
		t.Assert(len(results), 3)

		valGet, err := cmdGet.Result()
		t.AssertNil(err)
		t.Assert(valGet.String(), "10")

		valIncr, err := cmdIncr.Result()
		t.AssertNil(err)
		t.Assert(valIncr.Int64(), int64(11))
	})
}

// Test_Pipeline_MultipleGroups queues commands from different command groups
// (Hash.HSet, String.Set, Generic.Expire) in the same pipeline and verifies execution.
func Test_Pipeline_MultipleGroups(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		defer r.Do(ctx, "DEL", "test:pipe:multi:hash", "test:pipe:multi:str")

		pipe := r.Pipeline(ctx)
		t.AssertNE(pipe, nil)

		cmdHSet := pipe.PipelineGroupHash().HSet(ctx, "test:pipe:multi:hash", map[string]any{"k": "v"})
		cmdSet := pipe.PipelineGroupString().Set(ctx, "test:pipe:multi:str", "hello")
		cmdExpire := pipe.PipelineGroupGeneric().Expire(ctx, "test:pipe:multi:str", 3600)

		results, err := pipe.Exec(ctx)
		t.AssertNil(err)
		t.Assert(len(results), 3)

		valHSet, err := cmdHSet.Result()
		t.AssertNil(err)
		t.Assert(valHSet.Int64(), int64(1))

		valSet, err := cmdSet.Result()
		t.AssertNil(err)
		t.Assert(valSet.String(), "OK")

		valExpire, err := cmdExpire.Result()
		t.AssertNil(err)
		t.Assert(valExpire.Int64(), int64(1))

		v, err := r.HGet(ctx, "test:pipe:multi:hash", "k")
		t.AssertNil(err)
		t.Assert(v.String(), "v")

		v2, err := r.Do(ctx, "GET", "test:pipe:multi:str")
		t.AssertNil(err)
		t.Assert(v2.String(), "hello")
	})
}

// Test_Pipeline_Discard queues commands, calls Discard to cancel them,
// then verifies the keys were never created on the server.
func Test_Pipeline_Discard(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		pipe := r.Pipeline(ctx)
		t.AssertNE(pipe, nil)

		pipe.PipelineGroupString().Set(ctx, "test:pipe:discard", "should_not_exist")

		pipe.Discard()

		v, err := r.Do(ctx, "GET", "test:pipe:discard")
		t.AssertNil(err)
		t.Assert(v.IsNil(), true)

		defer r.Do(ctx, "DEL", "test:pipe:discard")
	})
}

// Test_Pipeline_DoRaw uses the low-level pipe.Do method to queue a SET command,
// then verifies the result after Exec.
func Test_Pipeline_DoRaw(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		defer r.Do(ctx, "DEL", "test:pipe:doraw")

		pipe := r.Pipeline(ctx)
		t.AssertNE(pipe, nil)

		cmd := pipe.Do(ctx, "SET", "test:pipe:doraw", "rawvalue")
		t.AssertNE(cmd, nil)

		results, err := pipe.Exec(ctx)
		t.AssertNil(err)
		t.Assert(len(results), 1)

		val, err := cmd.Result()
		t.AssertNil(err)
		t.Assert(val.String(), "OK")

		v, err := r.Do(ctx, "GET", "test:pipe:doraw")
		t.AssertNil(err)
		t.Assert(v.String(), "rawvalue")
	})
}

// Test_TxPipeline_Basic queues commands via TxPipeline (MULTI/EXEC),
// executes them, and verifies atomic execution.
func Test_TxPipeline_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		defer r.Do(ctx, "DEL", "test:tx:basic")

		tx := r.TxPipeline(ctx)
		t.AssertNE(tx, nil)

		cmdSet := tx.PipelineGroupString().Set(ctx, "test:tx:basic", "txval")
		cmdGet := tx.PipelineGroupString().Get(ctx, "test:tx:basic")

		results, err := tx.Exec(ctx)
		t.AssertNil(err)
		t.Assert(len(results), 2)

		valSet, err := cmdSet.Result()
		t.AssertNil(err)
		t.Assert(valSet.String(), "OK")

		valGet, err := cmdGet.Result()
		t.AssertNil(err)
		t.Assert(valGet.String(), "txval")
	})
}

// Test_TxPipeline_MultiCommands queues multiple HIncrBy commands in a transaction,
// verifies they are all applied atomically, and checks the final counter value.
func Test_TxPipeline_MultiCommands(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		defer r.Do(ctx, "DEL", "test:tx:counter")

		tx := r.TxPipeline(ctx)
		t.AssertNE(tx, nil)

		cmd1 := tx.PipelineGroupHash().HIncrBy(ctx, "test:tx:counter", "count", 10)
		cmd2 := tx.PipelineGroupHash().HIncrBy(ctx, "test:tx:counter", "count", 20)
		cmd3 := tx.PipelineGroupHash().HIncrBy(ctx, "test:tx:counter", "count", 30)

		results, err := tx.Exec(ctx)
		t.AssertNil(err)
		t.Assert(len(results), 3)

		val1, err := cmd1.Result()
		t.AssertNil(err)
		t.Assert(val1.Int64(), int64(10))

		val2, err := cmd2.Result()
		t.AssertNil(err)
		t.Assert(val2.Int64(), int64(30))

		val3, err := cmd3.Result()
		t.AssertNil(err)
		t.Assert(val3.Int64(), int64(60))

		v, err := r.HGet(ctx, "test:tx:counter", "count")
		t.AssertNil(err)
		t.Assert(v.Int64(), int64(60))
	})
}

// Test_ScanAll_Basic sets several keys with a common prefix, calls ScanAll
// with a matching pattern, and verifies all keys are returned.
func Test_ScanAll_Basic(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		prefix := "test:scanall:"
		for i := range 5 {
			_, err := r.Do(ctx, "SET", prefix+gconv.String(i), "val")
			t.AssertNil(err)
		}
		defer func() {
			for i := range 5 {
				r.Do(ctx, "DEL", prefix+gconv.String(i))
			}
		}()

		keys, err := r.GroupGeneric().ScanAll(ctx, gredis.ScanOption{
			Match: "test:scanall:*",
			Count: 10,
		})
		t.AssertNil(err)
		t.AssertGE(len(keys), 5)
	})
}

// Test_ScanAll_EmptyPattern verifies that ScanAll with a pattern matching
// no keys returns an empty slice without error.
func Test_ScanAll_EmptyPattern(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		keys, err := r.GroupGeneric().ScanAll(ctx, gredis.ScanOption{
			Match: "test:scanall:nonexistent:*",
			Count: 10,
		})
		t.AssertNil(err)
		t.Assert(len(keys), 0)
	})
}

// Test_Del_MultiKey verifies that Del removes multiple keys in a single call
// and returns the correct count of deleted keys.
func Test_Del_MultiKey(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		_, err = r.Do(ctx, "SET", "test:del:1", "a")
		t.AssertNil(err)
		_, err = r.Do(ctx, "SET", "test:del:2", "b")
		t.AssertNil(err)
		_, err = r.Do(ctx, "SET", "test:del:3", "c")
		t.AssertNil(err)

		n, err := r.GroupGeneric().Del(ctx, "test:del:1", "test:del:2", "test:del:3")
		t.AssertNil(err)
		t.Assert(n, int64(3))

		v, err := r.Do(ctx, "GET", "test:del:1")
		t.AssertNil(err)
		t.Assert(v.IsNil(), true)
	})
}

// Test_Watch_Success sets a key, starts a Watch transaction on it, queues
// a modification within the callback, and verifies the transaction commits
// successfully when no concurrent modification occurs.
func Test_Watch_Success(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		r, err := gredis.New(config)
		t.AssertNil(err)
		defer r.Close(ctx)

		defer r.Do(ctx, "DEL", "test:watch:key")

		_, err = r.Do(ctx, "SET", "test:watch:key", "initial")
		t.AssertNil(err)

		err = r.Watch(ctx, func(tx gredis.Tx) error {
			cmdSet := tx.PipelineGroupString().Set(ctx, "test:watch:key", "watched")
			cmdGet := tx.PipelineGroupString().Get(ctx, "test:watch:key")
			results, execErr := tx.Exec(ctx)
			if execErr != nil {
				return execErr
			}
			t.Assert(len(results), 2)

			valSet, _ := cmdSet.Result()
			t.Assert(valSet.String(), "OK")

			valGet, _ := cmdGet.Result()
			t.Assert(valGet.String(), "watched")
			return nil
		}, "test:watch:key")
		t.AssertNil(err)

		v, err := r.Do(ctx, "GET", "test:watch:key")
		t.AssertNil(err)
		t.Assert(v.String(), "watched")
	})
}
