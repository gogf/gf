package main

import (
    "fmt"
    "g/encoding/gmd5"
    "strings"
    "os"
    "g/os/glog"
    "g/net/gip"
    "sort"
)

func makeNodeId() string {
    hostname, err := os.Hostname()
    if err != nil {
        glog.Fatalln("getting local hostname failed:", err)
    }
    ips, err      := gip.IntranetIP()
    if err != nil {
        glog.Fatalln("getting local ips:", err)
    }
    // 如果有多个IP，那么将IP升序排序
    sort.Slice(ips, func(i, j int) bool { return ips[i] < ips[j] })
    return strings.ToUpper(gmd5.EncodeString(fmt.Sprintf("%s/%v", hostname, strings.Join(ips, ","))))
}

func main() {
    fmt.Println(makeNodeId())
}