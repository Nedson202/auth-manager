package auth_service_rest

import (
	"fmt"

	"github.com/gomodule/redigo/redis"
)

func (app App) handleCacheGet(key string) (value string, err error) {
	conn := app.cachePool.Get()
	defer conn.Close()

	value, err = redis.String(conn.Do("GET", key))
	if err == redis.ErrNil {
		return
	} else if err != nil {
		return "", fmt.Errorf("Redis:::Error getting key %s: %v", key, err)
	}

	return
}

func (app App) handleCacheSet(key string, value interface{}) (err error) {
	conn := app.cachePool.Get()
	defer conn.Close()

	_, err = conn.Do("SET", key, value)
	if err != nil {
		return fmt.Errorf("Redis:::Error setting key %s to %s: %v", key, value, err)
	}
	return
}
