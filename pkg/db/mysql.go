package db

import (
	"AIGO/internal/config"
	"AIGO/internal/model"
	"AIGO/pkg/log"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

// MysqlDB mysql数据库连接 全局配置
var MysqlDB *gorm.DB

// getMysqlConn 获取数据库连接
func getMysqlConn() (*gorm.DB, error) {
	if config.Cfg.MysqlCfg == nil {
		log.Fatalf("Mysql config is nil")
	}
	dsn := fmt.Sprintf("%v:%v@tcp(%v:%v)/%v?charset=utf8mb4&parseTime=true&loc=Local",
		config.Cfg.MysqlCfg.User, config.Cfg.MysqlCfg.Password,
		config.Cfg.MysqlCfg.Host, config.Cfg.MysqlCfg.Port,
		config.Cfg.MysqlCfg.DataBase)
	return gorm.Open(mysql.New(
		mysql.Config{
			DSN:                       dsn,
			SkipInitializeWithVersion: false,
		},
	), &gorm.Config{})
}

func init() {
	if db, err := getMysqlConn(); err != nil {
		log.Fatalf("Failed to connect to mysql: %v", err)
	} else {
		sqlDB, err := db.DB()
		if err != nil {
			log.Fatalf("Failed to get mysql db: %v", err)
		}
		if config.Cfg.MysqlCfg.ConnMaxLifeTime > 0 {
			sqlDB.SetConnMaxLifetime(time.Duration(config.Cfg.MysqlCfg.ConnMaxLifeTime) * time.Minute) // 设置连接最大生命周期
		}
		if config.Cfg.MysqlCfg.MaxOpenConns > 0 {
			sqlDB.SetMaxOpenConns(config.Cfg.MysqlCfg.MaxOpenConns) // 设置最大打开连接数
		}
		if config.Cfg.MysqlCfg.MaxIdleConns > 0 {
			sqlDB.SetMaxIdleConns(config.Cfg.MysqlCfg.MaxIdleConns) // 设置最大空闲连接数
		}
		if config.Cfg.MysqlCfg.ConnMaxIdleTime > 0 {
			sqlDB.SetConnMaxIdleTime(time.Duration(config.Cfg.MysqlCfg.ConnMaxIdleTime) * time.Minute) // 设置连接最大空闲时间
		}
		MysqlDB = db
		if err := db.AutoMigrate(&model.User{}); err != nil {
			log.Fatalf("Failed to auto migrate mysql: %v", err)
		}
	}
}
