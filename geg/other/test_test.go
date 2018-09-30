package main

import (
    "testing"
    "gitee.com/johng/gf/g/container/garray"
)

func TestGFArray001(t *testing.T) {
    var source = []string{"59705a2c1fd50736a4c768a1", "597a95ff1fd5073e48bb2272", "597a960f1fd5073e48bb2274"}
    var CacheChannelKeys = garray.NewSortedStringArray(0)
    CacheChannelKeys.Add(source...)
    t.Logf("%#v\n", CacheChannelKeys)

    CacheChannelKeys.Clear()
    CacheChannelKeys = garray.NewSortedStringArray(len(source))
    t.Logf("%#v\n", CacheChannelKeys)
    CacheChannelKeys.Add(source...)
    t.Logf("%#v\n", CacheChannelKeys)
}