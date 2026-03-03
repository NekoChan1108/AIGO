package db

import (
	"AIGO/internal/config"
	"AIGO/pkg/log"
	"time"

	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var MongoDB *mongo.Collection

// getMongoDBConn 获取MongoDB连接
func getMongoDBConn() (*mongo.Client, error) {
	if config.Cfg.MongoCfg == nil {
		log.Fatalf("MongoDB config is nil")
	}
	//客户端选项
	opts := &options.ClientOptions{}
	if config.Cfg.MongoCfg.Hosts != nil {
		opts.SetHosts(config.Cfg.MongoCfg.Hosts)
	}
	if config.Cfg.MongoCfg.User != "" && config.Cfg.MongoCfg.Password != "" {
		opts.SetAuth(options.Credential{
			Username: config.Cfg.MongoCfg.User,
			Password: config.Cfg.MongoCfg.Password,
		})
	}
	// 连接池设置
	// 校验并设置最大连接池数
	if config.Cfg.MongoCfg.MaxPoolSize > 0 {
		opts.SetMaxPoolSize(config.Cfg.MongoCfg.MaxPoolSize)
	}
	// 校验并设置最小连接池数
	if config.Cfg.MongoCfg.MinPoolSize > 0 {
		opts.SetMinPoolSize(config.Cfg.MongoCfg.MinPoolSize)
	}
	// 校验并设置最大空闲时间
	if config.Cfg.MongoCfg.MaxConnIdleTime > 0 {
		opts.SetMaxConnIdleTime(time.Duration(config.Cfg.MongoCfg.MaxConnIdleTime) * time.Minute)
	}
	// 校验并设置连接超时时间
	if config.Cfg.MongoCfg.ConnectTimeout > 0 {
		opts.SetConnectTimeout(time.Duration(config.Cfg.MongoCfg.ConnectTimeout) * time.Second)
	}
	// 校验并设置单个连接池最大连接数
	if config.Cfg.MongoCfg.MaxConnecting > 0 {
		opts.SetMaxConnecting(config.Cfg.MongoCfg.MaxConnecting)
	}
	return mongo.Connect(opts)
}

func init() {
	client, err := getMongoDBConn()
	if err != nil {
		log.Fatalf("Failed to connect to mongodb: %v", err)
	}
	MongoDB = client.Database(config.Cfg.MongoCfg.DataBase).Collection(config.Cfg.MongoCfg.Collection)
}
