package gtime

import (
    "time"
)

// 类似与js中的SetTimeout，一段时间后执行回调函数
func SetTimeout(t time.Duration, callback func()) {
    go func() {
        time.Sleep(t)
        callback()
    }()
}

// 类似与js中的SetInterval，每隔一段时间后执行回调函数，当回调函数返回true，那么继续执行，否则终止执行，该方法是异步的
// 注意：由于采用的是循环而不是递归操作，因此间隔时间将会以上一次回调函数执行完成的时间来计算
func SetInterval(t time.Duration, callback func() bool) {
    go func() {
        for {
            time.Sleep(t)
            r := callback()
            if !r {
                break;
            }
        }
    }()
}

// 获取当前的毫秒数
func Millisecond() int64 {
    return time.Now().UnixNano()/1e6
}
