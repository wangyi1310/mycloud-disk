package cache

import (
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/wangyi1310/mycloud-disk/pkg/log"
)

// RedisStore redis存储驱动
type RedisStore struct {
	pool *redis.Pool
}

// NewRedisStore 创建新的redis存储
func NewRedisStore(size int, network, address, user, password, database string) *RedisStore {
	return &RedisStore{
		pool: &redis.Pool{
			MaxIdle:     size,
			IdleTimeout: 240 * time.Second,
			TestOnBorrow: func(c redis.Conn, t time.Time) error {
				_, err := c.Do("PING")
				return err
			},
			Dial: func() (redis.Conn, error) {
				db, err := strconv.Atoi(database)
				if err != nil {
					return nil, err
				}

				c, err := redis.Dial(
					network,
					address,
					redis.DialDatabase(db),
					redis.DialUsername(user),
					redis.DialPassword(password),
				)
				if err != nil {
					log.Log().Panic("Failed to create Redis connection: %s", err)
				}
				return c, nil
			},
		},
	}
}

// Set 存储值
func (store *RedisStore) Set(key string, value interface{}, ttl int) error {
	rc := store.pool.Get()
	defer rc.Close()

	serialized, err := serializer(value)
	if err != nil {
		return err
	}

	if rc.Err() != nil {
		return rc.Err()
	}

	if ttl > 0 {
		_, err = rc.Do("SETEX", key, ttl, serialized)
	} else {
		_, err = rc.Do("SET", key, serialized)
	}

	if err != nil {
		return err
	}
	return nil

}

// Get 取值
func (store *RedisStore) Get(key string) (interface{}, bool) {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return nil, false
	}

	v, err := redis.Bytes(rc.Do("GET", key))
	if err != nil || v == nil {
		return nil, false
	}

	finalValue, err := deserializer(v)
	if err != nil {
		return nil, false
	}

	return finalValue, true

}

// Gets 批量取值
func (store *RedisStore) Gets(keys []string, prefix string) (map[string]interface{}, []string) {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return nil, keys
	}

	var queryKeys = make([]string, len(keys))
	for key, value := range keys {
		queryKeys[key] = prefix + value
	}

	v, err := redis.ByteSlices(rc.Do("MGET", redis.Args{}.AddFlat(queryKeys)...))
	if err != nil {
		return nil, keys
	}

	var res = make(map[string]interface{})
	var missed = make([]string, 0, len(keys))

	for key, value := range v {
		decoded, err := deserializer(value)
		if err != nil || decoded == nil {
			missed = append(missed, keys[key])
		} else {
			res[keys[key]] = decoded
		}
	}
	// 解码所得值
	return res, missed
}

// Sets 批量设置值
func (store *RedisStore) Sets(values map[string]interface{}, prefix string) error {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return rc.Err()
	}
	var setValues = make(map[string]interface{})

	// 编码待设置值
	for key, value := range values {
		serialized, err := serializer(value)
		if err != nil {
			return err
		}
		setValues[prefix+key] = serialized
	}

	_, err := rc.Do("MSET", redis.Args{}.AddFlat(setValues)...)
	if err != nil {
		return err
	}
	return nil

}

// Delete 批量删除给定的键
func (store *RedisStore) Delete(keys []string, prefix string) error {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return rc.Err()
	}

	// 处理前缀
	for i := 0; i < len(keys); i++ {
		keys[i] = prefix + keys[i]
	}

	_, err := rc.Do("DEL", redis.Args{}.AddFlat(keys)...)
	if err != nil {
		return err
	}
	return nil
}

// DeleteAll 批量所有键
func (store *RedisStore) DeleteAll() error {
	rc := store.pool.Get()
	defer rc.Close()
	if rc.Err() != nil {
		return rc.Err()
	}

	_, err := rc.Do("FLUSHDB")

	return err
}

// Persist Dummy implementation
func (store *RedisStore) Persist(path string) error {
	return nil
}

// Restore dummy implementation
func (store *RedisStore) Restore(path string) error {
	return nil
}
