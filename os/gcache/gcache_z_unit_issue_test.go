package gcache_test

import (
	"context"
	"testing"
	"time"

	"github.com/gogf/gf/v2/os/gcache"
	"github.com/gogf/gf/v2/test/gtest"
)

// https://github.com/gogf/gf/issues/4145
func Test_Issue4145(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			ctx        = context.Background()
			cache      = gcache.New()
			cacheKey1  = "GetTest-1"
			cacheKey2  = "GetTest2-1"
			cacheValue = "123456789"
		)

		// 定义需要测试的闭包函数
		getTestCached := func(ctx context.Context) (*string, error) {
			v, err := cache.GetOrSetFuncLock(ctx, cacheKey1, func(ctx context.Context) (interface{}, error) {
				str := cacheValue
				return &str, nil
			}, 1*time.Minute)

			if err != nil {
				return nil, err
			}

			var res *string
			if err := v.Struct(&res); err != nil {
				return nil, err
			}
			return res, nil
		}

		getTest2Cached := func(ctx context.Context) (*string, error) {
			v, err := cache.GetOrSetFuncLock(ctx, cacheKey2, func(ctx context.Context) (interface{}, error) {
				// 内部调用 getTestCached
				return getTestCached(ctx)
			}, 1*time.Minute)

			if err != nil {
				return nil, err
			}

			var res *string
			if err := v.Struct(&res); err != nil {
				return nil, err
			}
			return res, nil
		}

		// 测试用例
		// 第一次获取应该走实际逻辑
		value, err := getTestCached(ctx)
		t.AssertNil(err)
		t.Assert(*value, cacheValue)

		// 第二次获取应该走缓存
		v, err := cache.Get(ctx, cacheKey1)
		t.AssertNil(err)
		t.Assert(v, cacheValue)

		// 测试嵌套缓存调用
		value, err = getTest2Cached(ctx)
		t.AssertNil(err)
		t.Assert(*value, cacheValue)

		// 验证二级缓存
		v, err = cache.Get(ctx, cacheKey2)
		t.AssertNil(err)
		t.Assert(v, cacheValue)

		// 清理所有缓存
		_, err = cache.Remove(ctx, cacheKey1, cacheKey2)
		t.AssertNil(err)

		// 验证清理结果
		v1, _ := cache.Get(ctx, cacheKey1)
		v2, _ := cache.Get(ctx, cacheKey2)
		t.Assert(v1, nil)
		t.Assert(v2, nil)
	})
}
