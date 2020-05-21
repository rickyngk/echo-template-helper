package echor

import (
	"encoding/json"
	"fmt"
	"time"
	"github.com/go-redis/redis/v7"
)

// RedisRepoMetaInterface interface
type RedisRepoMetaInterface interface {
	DatasourceID() string
}

func getDriver(conf RedisRepoMetaInterface) *redis.Client {
	driver := getDatasourceDriver(conf.DatasourceID(), "redis")
	client := driver.(*redis.Client)
	return client
}

// // Expire func
// func (o *RedisCache) Expire(key string, expired int64) error {
// 	var dur = time.Duration(expired) * time.Minute
// 	return o.Client.Expire(key, dur).Err()
// }

// // Exist Check if key is exist
// func (o *RedisCache) Exist(key string) bool {
// 	n, _ := o.Client.Exists(key).Result()
// 	return n > 0
// }

// RedisFindKey func
func RedisFindKey(conf RedisRepoMetaInterface, keyPattern string) ([]string, error) {
	return getDriver(conf).Keys(keyPattern).Result()
}

// RedisSet func
func RedisSet(conf RedisRepoMetaInterface, key string, value string, expired int64) error {
	var dur = time.Duration(expired) * time.Minute
	err := getDriver(conf).Set(key, value, dur).Err()
	if err != nil {
		fmt.Println("[Redis SET ERROR]", key, value)
	}
	return err
}

// RedisSetObj func
func RedisSetObj(conf RedisRepoMetaInterface, key string, value interface{}, expired int64) error {
	v, jsonError := json.Marshal(value)
	if jsonError != nil {
		fmt.Println("[Warning]", "SetObj() fail to Marshal data", jsonError)
	} else {
		return RedisSet(conf, key, string(v), expired)
	}
	return jsonError
}

// RedisGet func
func RedisGet(conf RedisRepoMetaInterface, key string) (string, error) {
	return getDriver(conf).Get(key).Result()
}

// RedisGetObj func
func RedisGetObj(conf RedisRepoMetaInterface, key string, cachedData interface{}) bool {
	cached, err := RedisGet(conf, key)
	validCached := false
	if err == nil && cached != "" {
		tmp := []byte(cached)
		err2 := json.Unmarshal(tmp, cachedData)
		validCached = err2 == nil
	}
	return validCached
}

// // CreateHashIfNotExists func
// func (o *RedisCache) CreateHashIfNotExists(key string, expired int64) error {
// 	v, e := o.HGet(key, "__dur")
// 	if e != nil && e != redis.Nil {
// 		return e
// 	}
// 	if v == "" || e == redis.Nil {
// 		e = o.HSet(key, "__dur", strconv.FormatInt(expired, 10))
// 		if e != nil {
// 			return e
// 		}
// 	}
// 	e = o.Expire(key, expired)
// 	return e
// }

// // HSet func
// func (o *RedisCache) HSet(key string, field string, value string) error {
// 	err := o.Client.HSet(key, field, value).Err()
// 	if err != nil {
// 		fmt.Println("[Redis HSET ERROR]", key, value)
// 	}
// 	return err
// }

// // HGet func
// func (o *RedisCache) HGet(key string, field string) (string, error) {
// 	return o.Client.HGet(key, field).Result()
// }

// // HGetSafe func
// func (o *RedisCache) HGetSafe(key string, field string) (string, error) {
// 	r, e := o.Client.HGet(key, field).Result()
// 	if e == redis.Nil {
// 		return "", nil
// 	}
// 	return r, e
// }

// // Push func
// func (o *RedisCache) Push(key string, value string) (string, error) {
// 	err := o.Client.LPush(key, value).Err()
// 	if err != nil {
// 		fmt.Println("[Redis LPUSH ERROR]", key, value)
// 	}
// 	return "", err
// }

// // PushObj func
// func (o *RedisCache) PushObj(key string, value interface{}) error {
// 	v, jsonError := json.Marshal(value)
// 	if jsonError != nil {
// 		fmt.Println("[Warning]", "PushObj() fail to Marshal data", jsonError)
// 	} else {
// 		o.Push(key, string(v))
// 	}
// 	return jsonError
// }

// RedisRemove key
func RedisRemove(conf RedisRepoMetaInterface, key string) error {
	_, e := getDriver(conf).Del(key).Result()
	return e
}
