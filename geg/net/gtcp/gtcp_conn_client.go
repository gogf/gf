package main

import (
    "gitee.com/johng/gf/g/net/gtcp"
    "gitee.com/johng/gf/g/util/gconv"
    "gitee.com/johng/gf/g/os/glog"
    "time"
)

func main() {
    conn, err := gtcp.NewConn("127.0.0.1:8999")
    if err != nil {
        glog.Fatal(err)
    }
    for i := 0; i < 10000; i++ {
        if err := conn.Send([]byte(gconv.String(i))); err != nil {
            glog.Error(err)
        }
        time.Sleep(time.Second)
    }
}