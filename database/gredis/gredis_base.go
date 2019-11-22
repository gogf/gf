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


func typeInterfacess(i interface{}, err error) ([]interface{}, error) {
	if err != nil {
		return nil, err
	}
	return gconv.Interfaces(i), nil
}

//==========================================================================key
func (r *Redis) Del(key ...string) (int, error) {
	return typeInt(r.commandDo("DEL", gconv.Interfaces(key)...))
}

func (r *Redis) Exists(key string) (int, error) {
	return typeInt(r.commandDo("EXISTS", key))
}

func (r *Redis) Ttl(key string) (int64, error) {
	return typeInt64(r.commandDo("TTL", key))
}

func (r *Redis) Expire(key string, time int64) (int64, error) {
	return typeInt64(r.commandDo("EXPIRE", key, time))
}

func (r *Redis) Dump(key string) (string, error) {
	return typeString(r.commandDo("DUMP", key))
}

func (r *Redis) ExpireAt(key string, timestamp int64) (int, error) {
	return typeInt(r.commandDo("EXPIREAT", key, timestamp))
}

// Returns all keys matching pattern, but not for clustering
func (r *Redis) Keys(key string) ([]interface{}, error) {
	return typeInterfacess(r.commandDo("KEYS", key))
}

func (r *Redis) Object(action, key string) (interface{}, error) {
	return r.commandDo("OBJECT", action, key)
}

func (r *Redis) Persist(key string) (int, error) {
	return typeInt(r.commandDo("PERSIST", key))
}
func (r *Redis) PTTL(key string) (int64, error) {
	return typeInt64(r.commandDo("PTTL", key))
}
func (r *Redis) RandomKey() (interface{}, error) {
	return r.commandDo("RANDOMKEY")
}

func (r *Redis) Rename(oldkey, newkey string) (string, error) {
	return typeString(r.commandDo("RENAME", oldkey, newkey))
}

func (r *Redis) RenameNX(key, newkey string) (int, error) {
	return typeInt(r.commandDo("RENAMENX", key, newkey))
}

func (r *Redis) ReStore(key string, ttl int64, serializedValue string, replace ...string) (string, error) {
	str1 := ""
	if len(replace) > 0 {
		str1 = replace[0]
	}
	return typeString(r.commandDo("RESTORE", key, ttl, serializedValue, str1))
}

func (r *Redis) Sort(key string, params ...interface{}) ([]interface{}, error) {
	return typeInterfacess(r.commandDo("SORT", append([]interface{}{key}, params...)...))
}

func (r *Redis) Type(key string) (string, error) {
	return typeString(r.commandDo("type", key))
}

//============================================================================string
func (r *Redis) Append(key, value string) (int64, error) {
	return typeInt64(r.commandDo("append", key, value))
}

func (r *Redis) Set(key, value string) (interface{}, error) {
	return r.commandDo("set", key, value)
}

func (r *Redis) Get(key string) (string, error) {
	return typeString(r.commandDo("get", key))
}

func (r *Redis) BitCount(key string) (int, error) {
	return typeInt(r.commandDo("BITCOUNT", key))
}

func (r *Redis) BiTop(params ...string) (int, error) {
	return typeInt(r.commandDo("BITOP", gconv.Interfaces(params)...))
}

func (r *Redis) BitPos(key string, bit int, option ...int) (int, error) {
	return typeInt(r.commandDo("BITPOS", append([]interface{}{key, bit}, gconv.Interfaces(option)...)...))
}

func (r *Redis) BitField(option string) ([]interface{}, error) {
	return typeInterfacess(r.commandDo("BITFIELD", option))
}

func (r *Redis) Decr(key string) (int64, error) {
	return typeInt64(r.commandDo("DECR", key))
}

func (r *Redis) DecrBy(key string, decrement int64) (int64, error) {
	return typeInt64(r.commandDo("DECRBY", key, decrement))
}

func (r *Redis) GetBit(key string, offset int) (int, error) {
	return typeInt(r.commandDo("GETBIT", key, offset))
}

func (r *Redis) GetRange(key string, start, end int) (string, error) {
	return typeString(r.commandDo("GETRANGE", key, start, end))
}

func (r *Redis) GetSet(key string, value string) (string, error) {
	return typeString(r.commandDo("GETSET", key, value))
}

func (r *Redis) Incr(key string) (int64, error) {
	return typeInt64(r.commandDo("INCR", key))
}

func (r *Redis) IncrBy(key string, increment int64) (int64, error) {
	return typeInt64(r.commandDo("INCRBY", key, increment))
}

func (r *Redis) IncrByFloat(key string, increment float64) (string, error) {
	return typeString(r.commandDo("INCRBYFLOAT", key, increment))
}

func (r *Redis) MGet(key ...string) ([]string, error) {
	if len(key) < 1 {
		return nil, errors.New("there must be one key's name")
	}
	return typeStrings(r.commandDo("MGET", gconv.Interfaces(key)...))
}

func (r *Redis) MSet(params ...string) (string, error) {
	if len(params) < 2 {
		return "", errors.New("there must be one k-v ")
	}
	return typeString(r.commandDo("MSET", gconv.Interfaces(params)...))
}

func (r *Redis) MSetNx(params ...string) (int, error) {

	return typeInt(r.commandDo("MSETNX", gconv.Interfaces(params)...))
}

func (r *Redis) PSetEx(key string, milliseconds int64, value string) (string, error) {
	return typeString(r.commandDo("PSETEX", key, milliseconds, value))
}

func (r *Redis) SetBit(key string, offset, value int) (int, error) {
	return typeInt(r.commandDo("SETBIT", key, offset, value))
}

func (r *Redis) SetEx(key string, seconds int64, value string) (string, error) {
	return typeString(r.commandDo("SETEX", key, seconds, value))
}

func (r *Redis) SetNx(key string, value string) (int, error) {
	return typeInt(r.commandDo("SETNX", key, value))
}

func (r *Redis) SetRange(key string, offset int, value string) (int, error) {
	return typeInt(r.commandDo("SETRANGE", key, offset, value))
}

func (r *Redis) StrLen(key string) (int, error) {
	return typeInt(r.commandDo("STRLEN", key))
}

//=======================================================================Hash
func (r *Redis) HSet(key, fieldname string, value interface{}) (int, error) {
	return typeInt(r.commandDo("HSET", key, fieldname, value))
}

func (r *Redis) HSetNx(key, fieldname string, value interface{}) (int, error) {
	return typeInt(r.commandDo("HSETNX", key, fieldname, value))
}

func (r *Redis) HGet(key, fieldname string) (string, error) {
	return typeString(r.commandDo("HGET", key, fieldname))
}

func (r *Redis) HExists(key, fieldname string) (int, error) {
	return typeInt(r.commandDo("HEXISTS", key, fieldname))
}

func (r *Redis) HDel(key string, fields ...string) (int, error) {
	if len(fields) < 1 {
		return 0, errors.New("must have one field's name")
	}
	return typeInt(r.commandDo("HDEL", gconv.Interfaces(append([]string{key}, fields...))...))
}

func (r *Redis) HLen(key string) (int, error) {
	return typeInt(r.commandDo("HLEN", key))
}

func (r *Redis) HStrLen(key, field string) (int, error) {
	return typeInt(r.commandDo("HSTRLEN", key, field))

}

func (r *Redis) HIncrBy(key, field string, increment int64) (int64, error) {
	return typeInt64(r.commandDo("HINCRBY", key, field, increment))
}

func (r *Redis) HIncrByFloat(key, field string, increment float64) (string, error) {
	return typeString(r.commandDo("HINCRBYFLOAT", key, field, increment))
}

func (r *Redis) HMSet(key string, params ...interface{}) (string, error) {

	return typeString(r.commandDo("HMSET", append([]interface{}{key}, params...)...))
}

func (r *Redis) HMGet(key string, option ...string) ([]string, error) {
	return typeStrings(r.commandDo("HMGET", gconv.Interfaces(append([]string{key}, option...))...))
}

func (r *Redis) HKeys(key string) ([]string, error) {
	return typeStrings(r.commandDo("HKEYS", key))
}

func (r *Redis) HVals(key string) ([]string, error) {
	return typeStrings(r.commandDo("HVALS", key))
}

func (r *Redis) HGetAll(key string) ([]string, error) {
	return typeStrings(r.commandDo("HGETALL", key))
}

//==============================================================================list
func (r *Redis) LPush(key string, values ...interface{}) (int64, error) {
	return typeInt64(r.commandDo("LPUSH", append([]interface{}{key}, values...)...))

}

func (r *Redis) LPushX(key string, values interface{}) (int64, error) {
	return typeInt64(r.commandDo("LPUSHX", key, values))
}

func (r *Redis) RPush(key string, values ...interface{}) (int64, error) {

	return typeInt64(r.commandDo("RPUSH", append([]interface{}{key}, values...)...))
}

func (r *Redis) RPushX(key string, values interface{}) (int64, error) {
	return typeInt64(r.commandDo("RPUSHX", key, values))
}

func (r *Redis) LPop(key string) (string, error) {
	return typeString(r.commandDo("LPOP", key))
}

func (r *Redis) RPop(key string) (string, error) {
	return typeString(r.commandDo("RPOP", key))
}

func (r *Redis) RPopLPush(source, destination string) (string, error) {
	return typeString(r.commandDo("RPOPLPUSH", source, destination))
}

func (r *Redis) LRem(key string, count int, value interface{}) (int64, error) {
	return typeInt64(r.commandDo("LREM", key, count, value))
}

func (r *Redis) LLen(key string) (int64, error) {
	return typeInt64(r.commandDo("LLEN", key))
}

func (r *Redis) LIndex(key string, index int64) (string, error) {
	return typeString(r.commandDo("LINDEX", key, index))
}

func (r *Redis) LInsert(key, layout, pivot string, value interface{}) (int64, error) {
	return typeInt64(r.commandDo("LINSERT", key, layout, pivot, value))
}

func (r *Redis) LSet(key string, index int64, value interface{}) (string, error) {
	return typeString(r.commandDo("LSET", key, index, value))
}

func (r *Redis) LRange(key string, start, stop int64) ([]string, error) {
	return typeStrings(r.commandDo("LRANGE", key, start, stop))
}

func (r *Redis) BlPop(key string, params ...interface{}) ([]string, error) {
	return typeStrings(r.commandDo("BLPOP", append([]interface{}{key}, params...)...))
}

func (r *Redis) BrPop(key string, params ...interface{}) ([]string, error) {
	return typeStrings(r.commandDo("BRPOP", append([]interface{}{key}, params...)...))
}

func (r *Redis) BrPopLPush(source, destination string, timeout int) ([]string, error) {
	return typeStrings(r.commandDo("BRPOPLPUSH", source, destination, timeout))
}

//========================================================================================set
func (r *Redis) SAdd(key string, members ...interface{}) (int64, error) {

	return typeInt64(r.commandDo("SADD", append([]interface{}{key}, members...)...))
}

func (r *Redis) SisMember(key, member string) (int, error) {
	return typeInt(r.commandDo("SISMEMBER", key, member))
}

func (r *Redis) SPop(key string) (string, error) {
	return typeString(r.commandDo("SPOP", key))
}

func (r *Redis) SRandMember(key string, count ...int) ([]string, error) {
	if len(count) == 0 {
		return typeStrings(r.commandDo("SRANDMEMBER", key, 1))
	}
	return typeStrings(r.commandDo("SRANDMEMBER", key, count[0]))
}

func (r *Redis) SRem(key string, members ...string) (int, error) {
	return typeInt(r.commandDo("SREM", append([]interface{}{key}, gconv.Interfaces(members)...)...))
}

func (r *Redis) SMove(source, destination, member string) (int, error) {
	return typeInt(r.commandDo("SMOVE", source, destination, member))
}

func (r *Redis) SCard(key string) (int64, error) {
	return typeInt64(r.commandDo("SCARD", key))
}

func (r *Redis) SMembers(key string) ([]string, error) {
	return typeStrings(r.commandDo("SMEMBERS", key))
}

func (r *Redis) SInter(keys ...string) ([]string, error) {
	if len(keys) == 0 {
		return nil, errors.New("must have a key")
	}
	return typeStrings(r.commandDo("SINTER", gconv.Interfaces(keys)...))
}

func (r *Redis) SInterStore(destination string, key string, keys ...string) (int64, error) {
	return typeInt64(r.commandDo("SINTERSTORE", append([]interface{}{destination, key}, gconv.Interfaces(keys)...)...))
}

func (r *Redis) SUnion(key string, keys ...string) ([]string, error) {
	return typeStrings(r.commandDo("SUNION", append([]interface{}{key}, gconv.Interfaces(keys)...)...))
}

func (r *Redis) SUnionStore(destination string, key string, keys ...string) (int64, error) {
	return typeInt64(r.commandDo("SUNIONSTORE", append([]interface{}{destination, key}, gconv.Interfaces(keys)...)...))
}

func (r *Redis) SDiff(key string, keys ...string) ([]string, error) {
	return typeStrings(r.commandDo("SDIFF", append([]interface{}{key}, gconv.Interfaces(keys)...)...))
}

func (r *Redis) SDiffStore(destination string, key string, keys ...string) (int64, error) {

	return typeInt64(r.commandDo("SDIFFSTORE", append([]interface{}{destination, key}, gconv.Interfaces(keys)...)...))
}

//======================================================================================zset

func (r *Redis) ZAdd(params ...interface{}) (int, error) {
	return typeInt(r.commandDo("ZADD", params...))
}

func (r *Redis) ZScore(key string, member interface{}) (string, error) {
	return typeString(r.commandDo("ZSCORE", key, member))
}

func (r *Redis) ZinCrBy(key string, increment float64, member interface{}) (string, error) {
	return typeString(r.commandDo("ZINCRBY", key, increment, member))
}

func (r *Redis) ZCard(key string) (int64, error) {
	return typeInt64(r.commandDo("ZCARD", key))
}

func (r *Redis) ZCount(key string, min, max int64) (int64, error) {
	return typeInt64(r.commandDo("ZCOUNT", key, min, max))
}

func (r *Redis) ZRange(key string, start, stop int64, param ...string) ([]string, error) {
	if len(param) == 0 {
		return typeStrings(r.commandDo("ZRANGE", key, start, stop))
	}
	return typeStrings(r.commandDo("ZRANGE", key, start, stop, param[0]))
}

func (r *Redis) ZRevRange(key string, start, stop int64, param ...string) ([]string, error) {
	if len(param) == 0 {
		return typeStrings(r.commandDo("ZRANGE", key, start, stop))
	}
	return typeStrings(r.commandDo("ZREVRANGE", key, start, stop, param[0]))
}

func (r *Redis) ZRangeByScore(key, min, max string, options ...interface{}) ([]string, error) {

	return typeStrings(r.commandDo("ZRANGEBYSCORE", append([]interface{}{key, min, max}, options...)...))
}

func (r *Redis) ZRevRangeByScore(key string, min, max string, options ...interface{}) ([]string, error) {
	return typeStrings(r.commandDo("ZREVRANGEBYSCORE", append([]interface{}{key, min, max}, options...)...))
}

func (r *Redis) ZRank(key, member string) (int64, error) {
	return typeInt64(r.commandDo("ZRANK", key, member))
}

func (r *Redis) ZRevRank(key, member string) (int64, error) {
	return typeInt64(r.commandDo("ZREVRANK", key, member))
}

func (r *Redis) ZRem(key string, member ...interface{}) (int, error) {
	if len(member) == 0 {
		return 0, errors.New("must have an one key")
	}
	return typeInt(r.commandDo("ZREM", append([]interface{}{key}, member...)...))
}

func (r *Redis) ZRemRangeByRank(key string, start, stop int64) (int64, error) {
	return typeInt64(r.commandDo("ZREMRANGEBYRANK", key, start, stop))
}

func (r *Redis) ZRemRangeByScore(key string, min, max float64) (int64, error) {
	return typeInt64(r.commandDo("ZREMRANGEBYSCORE", key, min, max))
}

func (r *Redis) ZRangeByLex(key, min, max string, options ...interface{}) ([]string, error) {
	return typeStrings(r.commandDo("ZRANGEBYLEX", append([]interface{}{key, min, max}, options...)...))
}

func (r *Redis) ZLexCount(key, min, max string) (int64, error) {
	return typeInt64(r.commandDo("ZLEXCOUNT", key, min, max))
}

func (r *Redis) ZRemRangeByLex(key, min, max string) (int64, error) {
	return typeInt64(r.commandDo("ZREMRANGEBYLEX", key, min, max))
}

func (r *Redis) ZUnionStore(options ...interface{}) (int64, error) {
	if len(options) < 3 {
		return 0, errors.New("there must be three parameters")
	}
	return typeInt64(r.commandDo("ZUNIONSTORE", options...))
}
func (r *Redis) ZInterStore(options ...interface{}) (int64, error) {
	return typeInt64(r.commandDo("ZINTERSTORE", options...))
}

//================================================================HyperLogLog
func (r *Redis) PfAdd(key string, options ...interface{}) (int, error) {
	return typeInt(r.commandDo("PFADD", append([]interface{}{key}, options...)...))
}

func (r *Redis) PfCount(keys ...string) (int64, error) {
	return typeInt64(r.commandDo("PFCOUNT", gconv.Interfaces(keys)...))
}

func (r *Redis) PfMerge(keys ...string) (string, error) {
	if len(keys) < 2 {
		return "", errors.New("need at least two keys")
	}
	return typeString(r.commandDo("PFMERGE", gconv.Interfaces(keys)...))
}

//================================================================================GEO
func (r *Redis) GeoAdd(key string, params ...interface{}) (int, error) {
	return typeInt(r.commandDo("GEOADD", append([]interface{}{key}, params...)...))
}

func (r *Redis) GeoPos(key string, member ...interface{}) ([]*GeoLocation, error) {
	return typeGeoLocation(r.commandDo("GEOPOS", append([]interface{}{key}, member...)...))
}

func (r *Redis) GeoDist(key string, params ...string) (string, error) {
	return typeString(r.commandDo("GEODIST", append([]interface{}{key}, gconv.Interfaces(params)...)...))
}

func (r *Redis) GeoRadius(key string, member ...interface{}) ([]*GeoLocation, error) {
	if len(member) < 5 {
		return nil, errors.New("there are must have five keys")
	}
	return typeGeoLocationd(r.commandDo("GEORADIUS", append([]interface{}{key}, member...)...))
}

func (r *Redis) GeoRadiusByMember(key string, member ...interface{}) ([]*GeoLocation, error) {
	return typeGeoLocationd(r.commandDo("GEORADIUSBYMEMBER", append([]interface{}{key}, member...)...))
}

func (r *Redis) GeoHash(key string, member ...interface{}) ([]string, error) {
	return typeStrings(r.commandDo("GEOHASH", append([]interface{}{key}, member...)...))
}

//============================================================================channel
func (r *Redis) Publish(channel, message string) (int, error) {

	return typeInt(r.commandDo("PUBLISH", channel, message))
}

func (r *Redis) PubSub(channel string, member ...interface{}) ([]string, error) {
	return typeStrings(r.commandDo("PUBSUB", append([]interface{}{channel}, member...)...))
}

func (r *Redis) SubScribe(channel ...string) ([]string, error) {
	return typeStrings(r.commandDo("SUBSCRIBE", gconv.Interfaces(channel)...))
}

func (r *Redis) PSubScribe(pattern ...string) ([]string, error) {
	return typeStrings(r.commandDo("PSUBSCRIBE", gconv.Interfaces(pattern)...))
}

func (r *Redis) UnSubScribe(pattern ...string) ([]string, error) {
	return typeStrings(r.commandDo("UNSUBSCRIBE", gconv.Interfaces(pattern)...))
}

func (r *Redis) PunSubScribe(pattern ...string) ([]string, error) {
	return typeStrings(r.commandDo("PUNSUBSCRIBE", gconv.Interfaces(pattern)...))
}
