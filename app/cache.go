package app

import (
	"github.com/MeetFootball/Scaffold/util"
	"github.com/gomodule/redigo/redis"
)

// Cache 缓存
type Cache struct {
	Conf *RedisConfig // 配置
	Conn redis.Conn   // Redis 连接
}

// NewCache 实例化
func NewCache(name string) (*Cache, error) {
	var (
		err   error
		Cache = &Cache{}
		pool  *redis.Pool
	)
	path := util.GetConfigPath("redis")
	if Cache.Conf, err = GetRedisConfig(path, name); err != nil {
		return Cache, err
	}
	if pool, err = GetRedisPool(name); err != nil {
		return Cache, err
	}
	Cache.Conn = pool.Get()
	_, err = Cache.Conn.Do("SELECT", Cache.Conf.Db)
	return Cache, err
}

// GenerateRK 生成 Redis Key
func (c *Cache) GenerateRK(key string) string {
	return c.Conf.Prefix + key
}

// Exists 查看缓存是否存在
func (c *Cache) Exists(key string) (bool, error) {
	key = c.GenerateRK(key)
	return redis.Bool(c.Conn.Do("EXISTS", key))
}

// Get 获取数据
func (c *Cache) Get(key string) ([]byte, error) {
	key = c.GenerateRK(key)
	return redis.Bytes(c.Conn.Do("GET", key))
}

// GetString 获取数据(字符串类型)
func (c *Cache) GetString(key string) (string, error) {
	key = c.GenerateRK(key)
	return redis.String(c.Conn.Do("GET", key))
}

// Set 存储 String 数据，有过期时间
func (c *Cache) SetEX(key string, value []byte) (err error) {
	key = c.GenerateRK(key)
	_, err = redis.String(c.Conn.Do("SET", key, value, "EX", c.Conf.Expired))
	return
}

// RPush 从右侧推入数据
func (c *Cache) RPush(key string, value []byte) (err error) {
	key = c.GenerateRK(key)
	_, err = c.Conn.Do("RPUSH", key, value)
	return
}

// LPop 从左侧取出数据
func (c *Cache) LPop(key string) (value []byte, err error) {
	key = c.GenerateRK(key)
	value, err = redis.Bytes(c.Conn.Do("LPOP", key))
	return
}

// Del 删除缓存
func (c *Cache) Del(key string) (bool, error) {
	key = c.GenerateRK(key)
	return redis.Bool(c.Conn.Do("DEL", key))
}

// Close 关闭连接
func (c *Cache) Close() {
	if err := c.Conn.Close(); err != nil {
		WriteErrorLog("errors", "Console", "redis", err.Error())
	}
}

// CleanCacheByKey 根据 键 删除缓存
func CleanCacheByKey(conn, key string) (err error) {
	var Cache *Cache
	if Cache, err = NewCache(conn); err != nil {
		return err
	}
	defer Cache.Close() // 关闭连接池占用
	if _, err = Cache.Del(key); err != nil {
		return err
	}
	return
}

// CleanCacheByKeys 删除 多个缓存
func CleanCacheByKeys(conn string, keys []string) (err error) {
	var cache *Cache
	if cache, err = NewCache(conn); err != nil {
		return err
	}
	defer cache.Close() // 关闭连接池占用
	for _, key := range keys {
		_, err = cache.Del(key) // 删除缓存
		if err != nil {
			return err
		} // 返回错误
	}
	return
}
