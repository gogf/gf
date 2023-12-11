package main

import (
	"context"
	"fmt"

	"github.com/gogf/gf/contrib/nosql/redis/v2"
	_ "github.com/gogf/gf/contrib/nosql/redis/v2"

	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/os/gctx"
)

var (
	config = gredis.Config{
		Address: "127.0.0.1:6379",
		Db:      1,
	}
	group = "cache"
	ctx   = gctx.New()
)

// Myredis description
type Myredis struct {
	*redis.Redis
}

func init() {
	fmt.Println("init")
	gredis.RegisterAdapterFunc(func(config *gredis.Config) gredis.Adapter {
		return &Myredis{
			redis.New(config),
		}
	})
}

// GroupString is the redis group object for string operations.
func (r *Myredis) GroupString() gredis.IGroupString {
	fmt.Println("Myredis GroupString")
	return redis.GroupString{
		Redis: r,
	}
}

// GroupGeneric creates and returns GroupGeneric.
func (r *Myredis) GroupGeneric() gredis.IGroupGeneric {
	return redis.GroupGeneric{
		Redis: r,
	}
}

// GroupHash creates and returns a redis group object for hash operations.
func (r *Myredis) GroupHash() gredis.IGroupHash {
	return redis.GroupHash{
		Redis: r,
	}
}

// GroupList creates and returns a redis group object for list operations.
func (r *Myredis) GroupList() gredis.IGroupList {
	return redis.GroupList{
		Redis: r,
	}
}

// GroupPubSub creates and returns GroupPubSub.
func (r *Myredis) GroupPubSub() gredis.IGroupPubSub {
	return redis.GroupPubSub{
		Redis: r,
	}
}

// GroupScript creates and returns GroupScript.
func (r *Myredis) GroupScript() gredis.IGroupScript {
	return redis.GroupScript{
		Redis: r,
	}
}

// GroupSet creates and returns GroupSet.
func (r *Myredis) GroupSet() gredis.IGroupSet {
	return redis.GroupSet{
		Redis: r,
	}
}

// GroupSortedSet creates and returns GroupSortedSet.
func (r *Myredis) GroupSortedSet() gredis.IGroupSortedSet {
	return redis.GroupSortedSet{
		Redis: r,
	}
}

// Do send a command to the server and returns the received reply.
// It uses json.Marshal for struct/slice/map type values before committing them to redis.
func (r *Myredis) Do(ctx context.Context, command string, args ...interface{}) (*gvar.Var, error) {
	fmt.Println("Myredis Do", command, args)
	conn, err := r.Conn(ctx)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = conn.Close(ctx)
	}()
	return conn.Do(ctx, command, args...)
}

func main() {
	gredis.SetConfig(&config, group)

	_, err := g.Redis(group).Set(ctx, "key", "value")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	value, err := g.Redis(group).Get(ctx, "key")
	if err != nil {
		g.Log().Fatal(ctx, err)
	}
	fmt.Println(value.String())
}
