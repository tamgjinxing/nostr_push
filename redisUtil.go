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

// "address": "52.76.210.159:6379",
// "password": "cc6ee5619d1b",
// "db": 0

// GetRedisClient 获取 Redis 客户端连接实例
func GetRedisClient() (*redis.Client, error) {

	var err error
	once.Do(func() {
		redisAddr = config.Redis.Address
		redisPwd = config.Redis.Password
		redisDBNum = 0

		// 创建 Redis 客户端连接
		redisCli = redis.NewClient(&redis.Options{
			Addr:     redisAddr,
			Password: redisPwd,
			DB:       redisDBNum,
		})

		// 检查是否能够连接到 Redis
		_, err = redisCli.Ping(ctx).Result()
		if err != nil {
			fmt.Printf("failed to connect to Redis: %v\n", err)
			return
		}
	})

	return redisCli, err
}

// Set 设置键值对
func SetRedisKey(key, value string) error {
	// 获取 Redis 客户端连接实例
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	// 设置键值对
	return redisCli.Set(ctx, key, value, 0).Err()
}

// Get 获取键值对
func GetRedisByKey(key string) (string, error) {
	// 获取 Redis 客户端连接实例
	redisCli, err := GetRedisClient()
	if err != nil {
		return "", err
	}

	// 获取键值对
	return redisCli.Get(ctx, key).Result()
}

// Del 删除键
func DelRedisByKey(key string) error {
	// 获取 Redis 客户端连接实例
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	// 删除键
	return redisCli.Del(ctx, key).Err()
}

// Close 关闭 Redis 客户端连接
func Close() error {
	// 获取 Redis 客户端连接实例
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	// 关闭 Redis 客户端连接
	return redisCli.Close()
}

func SetRedisKeyByList(key string, values ...interface{}) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	// 设置键值对
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

// putJSON 将结构体序列化为 JSON 并存储到 Redis 中
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

// putHashList 将列表序列化为 JSON 并存储到 Redis 哈希表的特定字段中
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

// getHashList 从 Redis 哈希表的特定字段中获取 JSON 并反序列化为列表
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

// SaveRecentKey 保存最近 100 个键
func SaveRecentKey(key string) error {
	redisCli, err := GetRedisClient()
	if err != nil {
		return err
	}

	// 使用事务保证原子性
	_, err = redisCli.TxPipelined(ctx, func(pipe redis.Pipeliner) error {
		pipe.LPush(ctx, "recent_keys", key)   // 将新键推入列表头部
		pipe.LTrim(ctx, "recent_keys", 0, 99) // 保留列表中的前 100 个键
		return nil
	})
	return err
}

// 判断这个键最近是否有保存过
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

// 模糊删除以 prefix 开头的 Redis 键
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

		// if len(keys) > 0 {
		// 	for _, key := range keys {
		// 		log.Printf("Key:%s", key)
		// 	}
		// }

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
	fmt.Printf("Deleted %d keys with prefix %s\n", n, prefix)
	return nil
}
