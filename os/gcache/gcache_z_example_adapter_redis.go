package gcache

import (
	"context"
	"github.com/gogf/gf/v2/container/gvar"
	"github.com/gogf/gf/v2/database/gredis"
	"time"
)

// Redis is the gcache adapter implements using Redis server.
type Redis struct {
	redis *gredis.Redis
}

// newAdapterMemory creates and returns a new memory cache object.
func NewRedis(redis *gredis.Redis) Adapter {
	return &Redis{
		redis: redis,
	}
}

func (r Redis) Set(ctx context.Context, key interface{}, value interface{}, duration time.Duration) error {
	var err error
	if value == nil || duration < 0 {
		_, err = r.redis.Do(ctx,"DEL", key)
	} else {
		if duration == 0 {
			_, err = r.redis.Do(ctx,"SET", key, value)
		} else {
			_, err = r.redis.Do(ctx,"SETEX", key, uint64(duration.Seconds()), value)
		}
	}
	return err
}

func (r Redis) SetMap(ctx context.Context, data map[interface{}]interface{}, duration time.Duration) error {
	panic("implement me")
}

func (r Redis) SetIfNotExist(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (ok bool, err error) {
	panic("implement me")
}

func (r Redis) SetIfNotExistFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (ok bool, err error) {
	panic("implement me")
}

func (r Redis) SetIfNotExistFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (ok bool, err error) {
	panic("implement me")
}

func (r Redis) Get(ctx context.Context, key interface{}) (*gvar.Var, error) {
	v, err := r.redis.Do(ctx,"GET", key)
	if err != nil {
		return nil, err
	}

	return v, nil
}

func (r Redis) GetOrSet(ctx context.Context, key interface{}, value interface{}, duration time.Duration) (result *gvar.Var, err error) {
	panic("implement me")
}

func (r Redis) GetOrSetFunc(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (result *gvar.Var, err error) {
	panic("implement me")
}

func (r Redis) GetOrSetFuncLock(ctx context.Context, key interface{}, f func() (interface{}, error), duration time.Duration) (result *gvar.Var, err error) {
	panic("implement me")
}

func (r Redis) Contains(ctx context.Context, key interface{}) (bool, error) {
	panic("implement me")
}

func (r Redis) Size(ctx context.Context) (size int, err error) {
	panic("implement me")
}

func (r Redis) Data(ctx context.Context) (data map[interface{}]interface{}, err error) {
	panic("implement me")
}

func (r Redis) Keys(ctx context.Context) (keys []interface{}, err error) {
	panic("implement me")
}

func (r Redis) Values(ctx context.Context) (values []interface{}, err error) {
	panic("implement me")
}

func (r Redis) Update(ctx context.Context, key interface{}, value interface{}) (oldValue *gvar.Var, exist bool, err error) {
	panic("implement me")
}

func (r Redis) UpdateExpire(ctx context.Context, key interface{}, duration time.Duration) (oldDuration time.Duration, err error) {
	panic("implement me")
}

func (r Redis) GetExpire(ctx context.Context, key interface{}) (time.Duration, error) {
	panic("implement me")
}

func (r Redis) Remove(ctx context.Context, keys ...interface{}) (lastValue *gvar.Var, err error) {
	panic("implement me")
}

func (r Redis) Clear(ctx context.Context) error {
	panic("implement me")
}

func (r Redis) Close(ctx context.Context) error {
	panic("implement me")
}



