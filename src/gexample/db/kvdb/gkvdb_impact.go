// 分区碰撞测试，测试随机字符串使用哈希函数后的碰撞几率
package main

import (
    "g/util/gtime"
    "g/encoding/ghash"
    "strconv"
    "fmt"
    "g/util/grand"
    "sync"
)

var wg   sync.WaitGroup
var lock sync.RWMutex
var list []uint64 = make([]uint64, 0)

func main() {
    t1 := gtime.Second()
    // 生成1亿随机数据
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(n int) {
            data := make([]uint64, 0)
            for j := 1000000*n; j < 1000000*(n+1); j++ {
                key := ghash.SDBMHash64([]byte("key"+strconv.Itoa(j)+"_with_rand_"+grand.RandStr(10))) % 100000
                data = append(data, key)
            }
            fmt.Println("done", n)
            lock.Lock()
            list = append(list, data...)
            lock.Unlock()
            wg.Done()

        }(i)
    }
    wg.Wait()

    // 判断重复数据
    m  := make(map[uint64]uint32)
    for _, key := range list {
        if _, ok := m[key]; ok {
            m[key]++
        } else {
            m[key] = 1
        }
    }
    fmt.Println(gtime.Second() - t1)

    // 检查最大最富条数
    var max uint32 = 0
    for _, v := range m {
        if v > max || max == 0 {
            max = v
        }
    }
    fmt.Println("conflicts max:", max)
}