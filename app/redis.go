package app

import (
	"errors"
	"fmt"
	"time"

	"github.com/MeetFootball/Scaffold/model"
	"github.com/MeetFootball/Scaffold/util"
	"github.com/gomodule/redigo/redis"
)

// ConfigRedisMap 全局变量
var ConfigRedisMap *model.RedisMapConfig
var RedisPool map[string]*redis.Pool

// RedisMapConfig Redis 配置集合
type RedisMapConfig struct {
	List map[string]*RedisConfig `mapstructure:"list"`
}

// RedisConfig Redis 配置
type RedisConfig struct {
	ProxyList    string `mapstructure:"proxy_list"`
	Password     string `mapstructure:"password"`
	Prefix       string `mapstructure:"prefix"`
	Db           int    `mapstructure:"db"`
	Expired      int    `mapstructure:"expired"`
	MaxIdle      int    `mapstructure:"max_idle"`
	MaxActive    int    `mapstructure:"max_active"`
	ConnTimeout  int    `mapstructure:"conn_timeout"`
	IdelTimeout  int    `mapstructure:"idle_timeout"`
	ReadTimeout  int    `mapstructure:"read_timeout"`
	WriteTimeout int    `mapstructure:"write_timeout"`
}

// InitRedisConfig 加载 Redis 配置
func InitRedisConfig(path string) error {
	RedisConfigMap := &model.RedisMapConfig{}
	err := util.ParseConfig(path, RedisConfigMap)
	if err != nil {
		return err
	}
	if len(RedisConfigMap.List) == 0 {
		fmt.Printf("[INFO] %s%s\n", time.Now().Format(util.DateTimeFormat), " empty redis config.")
	}
	RedisPool = map[string]*redis.Pool{}
	for configName, config := range RedisConfigMap.List {
		dialector := &redis.Pool{
			MaxIdle:         config.MaxIdle,   // 最大空闲连接数
			MaxActive:       config.MaxActive, // 分配的最大连接数
			MaxConnLifetime: time.Duration(config.ConnTimeout),
			IdleTimeout:     time.Duration(config.IdelTimeout),
			Wait:            true,
			Dial: func() (redis.Conn, error) {
				options := redis.DialDatabase(config.Db)
				password := redis.DialPassword(config.Password)
				// **重要** 设置读写超时
				readTimeout := redis.DialReadTimeout(time.Second * time.Duration(config.ReadTimeout))
				writeTimeout := redis.DialReadTimeout(time.Second * time.Duration(config.WriteTimeout))
				conTimeout := redis.DialConnectTimeout(time.Second * time.Duration(config.ConnTimeout))
				c, err := redis.Dial("tcp", config.ProxyList, options, password, readTimeout, writeTimeout, conTimeout)
				if err != nil {
					panic(err.Error())
				}
				return c, err
			},
		}
		RedisPool[configName] = dialector
	}
	return nil
}

// GetRedisPool 获取 Redis 数据库连接
func GetRedisPool(name string) (*redis.Pool, error) {
	if pool, ok := RedisPool[name]; ok {
		return pool, nil
	}
	return nil, errors.New("GetRedisPoolError") // 获取 Redis 连接池错误
}

// GetRedisConfig 获取 Redis 配置
func GetRedisConfig(path, name string) (conf *RedisConfig, err error) {
	conf = &RedisConfig{}
	list := &RedisMapConfig{}
	// 解析 Redis 配置
	if err = util.ParseConfig(path, list); err != nil {
		err = errors.New("ReadRedisConfigError : " + path)
		return nil, err
	}
	for n, v := range list.List {
		if n == name {
			conf.Db = v.Db
			conf.Prefix = v.Prefix
			conf.Expired = v.Expired
		}
	}
	return
}

// CloseRedisDB 关闭 Redis 数据库
func CloseRedisDB() error {
	for _, pool := range RedisPool {
		err := pool.Close()
		if err != nil {
			return err
		}
	}
	return nil
}
