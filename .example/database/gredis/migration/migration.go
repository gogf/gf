package main

import (
	"fmt"

	"github.com/gogf/gf/util/gconv"

	"github.com/gogf/gf/os/gcmd"

	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/database/gredis"
	"github.com/gogf/gf/frame/g"
)

func main() {
	batchNumber := 1000
	redis1Config, err := gredis.ConfigFromStr("im-redis-slave:6379,9")
	if err != nil {
		panic(err)
	}
	redis2Config, err := gredis.ConfigFromStr("r-bp1f0a5d4efd8744.redis.rds.aliyuncs.com:6379,9")
	if err != nil {
		panic(err)
	}
	gredis.SetConfig(redis1Config)

	v, err := g.Redis().DoVar("keys", "*")
	if err != nil {
		panic(err)
	}
	array := garray.NewStrArrayFrom(v.Strings())
	for {
		slice := array.PopLefts(batchNumber)
		if len(slice) > 0 {
			// `migrate %s %d "" 0 2000 copy replace auth %s keys %s`,
			params := g.Slice{
				redis2Config.Host,
				redis2Config.Port,
				"",
				redis2Config.Db,
				2000,
				"copy",
				"replace",
				"keys",
			}
			params = append(params, gconv.Interfaces(slice)...)
			fmt.Println(params)
			if gcmd.GetOpt("dryrun") == "0" {
				if v, err := g.Redis().DoVar("migrate", params...); err != nil {
					panic(err)
				} else {
					fmt.Println(v.String())
				}
			}
		} else {
			break
		}
	}
	fmt.Println("done")
}
