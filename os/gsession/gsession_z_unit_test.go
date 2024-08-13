// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession

import (
	"fmt"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gctx"
	"testing"
	"time"

	"github.com/gogf/gf/v2/test/gtest"

	_ "github.com/gogf/gf/contrib/nosql/redis/v2"
)

func Test_NewSessionId(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {
		id1 := NewSessionId()
		id2 := NewSessionId()
		t.AssertNE(id1, id2)
		t.Assert(len(id1), 32)
	})
}

func TestSession_RegenSession(t *testing.T) {
	gtest.C(t, func(t *gtest.T) {

		config := &gredis.Config{
			Address:     "192.168.5.1:6379",
			Db:          10,
			Pass:        "123456",
			IdleTimeout: 600,
			MinIdle:     2,
			WaitTimeout: 30,
		}
		rds, erRds := gredis.New(config)
		if erRds != nil {
			fmt.Println("err:", erRds)
			return
		}
		storage := NewStorageRedisHashTable(rds, "tsid:")
		sManger := New(time.Second*600, storage)
		sid := "" //"19t8ieo1300d3eh4xqsnlbw100xj9q7q"
		session := sManger.New(gctx.New(), sid)
		sid, _ = session.Id()
		fmt.Println(sid)
		session.Set("score", 6666666)
		session.Set("name", "erretre")
		session.Close()
		score := session.MustGet("score")
		fmt.Println("score:", score.String())
		name := session.MustGet("name")
		fmt.Println("name:", name.String())
		data, _ := session.Data()
		fmt.Println(fmt.Sprintf("%+v", data))

		newSid, errNewSid := session.RegenSession(true)
		if errNewSid != nil {
			fmt.Println("regen session err:", errNewSid)
			return
		}
		fmt.Println("new sid:", newSid)
		session.Set("score", "99999999999")
		session.Set("tttt", 999999)
		session.Close()

		newScore, _ := session.Get("score")
		fmt.Println("newScore:", newScore, " newSid:", session.id)

		tttt, _ := session.Get("tttt")
		fmt.Println("tttt:", tttt, " newSid:", session.id)

		oldSid := sid
		sessionOld := sManger.New(gctx.New(), oldSid)
		oldData, _ := sessionOld.Data()
		fmt.Println("oldData:", oldData, " oldSid:", oldSid)
		//time.Sleep(time.Second * 2)
	})
}
