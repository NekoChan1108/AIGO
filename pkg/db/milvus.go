package db

import (
	"AIGO/internal/config"
	"AIGO/pkg/log"
	"context"
	"fmt"
	"time"

	"github.com/milvus-io/milvus/client/v2/milvusclient"
)

// MilvusDB milvus数据库连接
var MilvusDB *milvusclient.Client

// getMilvusConn 获取milvus连接
func getMilvusConn() (*milvusclient.Client, error) {
	ctx := context.Background()
	if config.Cfg.MilvusCfg == nil {
		log.Fatal("Milvus config is nil")
	}
	addr := fmt.Sprintf("%s:%s", config.Cfg.MilvusCfg.Host, config.Cfg.MilvusCfg.Port)
	cfg := &milvusclient.ClientConfig{
		Address: addr,
	}
	// 校验并设置用户名和密码
	if config.Cfg.MilvusCfg.User != "" && config.Cfg.MilvusCfg.Password != "" {
		cfg.Username = config.Cfg.MilvusCfg.User
		cfg.Password = config.Cfg.MilvusCfg.Password
	}
	// 校验并设置数据库
	if config.Cfg.MilvusCfg.DataBase != "" {
		cfg.DBName = config.Cfg.MilvusCfg.DataBase
	}
	// 校验并设置超时时间
	if config.Cfg.MilvusCfg.Timeout > 0 {
		timeCtx, cancel := context.WithTimeout(ctx, time.Duration(config.Cfg.MilvusCfg.Timeout)*time.Second)
		defer cancel()
		return milvusclient.New(timeCtx, cfg)
	}
	return milvusclient.New(ctx, cfg)
}

func init() {
	client, err := getMilvusConn()
	if err != nil {
		log.Fatalf("Failed to connect to milvus: %v", err)
	}
	MilvusDB = client
}
