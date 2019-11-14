package gredis

import (
	"errors"
	"github.com/gogf/gf/util/gconv"
	"reflect"
)

type GeoLocation struct {
	Name                string
	Longitude, Latitude string
	GeoHash             int64
	Dist                string
}

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
func typeStringss(i interface{}, err error) ([][]string, error) {
	if err != nil {
		return nil, err
	}
	ss := [][]string{}
	is := gconv.Interfaces(i)
	for _, v := range is {
		//fmt.Println(gconv.Strings(v))
		ss = append(ss, gconv.Strings(v))
	}

	return ss, nil
}

func typeGeoLocation(i interface{}, err error) ([]*GeoLocation, error) {
	if err != nil {
		return nil, err
	}
	var loc GeoLocation
	ss := []*GeoLocation{}
	is := gconv.Interfaces(i)
	for _, v := range is {
		s1 := gconv.Strings(v)
		loc.Longitude = s1[0]
		loc.Latitude = s1[1]
		ss = append(ss, &loc)
	}

	return ss, nil
}

func typeGeoLocationd(i interface{}, err error) ([]*GeoLocation, error) {
	if err != nil {
		return nil, err
	}

	var loc GeoLocation
	ss := []*GeoLocation{}
	is := gconv.Interfaces(i)

	for _, v := range is {
		if reflect.TypeOf(v).String() == "[]uint8" {
			loc.Name = gconv.String(v)
			ss = append(ss, &loc)
			continue
		}
		s1 := gconv.Interfaces(v)
		s1_length := len(s1)
		if s1_length == 3 {

			loc.Name = gconv.String(s1[0])
			loc.Dist = gconv.String(s1[1])
			s1_3 := gconv.Strings(s1[2])
			loc.Longitude = s1_3[0]
			loc.Latitude = s1_3[1]

		} else if s1_length == 2 {

			loc.Name = gconv.String(s1[0])
			if s1_2, ok := s1[1].(string); ok == true {
				loc.Dist = s1_2
			} else {
				s1_2s := gconv.Strings(s1[1])
				loc.Longitude = s1_2s[0]
				loc.Latitude = s1_2s[1]
			}
		}

		ss = append(ss, &loc)
	}

	return ss, nil
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
func (c *Redis) Del(key ...string) (int, error) {
	return typeInt(c.commandDo("DEL", gconv.Interfaces(key)...))
}

func (c *Redis) Exists(key string) (int, error) {
	return typeInt(c.commandDo("EXISTS", key))
}

func (c *Redis) Ttl(key string) (int64, error) {
	return typeInt64(c.commandDo("TTL", key))
}

func (c *Redis) Expire(key string, time int64) (int64, error) {
	return typeInt64(c.commandDo("EXPIRE", key, time))
}

func (c *Redis) Dump(key string) (string, error) {
	return typeString(c.commandDo("DUMP", key))
}

func (c *Redis) ExpireAt(key string, timestamp int64) (int, error) {
	return typeInt(c.commandDo("EXPIREAT", key, timestamp))
}

// Returns all keys matching pattern, but not for clustering
func (c *Redis) Keys(key string) ([]interface{}, error) {
	return typeInterfacess(c.commandDo("KEYS", key))
}

func (c *Redis) Object(action, key string) (interface{}, error) {
	return c.commandDo("OBJECT", action, key)
}

func (c *Redis) Persist(key string) (int, error) {
	return typeInt(c.commandDo("PERSIST", key))
}
func (c *Redis) PTTL(key string) (int64, error) {
	return typeInt64(c.commandDo("PTTL", key))
}
func (c *Redis) RandomKey() (interface{}, error) {
	return c.commandDo("RANDOMKEY")
}

func (c *Redis) Rename(oldkey, newkey string) (string, error) {
	return typeString(c.commandDo("RENAME", oldkey, newkey))
}

func (c *Redis) RenameNX(key, newkey string) (int, error) {
	return typeInt(c.commandDo("RENAMENX", key, newkey))
}

func (c *Redis) ReStore(key string, ttl int64, serializedValue string, replace ...string) (string, error) {
	str1 := ""
	if len(replace) > 0 {
		str1 = replace[0]
	}
	return typeString(c.commandDo("RESTORE", key, ttl, serializedValue, str1))
}

func (c *Redis) Sort(key string, params ...interface{}) ([]interface{}, error) {
	return typeInterfacess(c.commandDo("SORT", append([]interface{}{key}, params...)...))
}

func (c *Redis) Type(key string) (string, error) {
	return typeString(c.commandDo("type", key))
}

//============================================================================string
func (c *Redis) Append(key, value string) (int64, error) {
	return typeInt64(c.commandDo("append", key, value))
}

func (c *Redis) Set(key, value string) (interface{}, error) {
	return c.commandDo("set", key, value)
}

func (c *Redis) Get(key string) (string, error) {
	return typeString(c.commandDo("get", key))
}

func (c *Redis) BitCount(key string) (int, error) {
	return typeInt(c.commandDo("BITCOUNT", key))
}

func (c *Redis) BiTop(params ...string) (int, error) {
	return typeInt(c.commandDo("BITOP", gconv.Interfaces(params)...))
}

func (c *Redis) BitPos(key string, bit int, option ...int) (int, error) {
	return typeInt(c.commandDo("BITPOS", append([]interface{}{key, bit}, gconv.Interfaces(option)...)...))
}

func (c *Redis) BitField(option string) ([]interface{}, error) {
	return typeInterfacess(c.commandDo("BITFIELD", option))
}

func (c *Redis) Decr(key string) (int64, error) {
	return typeInt64(c.commandDo("DECR", key))
}

func (c *Redis) DecrBy(key string, decrement int64) (int64, error) {
	return typeInt64(c.commandDo("DECRBY", key, decrement))
}

func (c *Redis) GetBit(key string, offset int) (int, error) {
	return typeInt(c.commandDo("GETBIT", key, offset))
}

func (c *Redis) GetRange(key string, start, end int) (string, error) {
	return typeString(c.commandDo("GETRANGE", key, start, end))
}

func (c *Redis) GetSet(key string, value string) (string, error) {
	return typeString(c.commandDo("GETSET", key, value))
}

func (c *Redis) Incr(key string) (int64, error) {
	return typeInt64(c.commandDo("INCR", key))
}

func (c *Redis) IncrBy(key string, increment int64) (int64, error) {
	return typeInt64(c.commandDo("INCRBY", key, increment))
}

func (c *Redis) IncrByFloat(key string, increment float64) (string, error) {
	return typeString(c.commandDo("INCRBYFLOAT", key, increment))
}

func (c *Redis) MGet(key ...string) ([]string, error) {
	if len(key) < 1 {
		return nil, errors.New("there must be one key's name")
	}
	return typeStrings(c.commandDo("MGET", gconv.Interfaces(key)...))
}

func (c *Redis) MSet(params ...string) (string, error) {
	if len(params) < 2 {
		return "", errors.New("there must be one k-v ")
	}
	return typeString(c.commandDo("MSET", gconv.Interfaces(params)...))
}

func (c *Redis) MSetNx(params ...string) (int, error) {

	return typeInt(c.commandDo("MSETNX", gconv.Interfaces(params)...))
}

func (c *Redis) PSetEx(key string, milliseconds int64, value string) (string, error) {
	return typeString(c.commandDo("PSETEX", key, milliseconds, value))
}

func (c *Redis) SetBit(key string, offset, value int) (int, error) {
	return typeInt(c.commandDo("SETBIT", key, offset, value))
}

func (c *Redis) SetEx(key string, seconds int64, value string) (string, error) {
	return typeString(c.commandDo("SETEX", key, seconds, value))
}

func (c *Redis) SetNx(key string, value string) (int, error) {
	return typeInt(c.commandDo("SETNX", key, value))
}

func (c *Redis) SetRange(key string, offset int, value string) (int, error) {
	return typeInt(c.commandDo("SETRANGE", key, offset, value))
}

func (c *Redis) StrLen(key string) (int, error) {
	return typeInt(c.commandDo("STRLEN", key))
}

//=======================================================================Hash
func (c *Redis) HSet(key, fieldname string, value interface{}) (int, error) {
	return typeInt(c.commandDo("HSET", key, fieldname, value))
}

func (c *Redis) HSetNx(key, fieldname string, value interface{}) (int, error) {
	return typeInt(c.commandDo("HSETNX", key, fieldname, value))
}

func (c *Redis) HGet(key, fieldname string) (string, error) {
	return typeString(c.commandDo("HGET", key, fieldname))
}

func (c *Redis) HExists(key, fieldname string) (int, error) {
	return typeInt(c.commandDo("HEXISTS", key, fieldname))
}

func (c *Redis) HDel(key string, fields ...string) (int, error) {
	if len(fields) < 1 {
		return 0, errors.New("must have one field's name")
	}
	return typeInt(c.commandDo("HDEL", gconv.Interfaces(append([]string{key}, fields...))...))
}

func (c *Redis) HLen(key string) (int, error) {
	return typeInt(c.commandDo("HLEN", key))
}

func (c *Redis) HStrLen(key, field string) (int, error) {
	return typeInt(c.commandDo("HSTRLEN", key, field))

}

func (c *Redis) HIncrBy(key, field string, increment int64) (int64, error) {
	return typeInt64(c.commandDo("HINCRBY", key, field, increment))
}

func (c *Redis) HIncrByFloat(key, field string, increment float64) (string, error) {
	return typeString(c.commandDo("HINCRBYFLOAT", key, field, increment))
}

func (c *Redis) HMSet(key string, params ...interface{}) (string, error) {

	return typeString(c.commandDo("HMSET", append([]interface{}{key}, params...)...))
}

func (c *Redis) HMGet(key string, option ...string) ([]string, error) {
	return typeStrings(c.commandDo("HMGET", gconv.Interfaces(append([]string{key}, option...))...))
}

func (c *Redis) HKeys(key string) ([]string, error) {
	return typeStrings(c.commandDo("HKEYS", key))
}

func (c *Redis) HVals(key string) ([]string, error) {
	return typeStrings(c.commandDo("HVALS", key))
}

func (c *Redis) HGetAll(key string) ([]string, error) {
	return typeStrings(c.commandDo("HGETALL", key))
}

//==============================================================================list
func (c *Redis) LPush(key string, values ...interface{}) (int64, error) {
	return typeInt64(c.commandDo("LPUSH", append([]interface{}{key}, values...)...))

}

func (c *Redis) LPushX(key string, values interface{}) (int64, error) {
	return typeInt64(c.commandDo("LPUSHX", key, values))
}

func (c *Redis) RPush(key string, values ...interface{}) (int64, error) {

	return typeInt64(c.commandDo("RPUSH", append([]interface{}{key}, values...)...))
}

func (c *Redis) RPushX(key string, values interface{}) (int64, error) {
	return typeInt64(c.commandDo("RPUSHX", key, values))
}

func (c *Redis) LPop(key string) (string, error) {
	return typeString(c.commandDo("LPOP", key))
}

func (c *Redis) RPop(key string) (string, error) {
	return typeString(c.commandDo("RPOP", key))
}

func (c *Redis) RPoplPush(source, destination string) (string, error) {
	return typeString(c.commandDo("RPOPLPUSH", source, destination))
}

func (c *Redis) LRem(key string, count int, value interface{}) (int64, error) {
	return typeInt64(c.commandDo("LREM", key, count, value))
}

func (c *Redis) LLen(key string) (int64, error) {
	return typeInt64(c.commandDo("LLEN", key))
}

func (c *Redis) LIndex(key string, index int64) (string, error) {
	return typeString(c.commandDo("LINDEX", key, index))
}

func (c *Redis) LInsert(key, layout, pivot string, value interface{}) (int64, error) {
	return typeInt64(c.commandDo("LINSERT", key, layout, pivot, value))
}

func (c *Redis) LSet(key string, index int64, value interface{}) (string, error) {
	return typeString(c.commandDo("LSET", key, index, value))
}

func (c *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return typeStrings(c.commandDo("LRANGE", key, start, stop))
}

func (c *Redis) BlPop(key string, params ...interface{}) ([]string, error) {
	return typeStrings(c.commandDo("BLPOP", append([]interface{}{key}, params...)...))
}

func (c *Redis) BrPop(key string, params ...interface{}) ([]string, error) {
	return typeStrings(c.commandDo("BRPOP", append([]interface{}{key}, params...)...))
}

func (c *Redis) BrPopLPush(source, destination string, timeout int) ([]string, error) {
	return typeStrings(c.commandDo("BRPOPLPUSH", source, destination, timeout))
}

//========================================================================================set
func (c *Redis) SAdd(key string, members ...interface{}) (int64, error) {

	return typeInt64(c.commandDo("SADD", append([]interface{}{key}, members...)...))
}

func (c *Redis) SisMember(key, member string) (int, error) {
	return typeInt(c.commandDo("SISMEMBER", key, member))
}

func (c *Redis) SPop(key string) (string, error) {
	return typeString(c.commandDo("SPOP", key))
}

func (c *Redis) SRandMember(key string, count ...int) ([]string, error) {
	if len(count) == 0 {
		return typeStrings(c.commandDo("SRANDMEMBER", key, 1))
	}
	return typeStrings(c.commandDo("SRANDMEMBER", key, count[0]))
}

func (c *Redis) SRem(key string, members ...string) (int, error) {
	return typeInt(c.commandDo("SREM", append([]interface{}{key}, gconv.Interfaces(members)...)...))
}

func (c *Redis) SMove(source, destination, member string) (int, error) {
	return typeInt(c.commandDo("SMOVE", source, destination, member))
}

func (c *Redis) SCard(key string) (int64, error) {
	return typeInt64(c.commandDo("SCARD", key))
}

func (c *Redis) SMembers(key string) ([]string, error) {
	return typeStrings(c.commandDo("SMEMBERS", key))
}

func (c *Redis) SInter(keys ...string) ([]string, error) {
	if len(keys) == 0 {
		return nil, errors.New("must have a key")
	}
	return typeStrings(c.commandDo("SINTER", gconv.Interfaces(keys)...))
}

func (c *Redis) SInterStore(destination string, key string, keys ...string) (int64, error) {
	return typeInt64(c.commandDo("SINTERSTORE", append([]interface{}{destination, key}, gconv.Interfaces(keys)...)...))
}

func (c *Redis) SUnion(key string, keys ...string) ([]string, error) {
	return typeStrings(c.commandDo("SUNION", append([]interface{}{key}, gconv.Interfaces(keys)...)...))
}

func (c *Redis) SUnionStore(destination string, key string, keys ...string) (int64, error) {
	return typeInt64(c.commandDo("SUNIONSTORE", append([]interface{}{destination, key}, gconv.Interfaces(keys)...)...))
}

func (c *Redis) SDiff(key string, keys ...string) ([]string, error) {
	return typeStrings(c.commandDo("SDIFF", append([]interface{}{key}, gconv.Interfaces(keys)...)...))
}

func (c *Redis) SDiffStore(destination string, key string, keys ...string) (int64, error) {

	return typeInt64(c.commandDo("SDIFFSTORE", append([]interface{}{destination, key}, gconv.Interfaces(keys)...)...))
}

//======================================================================================zset

func (c *Redis) ZAdd(params ...interface{}) (int, error) {
	return typeInt(c.commandDo("ZADD", params...))
}

func (c *Redis) ZScore(key string, member interface{}) (string, error) {
	return typeString(c.commandDo("ZSCORE", key, member))
}

func (c *Redis) ZinCrBy(key string, increment float64, member interface{}) (string, error) {
	return typeString(c.commandDo("ZINCRBY", key, increment, member))
}

func (c *Redis) ZCard(key string) (int64, error) {
	return typeInt64(c.commandDo("ZCARD", key))
}

func (c *Redis) ZCount(key string, min, max int64) (int64, error) {
	return typeInt64(c.commandDo("ZCOUNT", key, min, max))
}

func (c *Redis) ZRange(key string, start, stop int64, param ...string) ([]string, error) {
	if len(param) == 0 {
		return typeStrings(c.commandDo("ZRANGE", key, start, stop))
	}
	return typeStrings(c.commandDo("ZRANGE", key, start, stop, param[0]))
}

func (c *Redis) ZRevRange(key string, start, stop int64, param ...string) ([]string, error) {
	if len(param) == 0 {
		return typeStrings(c.commandDo("ZRANGE", key, start, stop))
	}
	return typeStrings(c.commandDo("ZREVRANGE", key, start, stop, param[0]))
}

func (c *Redis) ZRangeByScore(key, min, max string, options ...interface{}) ([]string, error) {

	return typeStrings(c.commandDo("ZRANGEBYSCORE", append([]interface{}{key, min, max}, options...)...))
}

func (c *Redis) ZRevRangeByScore(key string, min, max string, options ...interface{}) ([]string, error) {
	return typeStrings(c.commandDo("ZREVRANGEBYSCORE", append([]interface{}{key, min, max}, options...)...))
}

func (c *Redis) ZRank(key, member string) (int64, error) {
	return typeInt64(c.commandDo("ZRANK", key, member))
}

func (c *Redis) ZRevRank(key, member string) (int64, error) {
	return typeInt64(c.commandDo("ZREVRANK", key, member))
}

func (c *Redis) ZRem(key string, member ...interface{}) (int, error) {
	if len(member) == 0 {
		return 0, errors.New("must have an one key")
	}
	return typeInt(c.commandDo("ZREM", append([]interface{}{key}, member...)...))
}

func (c *Redis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return typeInt64(c.commandDo("ZREMRANGEBYRANK", key, start, stop))
}

func (c *Redis) ZRemRangeByScore(key string, min, max float64) (int64, error) {
	return typeInt64(c.commandDo("ZREMRANGEBYSCORE", key, min, max))
}

func (c *Redis) ZRangeByLex(key, min, max string, options ...interface{}) ([]string, error) {
	return typeStrings(c.commandDo("ZRANGEBYLEX", append([]interface{}{key, min, max}, options...)...))
}

func (c *Redis) ZLexCount(key, min, max string) (int64, error) {
	return typeInt64(c.commandDo("ZLEXCOUNT", key, min, max))
}

func (c *Redis) ZRemRangeByLex(key, min, max string) (int64, error) {
	return typeInt64(c.commandDo("ZREMRANGEBYLEX", key, min, max))
}

func (c *Redis) ZUnionStore(options ...interface{}) (int64, error) {
	if len(options) < 3 {
		return 0, errors.New("there must be three parameters")
	}
	return typeInt64(c.commandDo("ZUNIONSTORE", options...))
}
func (c *Redis) ZInterStore(options ...interface{}) (int64, error) {
	return typeInt64(c.commandDo("ZINTERSTORE", options...))
}

//================================================================HyperLogLog
func (c *Redis) PfAdd(key string, options ...interface{}) (int, error) {
	return typeInt(c.commandDo("PFADD", append([]interface{}{key}, options...)...))
}

func (c *Redis) PfCount(keys ...string) (int64, error) {
	return typeInt64(c.commandDo("PFCOUNT", gconv.Interfaces(keys)...))
}

func (c *Redis) PfMerge(keys ...string) (string, error) {
	if len(keys) < 2 {
		return "", errors.New("need at least two keys")
	}
	return typeString(c.commandDo("PFMERGE", gconv.Interfaces(keys)...))
}

//================================================================================GEO
func (c *Redis) GeoAdd(key string, params ...interface{}) (int, error) {
	return typeInt(c.commandDo("GEOADD", append([]interface{}{key}, params...)...))
}

func (c *Redis) GeoPos(key string, member ...interface{}) ([]*GeoLocation, error) {
	return typeGeoLocation(c.commandDo("GEOPOS", append([]interface{}{key}, member...)...))
}

func (c *Redis) GeoDist(key string, params ...string) (string, error) {
	return typeString(c.commandDo("GEODIST", append([]interface{}{key}, gconv.Interfaces(params)...)...))
}

func (c *Redis) GeoRadius(key string, member ...interface{}) ([]*GeoLocation, error) {
	if len(member) < 5 {
		return nil, errors.New("there are must have five keys")
	}
	return typeGeoLocationd(c.commandDo("GEORADIUS", append([]interface{}{key}, member...)...))
}

func (c *Redis) GeoRadiusByMember(key string, member ...interface{}) ([]*GeoLocation, error) {
	return typeGeoLocationd(c.commandDo("GEORADIUSBYMEMBER", append([]interface{}{key}, member...)...))
}

func (c *Redis) GeoHash(key string, member ...interface{}) ([]string, error) {
	return typeStrings(c.commandDo("GEOHASH", append([]interface{}{key}, member...)...))
}

//============================================================================channel
func (c *Redis) Publish(channel, message string) (int, error) {

	return typeInt(c.commandDo("PUBLISH", channel, message))
}

func (c *Redis) PubSub(channel string, member ...interface{}) ([]string, error) {
	return typeStrings(c.commandDo("PUBSUB", append([]interface{}{channel}, member...)...))
}

func (c *Redis) SubScribe(channel ...string) ([]string, error) {
	return typeStrings(c.commandDo("SUBSCRIBE", gconv.Interfaces(channel)...))
}

func (c *Redis) PSubScribe(pattern ...string) ([]string, error) {
	return typeStrings(c.commandDo("PSUBSCRIBE", gconv.Interfaces(pattern)...))
}

func (c *Redis) UnSubScribe(pattern ...string) ([]string, error) {
	return typeStrings(c.commandDo("UNSUBSCRIBE", gconv.Interfaces(pattern)...))
}

func (c *Redis) PunSubScribe(pattern ...string) ([]string, error) {
	return typeStrings(c.commandDo("PUNSUBSCRIBE", gconv.Interfaces(pattern)...))
}
