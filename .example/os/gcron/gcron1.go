package main

import (
	"time"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/glog"
)

func main() {
	gcron.Add("0 30 * * * *", func() { glog.Println("Every hour on the half hour") })
	gcron.Add("* * * * * *", func() { glog.Println("Every second, pattern") }, "second-cron")
	gcron.Add("*/5 * * * * *", func() { glog.Println("Every 5 seconds, pattern") })

	gcron.Add("@hourly", func() { glog.Println("Every hour") })
	gcron.Add("@every 1h30m", func() { glog.Println("Every hour thirty") })
	gcron.Add("@every 1s", func() { glog.Println("Every 1 second") })
	gcron.Add("@every 5s", func() { glog.Println("Every 5 seconds") })

	time.Sleep(3 * time.Second)

	gcron.Stop("second-cron")

	time.Sleep(3 * time.Second)

	gcron.Start("second-cron")

	time.Sleep(10 * time.Second)
}
