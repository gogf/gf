package gredis

import (
	"github.com/gogf/gf/container/garray"
	"github.com/gogf/gf/util/gconv"
)

func typeInt64(i interface{}, err error) (int64, error) {
	if err != nil {
		return 0, err
	}
	return gconv.Int64(i), nil
}

func typeInt(i interface{}, err error) (int, error) {
	if err != nil {
		return 0, err
	}
	return gconv.Int(i), nil

}
func typeFloat64(i interface{}, err error) (float64, error) {
	if err != nil {
		return 0, err
	}
	return gconv.Float64(i), nil
}

func typeString(i interface{}, err error) (string, error) {
	if err != nil {
		return "", err
	}
	return gconv.String(i), nil
}

func typeStrings(i interface{}, err error) ([]string, error) {
	if err != nil {
		return nil, err
	}
	return gconv.Strings(i), nil
}

func typeBool(i interface{}, err error) (bool, error) {
	if err != nil {
		return false, err
	}
	return gconv.Bool(i), nil
}

func typeInterfacess(i interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	return gconv.Interfaces(i), nil
}

//==========================================================================key

func (c *Redis) Del(key ...string) (interface{}, error) {
	return c.commnddo("DEL", gconv.Interfaces(key)...)
}

func (c *Redis) Exists(key string) (interface{}, error) {
	return c.commnddo("EXISTS", key)
}

func (c *Redis) Ttl(key string) (interface{}, error) {
	return c.commnddo("TTL", key)
}

func (c *Redis) Expire(key string) (interface{}, error) {
	return c.commnddo("EXPIRE", key)
}

func (c *Redis) Dump(key string) (interface{}, error) {
	return c.commnddo("DUMP", key)
}

func (c *Redis) Expireat(key string, time int64) (interface{}, error) {
	return c.commnddo("EXPIREAT", key, time)
}

func (c *Redis) Keys(key string) (interface{}, error) {
	return c.commnddo("KEYS", key)
}

func (c *Redis) Object(action, key string) (interface{}, error) {
	return c.commnddo("OBJECT", action, key)
}

func (c *Redis) Persist(key string) (interface{}, error) {
	return c.commnddo("PERSIST", key)
}
func (c *Redis) Pttl(key string) (interface{}, error) {
	return c.commnddo("PTTL", key)
}
func (c *Redis) Randomkey() (interface{}, error) {
	return c.commnddo("RANDOMKEY")
}

func (c *Redis) Rename(oldkey, newkey string) (interface{}, error) {
	return c.commnddo("RENAME", oldkey, newkey)
}

func (c *Redis) Renamenx(oldkey, newkey string) (interface{}, error) {
	return c.commnddo("RENAMENX", oldkey, newkey)
}

func (c *Redis) Restore(key string, ttl int64, serializedvalue string) (interface{}, error) {
	return c.commnddo("RESTORE", key, ttl, serializedvalue)
}

func (c *Redis) Sort(key string, params ...string) (interface{}, error) {
	param := garray.NewStrArrayFrom(params)
	return c.commnddo("SORT", gconv.Interfaces(param.InsertBefore(0, key))...)
}

func (c *Redis) Type(key string) (interface{}, error) {
	return c.commnddo("type", key)
}

//============================================================================string
func (c *Redis) Append(key, value string) (interface{}, error) {
	return c.commnddo("append", key, value)
}

func (c *Redis) Set(key, value string) (interface{}, error) {
	return c.commnddo("set", key, value)
}

func (c *Redis) Get(key string) (interface{}, error) {
	return c.commnddo("get", key)
}

func (c *Redis) Bitcount(key string) (interface{}, error) {
	return c.commnddo("BITCOUNT", key)
}

func (c *Redis) Bitop(params ...string) (interface{}, error) {
	return c.commnddo("BITOP", gconv.Interfaces(params)...)
}

func (c *Redis) BITPOS(key string, bit int, option ...int) (int, error) {
	param := garray.NewIntArrayFrom(option).InsertBefore(0, bit)
	return typeInt(c.commnddo("BITPOS", gconv.Interfaces(param)...))
}

func (c *Redis) Bitfield(key string, option ...interface{}) ([]interface{}, error) {
	param := garray.NewArrayFrom(option).InsertBefore(0, key)
	return typeInterfacess(c.commnddo("BITFIELD", gconv.Interfaces(param)...))
}

func (c *Redis) Decr(key string) (interface{}, error) {
	return c.commnddo("DECR", key)
}

func (c *Redis) Decrby(key string, decrement int64) (interface{}, error) {
	return c.commnddo("DECRBY", key, decrement)
}

func (c *Redis) Getbit(key string, offset int64) (interface{}, error) {
	return c.commnddo("GETBIT", key, offset)
}

func (c *Redis) Getrange(key string, start, end int) (interface{}, error) {
	return c.commnddo("GETRANGE", key, start, end)
}

func (c *Redis) Getset(key string, value string) (interface{}, error) {
	return c.commnddo("GETSET", key, value)
}

func (c *Redis) Incr(key string) (interface{}, error) {
	return c.commnddo("INCR", key)
}

func (c *Redis) Incrby(key string, increment int64) (interface{}, error) {
	return c.commnddo("INCRBY", key, increment)
}

func (c *Redis) Incrbyfloat(key string, increment float64) (interface{}, error) {
	return c.commnddo("INCRBYFLOAT", key, increment)
}

func (c *Redis) Mget(key ...string) (interface{}, error) {

	return c.commnddo("MGET", gconv.Interfaces(key)...)
}

func (c *Redis) Mset(params ...string) (interface{}, error) {
	return c.commnddo("MSET", gconv.Interfaces(params)...)
}

func (c *Redis) Msetnx(params ...string) (interface{}, error) {

	return c.commnddo("MSETNX", gconv.Interfaces(params)...)
}

func (c *Redis) PSETEX(key string, milliseconds int64, value string) (interface{}, error) {
	return c.commnddo("PSETEX", key, milliseconds, value)
}

func (c *Redis) Setbit(key string, offset int, value string) (interface{}, error) {
	return c.commnddo("SETBIT", key, offset, value)
}

func (c *Redis) Setex(key string, seconds int64, value string) (interface{}, error) {
	return c.commnddo("SETEX", key, seconds, value)
}

func (c *Redis) Setnx(key string, value string) (interface{}, error) {
	return c.commnddo("SETNX", key, value)
}

func (c *Redis) Setrange(key string, offset int, value string) (interface{}, error) {
	return c.commnddo("SETRANGE", key, offset, value)
}

func (c *Redis) Strlen(key string) (interface{}, error) {
	return c.commnddo("STRLEN", key)
}

//=======================================================================Hash
func (c *Redis) Hset(key, fieldname string, value interface{}) (interface{}, error) {
	return c.commnddo("HSET", key, fieldname, fieldname)
}

func (c *Redis) Hsetnx(key, fieldname string, value interface{}) (interface{}, error) {
	return c.commnddo("HSETNX", key, fieldname, fieldname)
}

func (c *Redis) Hget(key, fieldname string) (interface{}, error) {
	return c.commnddo("HGET", key, fieldname)
}

func (c *Redis) Hexists(key, fieldname string) (bool, error) {
	return typeBool(c.commnddo("HEXISTS", key, fieldname))
}

func (c *Redis) Hdel(key string, fields ...string) (int64, error) {
	param := garray.NewStrArrayFrom(fields)
	return typeInt64(c.commnddo("HDEL", gconv.Interfaces(param.InsertBefore(0, key))...))
}

func (c *Redis) Hlen(key string) (int64, error) {
	return typeInt64(c.commnddo("HLEN", key))
}

func (c *Redis) Hstrlen(key, field string) (int64, error) {
	return typeInt64(c.commnddo("HSTRLEN", key, field))

}

func (c *Redis) Hincrby(key, field string, increment int64) (int64, error) {
	return typeInt64(c.commnddo("HINCRBY", key, field, increment))
}

func (c *Redis) Hincrbyfloat(key, field string, increment float64) (float64, error) {
	return typeFloat64(c.commnddo("HINCRBYFLOAT", key, field, increment))
}

func (c *Redis) HMSET(key string, params ...interface{}) (string, error) {
	param := garray.NewArrayFrom(params)
	return typeString(c.commnddo("HMSET", gconv.Interfaces(param.InsertBefore(0, key))...))
}

func (c *Redis) HMGET(keys ...string) (interface{}, error) {
	return c.commnddo("HMGET", keys)
}

func (c *Redis) Hkeys(key string) ([]string, error) {
	return typeStrings(c.commnddo("HKEYS", key))
}

func (c *Redis) Hvals(key string) (interface{}, error) {
	return c.commnddo("HVALS", key)
}

func (c *Redis) Hgetall(key string) (interface{}, error) {
	return c.commnddo("HGETALL", key)
}

//==============================================================================list
func (c *Redis) LPUSH(key string, values ...interface{}) (int64, error) {
	param := garray.NewArrayFrom(values)
	return typeInt64(c.commnddo("LPUSH", gconv.Interfaces(param.InsertBefore(0, key))...))
}

func (c *Redis) Lpushx(key string, values interface{}) (int64, error) {
	return typeInt64(c.commnddo("LPUSHX", key, values))
}

func (c *Redis) Rpush(key string, values ...interface{}) (int64, error) {
	param := garray.NewArrayFrom(values)
	return typeInt64(c.commnddo("RPUSH", gconv.Interfaces(param.InsertBefore(0, key))...))
}

func (c *Redis) Rpushx(key string, values interface{}) (int64, error) {
	return typeInt64(c.commnddo("RPUSHX", key, values))
}

func (c *Redis) Lpop(key string) (interface{}, error) {
	return c.commnddo("LPOP", key)
}

func (c *Redis) Rpop(key string) (interface{}, error) {
	return c.commnddo("RPOP", key)
}

func (c *Redis) Rpoplpush(key string, source, destination interface{}) (interface{}, error) {
	return c.commnddo("RPOPLPUSH", key, source, destination)
}

func (c *Redis) Lrem(key string, count int, value interface{}) (int64, error) {
	return typeInt64(c.commnddo("LREM", key, count, value))
}

func (c *Redis) Llen(key string) (int64, error) {
	return typeInt64(c.commnddo("LLEN", key))
}

func (c *Redis) Lindex(key string, index int64) (interface{}, error) {
	return c.commnddo("LINDEX", key, index)
}

func (c *Redis) Linsert(key, layout, pivot string, value interface{}) (int64, error) {
	return typeInt64(c.commnddo("LINSERT", key, layout, pivot, value))
}

func (c *Redis) Lset(key, string, index int64, value interface{}) (string, error) {
	return typeString(c.commnddo("LSET", key, index, value))
}

func (c *Redis) Lrange(key, string, start, stop int64) (interface{}, error) {
	return c.commnddo("LRANGE", key, start, stop)
}

func (c *Redis) Blpop(key string, params ...interface{}) (interface{}, error) {
	param := garray.NewArrayFrom(params)
	return c.commnddo("BLPOP", gconv.Interfaces(param.InsertBefore(0, key))...)
}

func (c *Redis) Brpop(key string, params ...interface{}) (interface{}, error) {
	param := garray.NewArrayFrom(params)
	return c.commnddo("BRPOP", gconv.Interfaces(param.InsertBefore(0, key))...)
}

func (c *Redis) Brpoplpush(key, source, destination string, timeout int) (interface{}, error) {
	return c.commnddo("BRPOPLPUSH", key, source, destination, timeout)
}

//========================================================================================set
func (c *Redis) Sadd(key string, members ...interface{}) (int64, error) {
	param := garray.NewArrayFrom(members)
	return typeInt64(c.commnddo("SADD", gconv.Interfaces(param.InsertBefore(0, key))...))
}

func (c *Redis) Sismember(key, member string) (bool, error) {
	return typeBool(c.commnddo("SISMEMBER", key, member))
}

func (c *Redis) Spop(key string) (interface{}, error) {
	return c.commnddo("SPOP", key)
}

func (c *Redis) Srandmember(key string, count ...int) (interface{}, error) {
	return c.commnddo("SRANDMEMBER", key, count[0])
}

func (c *Redis) Srem(keys ...string) (int64, error) {
	return typeInt64(c.commnddo("SREM", gconv.Interfaces(keys)...))
}

func (c *Redis) Smove(source, destination, member string) (bool, error) {
	return typeBool(c.commnddo("SMOVE", source, destination, member))
}

func (c *Redis) Scard(key string) (int64, error) {
	return typeInt64(c.commnddo("SCARD ", key))
}

func (c *Redis) Smembers(key string) (interface{}, error) {
	return c.commnddo("SMEMBERS ", key)
}

func (c *Redis) Sinter(keys ...string) (interface{}, error) {
	return c.commnddo("SINTER ", gconv.Interfaces(keys)...)
}

func (c *Redis) Sinterstore(destination string, key string, keys ...string) (int64, error) {
	param := garray.NewStrArrayFrom(keys)
	param = param.InsertBefore(0, key).InsertBefore(0, destination)
	return typeInt64(c.commnddo("SINTERSTORE ", gconv.Interfaces(param)...))
}

func (c *Redis) SUNION(key string, keys ...string) (interface{}, error) {
	param := garray.NewStrArrayFrom(keys)
	param = param.InsertBefore(0, key)
	return c.commnddo("SUNION ", gconv.Interfaces(param)...)
}

func (c *Redis) Sunionstore(destination string, key string, keys ...string) (int64, error) {
	param := garray.NewStrArrayFrom(keys)
	param = param.InsertBefore(0, key).InsertBefore(0, destination)
	return typeInt64(c.commnddo("SUNIONSTORE ", gconv.Interfaces(param)...))
}

func (c *Redis) Sdiff(key string, keys ...string) (interface{}, error) {
	param := garray.NewStrArrayFrom(keys)
	param = param.InsertBefore(0, key)
	return c.commnddo("SDIFF ", gconv.Interfaces(param)...)
}

func (c *Redis) Sdiffstore(destination string, key string, keys ...string) (int64, error) {
	param := garray.NewStrArrayFrom(keys)
	param = param.InsertBefore(0, key).InsertBefore(0, destination)
	return typeInt64(c.commnddo("SDIFFSTORE ", gconv.Interfaces(param)...))
}

//======================================================================================zset

func (c *Redis) Zadd(params ...interface{}) (int64, error) {
	return typeInt64(c.commnddo("ZADD ", params...))
}

func (c *Redis) Zscore(key string, member interface{}) (string, error) {
	return typeString(c.commnddo("ZSCORE ", key, member))
}

func (c *Redis) Zincrby(key string, increment int, member interface{}) (string, error) {
	return typeString(c.commnddo("ZINCRBY ", key, increment, member))
}

func (c *Redis) ZCARD(key string) (int64, error) {
	return typeInt64(c.commnddo("ZCARD ", key))
}

func (c *Redis) Zcount(key string, min, max int64) (int64, error) {
	return typeInt64(c.commnddo("ZCOUNT ", min, max))
}

func (c *Redis) ZRANGE(key string, start, stop int64, param ...string) (interface{}, error) {
	return c.commnddo("ZRANGE ", start, stop, param[0])
}

func (c *Redis) Zrevrange(key string, start, stop int64, options ...string) (interface{}, error) {
	return c.commnddo("ZREVRANGE ", start, stop, options[0])
}

func (c *Redis) Zrangbyscore(key string, start, stop int64, options ...string) (interface{}, error) {
	return c.commnddo("ZRANGEBYSCORE ", start, stop, options[0])
}

func (c *Redis) Zrevrangebyscore(key string, start, stop int64, options ...string) (interface{}, error) {
	return c.commnddo("ZREVRANGEBYSCORE ", start, stop, options[0])
}

func (c *Redis) Zrank(key, member string) (int64, error) {
	return typeInt64(c.commnddo("ZRANK ", member))
}

func (c *Redis) ZREVRANK(key, member string) (int64, error) {
	return typeInt64(c.commnddo("ZREVRANK ", member))
}

func (c *Redis) ZREM(key string, member ...string) (int64, error) {
	param := garray.NewStrArrayFrom(member)
	param = param.InsertBefore(0, key)
	return typeInt64(c.commnddo("ZREM ", param))
}

func (c *Redis) Zremrangebyrank(key string, start, stop int64) (int64, error) {
	return typeInt64(c.commnddo("ZREMRANGEBYRANK ", key, start, stop))
}

func (c *Redis) Zremrangebyscore(key string, min, max int64) (int64, error) {
	return typeInt64(c.commnddo("ZREMRANGEBYSCORE ", key, min, max))
}

func (c *Redis) Zrangebylex(key, min, max string, options ...string) ([]interface{}, error) {
	param := garray.NewStrArrayFrom(options)
	param = param.InsertBefore(0, max).InsertBefore(0, min).InsertBefore(0, key)
	return typeInterfacess(c.commnddo("ZRANGEBYLEX ", param))
}

func (c *Redis) Zlexcount(key, min, max string) (int64, error) {
	return typeInt64(c.commnddo("ZLEXCOUNT ", key, min, max))
}

func (c *Redis) Zremrangebylex(key, min, max string) (int64, error) {
	return typeInt64(c.commnddo("ZREMRANGEBYLEX ", key, min, max))
}

func (c *Redis) Zunionstore(options ...interface{}) (int64, error) {
	return typeInt64(c.commnddo("ZUNIONSTORE ", options...))
}
func (c *Redis) Zinterstore(options ...interface{}) (int64, error) {
	return typeInt64(c.commnddo("ZINTERSTORE ", options...))
}

//================================================================HyperLogLog
func (c *Redis) Pfadd(key string, options ...interface{}) (bool, error) {
	param := garray.NewArrayFrom(options)
	param = param.InsertBefore(0, key)
	return typeBool(c.commnddo("PFADD ", gconv.Interfaces(param)...))
}

func (c *Redis) Pfcount(keys ...string) (int64, error) {
	return typeInt64(c.commnddo("PFCOUNT ", gconv.Interfaces(keys)...))
}

func (c *Redis) Pfmerge(keys ...string) (string, error) {
	return typeString(c.commnddo("PFMERGE ", gconv.Interfaces(keys)...))
}

//================================================================================GEO
func (c *Redis) Geoadd(key string, params ...interface{}) (int64, error) {
	param := garray.NewArrayFrom(params)
	param = param.InsertBefore(0, key)
	return typeInt64(c.commnddo("GEOADD ", gconv.Interfaces(param)...))
}

func (c *Redis) Geopos(key string, member ...interface{}) ([]interface{}, error) {
	param := garray.NewArrayFrom(member)
	param = param.InsertBefore(0, key)
	return typeInterfacess(c.commnddo("GEOPOS ", gconv.Interfaces(param)...))
}

func (c *Redis) Geodist(key string, params ...string) (interface{}, error) {
	param := garray.NewStrArrayFrom(params).InsertBefore(0, key)
	return c.commnddo("GEODIST ", gconv.Interfaces(param)...)
}

func (c *Redis) Georadius(key string, member ...interface{}) ([]interface{}, error) {
	param := garray.NewArrayFrom(member).InsertBefore(0, key)
	return typeInterfacess(c.commnddo("GEORADIUS ", gconv.Interfaces(param)...))
}

func (c *Redis) Georadiusbymember(key string, member ...interface{}) ([]interface{}, error) {
	param := garray.NewArrayFrom(member).InsertBefore(0, key)
	return typeInterfacess(c.commnddo("GEORADIUSBYMEMBER ", gconv.Interfaces(param)...))
}

func (c *Redis) Geohash(key string, member ...interface{}) ([]interface{}, error) {
	param := garray.NewArrayFrom(member).InsertBefore(0, key)
	return typeInterfacess(c.commnddo("GEOHASH ", gconv.Interfaces(param)...))
}

//============================================================================channel
func (c *Redis) Publist(channel, message string) (int, error) {
	return typeInt(c.commnddo("PUBLISH ", channel, message))
}

func (c *Redis) Subscribe(channel ...string) (interface{}, error) {
	return c.commnddo("SUBSCRIBE", gconv.Interfaces(channel)...)
}

func (c *Redis) Psubscribe(pattern ...string) (interface{}, error) {
	return c.commnddo("PSUBSCRIBE", gconv.Interfaces(pattern)...)
}

func (c *Redis) Unsubscribe(pattern ...string) (interface{}, error) {
	return c.commnddo("UNSUBSCRIBE", gconv.Interfaces(pattern)...)
}

func (c *Redis) Pubsubscribe(pattern ...string) (interface{}, error) {
	return c.commnddo("PUNSUBSCRIBE", gconv.Interfaces(pattern)...)
}
