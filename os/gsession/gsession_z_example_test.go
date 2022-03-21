// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gsession_test

import (
	"fmt"
	"time"

	"github.com/gogf/gf/v2/container/gmap"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/os/gctx"
	"github.com/gogf/gf/v2/os/gsession"
)

func ExampleNew() {
	manager := gsession.New(time.Second)
	fmt.Println(manager.TTL())

	// Output:
	// 1s
}

func ExampleManager_SetStorage() {
	manager := gsession.New(time.Second)
	manager.SetStorage(gsession.NewStorageMemory())
	fmt.Println(manager.TTL())

	// Output:
	// 1s
}

func ExampleManager_GetStorage() {
	manager := gsession.New(time.Second, gsession.NewStorageMemory())
	size, _ := manager.GetStorage().GetSize(gctx.New(), "id")
	fmt.Println(size)

	// Output:
	// -1
}

func ExampleManager_SetTTL() {
	manager := gsession.New(time.Second)
	manager.SetTTL(time.Minute)
	fmt.Println(manager.TTL())

	// Output:
	// 1m0s
}

func ExampleSession_Set() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())
	s := manager.New(gctx.New())
	fmt.Println(s.Set("key", "val") == nil)

	// Output:
	// true
}

func ExampleSession_SetMap() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())
	s := manager.New(gctx.New())
	fmt.Println(s.SetMap(map[string]interface{}{}) == nil)

	// Output:
	// true
}

func ExampleSession_Remove() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())
	s1 := manager.New(gctx.New())
	fmt.Println(s1.Remove("key"))

	s2 := manager.New(gctx.New(), "Remove")
	fmt.Println(s2.Remove("key"))

	// Output:
	// <nil>
	// <nil>
}

func ExampleSession_RemoveAll() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())
	s1 := manager.New(gctx.New())
	fmt.Println(s1.RemoveAll())

	s2 := manager.New(gctx.New(), "Remove")
	fmt.Println(s2.RemoveAll())

	// Output:
	// <nil>
	// <nil>
}

func ExampleSession_Id() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())
	s := manager.New(gctx.New(), "Id")
	id, _ := s.Id()
	fmt.Println(id)

	// Output:
	// Id
}

func ExampleSession_SetId() {
	nilSession := &gsession.Session{}
	fmt.Println(nilSession.SetId("id"))

	manager := gsession.New(time.Second, gsession.NewStorageFile())
	s := manager.New(gctx.New())
	s.Id()
	fmt.Println(s.SetId("id"))

	// Output:
	// <nil>
	// session already started
}

func ExampleSession_SetIdFunc() {
	nilSession := &gsession.Session{}
	fmt.Println(nilSession.SetIdFunc(func(ttl time.Duration) string {
		return "id"
	}))

	manager := gsession.New(time.Second, gsession.NewStorageFile())
	s := manager.New(gctx.New())
	s.Id()
	fmt.Println(s.SetIdFunc(func(ttl time.Duration) string {
		return "id"
	}))

	// Output:
	// <nil>
	// session already started
}

func ExampleSession_Data() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())

	s1 := manager.New(gctx.New())
	data1, _ := s1.Data()
	fmt.Println(data1)

	s2 := manager.New(gctx.New(), "id_data")
	data2, _ := s2.Data()
	fmt.Println(data2)

	// Output:
	// map[]
	// map[]
}

func ExampleSession_Size() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())

	s1 := manager.New(gctx.New())
	size1, _ := s1.Size()
	fmt.Println(size1)

	s2 := manager.New(gctx.New(), "Size")
	size2, _ := s2.Size()
	fmt.Println(size2)

	// Output:
	// 0
	// 0
}

func ExampleSession_Contains() {
	manager := gsession.New(time.Second, gsession.NewStorageFile())

	s1 := manager.New(gctx.New())
	notContains, _ := s1.Contains("Contains")
	fmt.Println(notContains)

	s2 := manager.New(gctx.New(), "Contains")
	contains, _ := s2.Contains("Contains")
	fmt.Println(contains)

	// Output:
	// false
	// false
}

func ExampleStorageFile_SetCryptoKey() {
	storage := gsession.NewStorageFile()
	storage.SetCryptoKey([]byte("key"))

	size, _ := storage.GetSize(gctx.New(), "id")
	fmt.Println(size)

	// Output:
	// -1
}

func ExampleStorageFile_SetCryptoEnabled() {
	storage := gsession.NewStorageFile()
	storage.SetCryptoEnabled(true)

	size, _ := storage.GetSize(gctx.New(), "id")
	fmt.Println(size)

	// Output:
	// -1
}

func ExampleStorageFile_UpdateTTL() {
	var (
		ctx = gctx.New()
	)

	storage := gsession.NewStorageFile()
	fmt.Println(storage.UpdateTTL(ctx, "id", time.Second*15))

	time.Sleep(time.Second * 11)

	// Output:
	// <nil>
}

func ExampleStorageRedis_Get() {
	storage := gsession.NewStorageRedis(&gredis.Redis{})
	val, _ := storage.Get(gctx.New(), "id", "key")
	fmt.Println(val)

	// Output:
	// <nil>
}

func ExampleStorageRedis_Data() {
	storage := gsession.NewStorageRedis(&gredis.Redis{})
	val, _ := storage.Data(gctx.New(), "id")
	fmt.Println(val)

	// Output:
	// map[]
}

func ExampleStorageRedis_GetSize() {
	storage := gsession.NewStorageRedis(&gredis.Redis{})
	val, _ := storage.GetSize(gctx.New(), "id")
	fmt.Println(val)

	// Output:
	// -1
}

func ExampleStorageRedis_Remove() {
	storage := gsession.NewStorageRedis(&gredis.Redis{})
	err := storage.Remove(gctx.New(), "id", "key")
	fmt.Println(err != nil)

	// Output:
	// true
}

func ExampleStorageRedis_RemoveAll() {
	storage := gsession.NewStorageRedis(&gredis.Redis{})
	err := storage.RemoveAll(gctx.New(), "id")
	fmt.Println(err != nil)

	// Output:
	// true
}

func ExampleStorageRedis_UpdateTTL() {
	storage := gsession.NewStorageRedis(&gredis.Redis{})
	err := storage.UpdateTTL(gctx.New(), "id", time.Second*15)
	fmt.Println(err)

	time.Sleep(time.Second * 11)

	// Output:
	// <nil>
}

func ExampleStorageRedisHashTable_Get() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	v, err := storage.Get(gctx.New(), "id", "key")

	fmt.Println(v)
	fmt.Println(err)

	// Output:
	// <nil>
	// redis adapter not initialized, missing configuration or adapter register?
}

func ExampleStorageRedisHashTable_Data() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	data, err := storage.Data(gctx.New(), "id")

	fmt.Println(data)
	fmt.Println(err)

	// Output:
	// map[]
	// redis adapter not initialized, missing configuration or adapter register?
}

func ExampleStorageRedisHashTable_GetSize() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	size, err := storage.GetSize(gctx.New(), "id")

	fmt.Println(size)
	fmt.Println(err)

	// Output:
	// -1
	// redis adapter not initialized, missing configuration or adapter register?
}

func ExampleStorageRedisHashTable_Remove() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	err := storage.Remove(gctx.New(), "id", "key")

	fmt.Println(err)

	// Output:
	// redis adapter not initialized, missing configuration or adapter register?
}

func ExampleStorageRedisHashTable_RemoveAll() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	err := storage.RemoveAll(gctx.New(), "id")

	fmt.Println(err)

	// Output:
	// redis adapter not initialized, missing configuration or adapter register?
}

func ExampleStorageRedisHashTable_GetSession() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	strAnyMap := gmap.StrAnyMap{}

	data, err := storage.GetSession(gctx.New(), "id", time.Second, &strAnyMap)

	fmt.Println(data)
	fmt.Println(err)

	// Output:
	//
	// redis adapter not initialized, missing configuration or adapter register?
}

func ExampleStorageRedisHashTable_SetSession() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	strAnyMap := gmap.StrAnyMap{}

	err := storage.SetSession(gctx.New(), "id", &strAnyMap, time.Second)

	fmt.Println(err)

	// Output:
	// redis adapter not initialized, missing configuration or adapter register?
}

func ExampleStorageRedisHashTable_UpdateTTL() {
	storage := gsession.NewStorageRedisHashTable(&gredis.Redis{})

	err := storage.UpdateTTL(gctx.New(), "id", time.Second)

	fmt.Println(err)

	// Output:
	// redis adapter not initialized, missing configuration or adapter register?
}
