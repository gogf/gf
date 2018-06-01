package sleep

import (
    "time"
    "gitee.com/johng/gf/g/os/glog"
)

func init () {
    glog.Println("sleep package importing")
    time.Sleep(3*time.Second)
    glog.Println("sleep package imported")
}

func Test() {
    glog.Println("Test")
}
