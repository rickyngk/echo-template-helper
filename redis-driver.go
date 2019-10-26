package echor

import (
	"fmt"

	"github.com/go-redis/redis"
)

// RedisDriver class
type RedisDriver struct {
	connector *redis.Client
}

// NewRedisDriver - create new connection
func NewRedisDriver(conf DataSourceSpecs) DbDriverInterface {

	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%d", conf.Address, conf.Port),
		Password: conf.Password,
		DB:       conf.DbIndex,
	})
	driver := &RedisDriver{connector: client}
	return driver
}

// GetNativeConnector func - impl interface DbNativeInterface
func (d *RedisDriver) GetNativeConnector() interface{} {
	return d.connector
}
