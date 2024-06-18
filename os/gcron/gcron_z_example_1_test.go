// Copyright GoFrame Author(https://goframe.org). All Rights Reserved.
//
// This Source Code Form is subject to the terms of the MIT License.
// If a copy of the MIT was not distributed with this file,
// You can obtain one at https://github.com/gogf/gf.

package gcron_test

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gcron"
	"github.com/gogf/gf/v2/os/glog"
)

func ExampleCronAddSingleton() {
	gcron.AddSingleton(ctx, "* * * * * *", func(ctx context.Context) {
		glog.Print(context.TODO(), "doing")
		time.Sleep(2 * time.Second)
	})
	select {}
}

func ExampleCronGracefulShutdown() {
	_, err := gcron.Add(ctx, "*/2 * * * * *", func(ctx context.Context) {
		g.Log().Debug(ctx, "Every 2s job start")
		time.Sleep(5 * time.Second)
		g.Log().Debug(ctx, "Every 2s job after 5 second end")
	}, "MyCronJob")
	if err != nil {
		panic(err)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	sig := <-quit
	glog.Printf(ctx, "Signal received: %s, stopping cron", sig)

	glog.Print(ctx, "Waiting for all cron jobs to complete...")
	gcron.StopGracefully()
	glog.Print(ctx, "All cron jobs completed")
}
