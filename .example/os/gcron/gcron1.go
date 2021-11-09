package main

import (
	"time"

	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/glog"
)

func main() {
	gcron.Add("0 30 * * * *", func() { glog.Print("Every hour on the half hour") })
	gcron.Add("* * * * * *", func() { glog.Print("Every second, pattern") }, "second-cron")
	gcron.Add("*/5 * * * * *", func() { glog.Print("Every 5 seconds, pattern") })

	gcron.Add("@hourly", func() { glog.Print("Every hour") })
	gcron.Add("@every 1h30m", func() { glog.Print("Every hour thirty") })
	gcron.Add("@every 1s", func() { glog.Print("Every 1 second") })
	gcron.Add("@every 5s", func() { glog.Print("Every 5 seconds") })

	time.Sleep(3 * time.Second)

	gcron.Stop("second-cron")

	time.Sleep(3 * time.Second)

	gcron.Start("second-cron")

	time.Sleep(10 * time.Second)
}
