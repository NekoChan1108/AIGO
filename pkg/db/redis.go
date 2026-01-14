package db

import (
	"AIGO/internal/config"
	"AIGO/pkg/log"
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisDB 全局redis数据库连接
var RedisDB *redis.Client

// getRedisConn 获取redis连接
func getRedisConn() *redis.Client {
	if config.Cfg.RedisCfg == nil {
		log.Fatalf("Redis config is nil")
	}
	addr := fmt.Sprintf("%s:%s", config.Cfg.RedisCfg.Host, config.Cfg.RedisCfg.Port)
	opts := &redis.Options{
		Addr: addr,
	}
	// 校验并设置数据库
	if config.Cfg.RedisCfg.DataBase >= 0 {
		opts.DB = config.Cfg.RedisCfg.DataBase
	}
	// 校验并设置用户名
	if config.Cfg.RedisCfg.User != "" {
		opts.Username = config.Cfg.RedisCfg.User
	}
	// 校验并设置密码
	if config.Cfg.RedisCfg.Password != "" {
		opts.Password = config.Cfg.RedisCfg.Password
	}
	// 校验并设置超时时间
	if config.Cfg.RedisCfg.Timeout > 0 {
		opts.DialTimeout = time.Duration(config.Cfg.RedisCfg.Timeout) * time.Second
	}
	// 校验并设置连接池
	if config.Cfg.RedisCfg.PoolSize > 0 {
		opts.PoolSize = config.Cfg.RedisCfg.PoolSize
	}
	// 校验并设置连接空闲时间
	if config.Cfg.RedisCfg.ConnMaxIdleTime > 0 {
		opts.ConnMaxIdleTime = time.Duration(config.Cfg.RedisCfg.ConnMaxIdleTime) * time.Minute
	}
	// 校验并设置连接最大生命周期
	if config.Cfg.RedisCfg.ConnMaxLifeTime > 0 {
		opts.ConnMaxLifetime = time.Duration(config.Cfg.RedisCfg.ConnMaxLifeTime) * time.Minute
	}
	// 校验并设置最大打开连接数
	if config.Cfg.RedisCfg.MaxActiveConns > 0 {
		opts.MaxActiveConns = config.Cfg.RedisCfg.MaxActiveConns
	}
	// 校验并设置最大空闲连接数
	if config.Cfg.RedisCfg.MaxIdleConns > 0 {
		opts.MaxIdleConns = config.Cfg.RedisCfg.MaxIdleConns
	}
	// 校验并设置最小空闲连接数
	if config.Cfg.RedisCfg.MinIdleConns > 0 {
		opts.MinIdleConns = config.Cfg.RedisCfg.MinIdleConns
	}
	return redis.NewClient(opts)

}

func init() {
	if err := getRedisConn().Ping(context.Background()).Err(); err != nil {
		log.Fatalf("Failed to connect to redis: %v", err)
	}
	RedisDB = getRedisConn()
}
