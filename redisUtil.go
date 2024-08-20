package main

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
)

const (
	// UserKeyPrefix 表示用户键的前缀
	UserKeyPrefix = "user:"

	// SessionKeyPrefix 表示会话键的前缀
	SessionKeyPrefix = "session:"

	// CounterKey 表示计数器键
	CounterKey = "counter"
)

var (
	once       sync.Once
	redisCli   *redis.Client
	ctx        = context.Background()
	redisAddr  string
	redisPwd   string
	redisDBNum int
)

func GetRedisClient() (*redis.Client, error) {

	var err error
	once.Do(func() {
		redisAddr = config.Redis.Address
		redisPwd = config.Redis.Password
		redisDBNum = 0

		redisCli = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPwd,
			DB:       redisDBNum,
		})

		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			fmt.Printf("failed to connect to Redis: %v\n", err)
			return
		}
	})

	return redisCli, err
}

func SetRedisKey(key, value string) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	return redisCli.Set(ctx, key, value, 0).Err()
}

func GetRedisByKey(key string) (string, error) {
	redisCli, err := GetRedisClient()
	if err != nil {
		return "", err
	}

	return redisCli.Get(ctx, key).Result()
}

func DelRedisByKey(key string) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	return redisCli.Del(ctx, key).Err()
}

func Close() error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	return redisCli.Close()
}

func SetRedisKeyByList(key string, values ...interface{}) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	return redisCli.RPush(ctx, key, values).Err()
}

func GetRedisListByKey(key string) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	return redisCli.LRange(ctx, key, 0, -1).Err()
}

func PopByKey(key string) (string, error) {
	redisCli, err := GetRedisClient()
	if err != nil {
		return "", err
	}

	val, err := redisCli.RPop(ctx, key).Result()

	if err != nil {
		return "", err
	}

	return val, nil
}

// HGetJSON gets a hash field and unmarshals the JSON encoded value
func HGetJSON(key, field string, dest interface{}) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}
	data, err := redisCli.HGet(ctx, key, field).Result()
	if err != nil {
		return err
	}
	return json.Unmarshal([]byte(data), dest)
}

func GetJSON(key string, dest interface{}) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}
	data, err := redisCli.Get(ctx, key).Result()
	if err != nil {
		return err
	}

	return json.Unmarshal([]byte(data), dest)
}

func PutJSON(key string, value interface{}) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	data, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = redisCli.Set(ctx, key, data, 0).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func PutHashList(key, field string, list []string) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	data, err := json.Marshal(list)
	if err != nil {
		return err
	}

	err = redisCli.HSet(ctx, key, field, data).Err()
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

func GetHashList(key, field string) ([]string, error) {
	redisCli, err := GetRedisClient()
	if err != nil {
		return nil, err
	}

	data, err := redisCli.HGet(ctx, key, field).Result()
	if err != nil {
		return nil, err
	}

	var list []string
	err = json.Unmarshal([]byte(data), &list)
	if err != nil {
		return nil, err
	}

	return list, nil
}

func SaveRecentKey(key string) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	_, err = redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.LPush(ctx, "recent_keys", key)
		pipe.LTrim(ctx, "recent_keys", 0, 99)
		return nil
	})
	return err
}

func RecentKeysExist(key string) (bool, error) {
	redisCli, err := GetRedisClient()
	if err != nil {
		return false, err
	}

	keys, err := redisCli.LRange(ctx, "recent_keys", 0, 99).Result()

	if err != nil {
		return false, nil
	}

	for _, value := range keys {
		if key == value {
			return true, nil
		}
	}

	return false, nil
}

func DeleteKeysWithPrefix(prefix string) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	var cursor uint64
	var n int
	for {
		var keys []string
		var err error
		keys, cursor, err = redisCli.Scan(ctx, cursor, prefix+"*", 100).Result()
		if err != nil {
			return err
		}

		if len(keys) > 0 {
			if err := redisCli.Del(ctx, keys...).Err(); err != nil {
				return err
			}
		}
		n += len(keys)

		if cursor == 0 {
			break
		}
	}
	log.Printf("Deleted %d keys with prefix %s\n", n, prefix)
	return nil
}
