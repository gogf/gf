// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
	"github.com/gogf/gf/v2/container/gvar"
	"testing"
	"time"

	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/test/gtest"
)

func Test_ConfigFromMap(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		c, err := gredis.ConfigFromMap(g.Map{
			`address`:     `127.0.0.1:6379`,
			`db`:          `10`,
			`pass`:        `&*^%$#65Gv`,
			`minIdle`:     `10`,
			`MaxIdle`:     `100`,
			`ReadTimeout`: `10s`,
		})
		t.AssertNil(err)
		t.Assert(c.Address, `127.0.0.1:6379`)
		t.Assert(c.Db, `10`)
		t.Assert(c.Pass, `&*^%$#65Gv`)
		t.Assert(c.MinIdle, 10)
		t.Assert(c.MaxIdle, 100)
		t.Assert(c.ReadTimeout, 10*time.Second)
	})
}

func Test_ConfigAddUser(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		var (
			c   *gredis.Redis
			err error
			r   *gvar.Var
		)

		c, err = gredis.New(&gredis.Config{
			Address: `127.0.0.1`,
			Db:      1,
			User:    "root",
			Pass:    "",
		})
		t.AssertNil(err)

		_, err = c.Conn(ctx)
		t.AssertNil(err)

		_, err = redis.Do(ctx, "SET", "k", "v")
		t.AssertNil(err)

		r, err = redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, []byte("v"))

		_, err = redis.Do(ctx, "DEL", "k")
		t.AssertNil(err)

		r, err = redis.Do(ctx, "GET", "k")
		t.AssertNil(err)
		t.Assert(r, nil)
	})
}
