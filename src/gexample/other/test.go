package main

import (
    "g/os/glog"
)


func main() {
    glog.SetLogPath("/root")
    glog.Info("test")
    //glog.Error("test")
}