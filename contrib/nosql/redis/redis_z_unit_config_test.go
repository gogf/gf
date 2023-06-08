// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package redis_test

import (
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
			c    *gredis.Redis
			conn gredis.Conn
			res  *gvar.Var
			err  error
		)

		c, err = gredis.New(&gredis.Config{
			Address: `redis-cloud.aws-west-xx.ec2.redislabs.com:18311`,
			Db:      1,
			User:    "userNameOrEmpty",
			Pass:    "passwordOrEmpty",
		})
		t.AssertNil(err)

		conn, err = c.Conn(ctx)
		t.AssertNil(err)

		res, err = conn.Do(ctx, "ping")
		t.AssertNil(err)
		t.Assert(res.String(), "PONG")
	})
}
