package config

import (
	"AIGO/pkg/log"
	"path/filepath"
	"runtime"

	"github.com/spf13/viper"
)

// TODO: 添加到配置中心例如etcd nacos consul 云端读不到读本地配置进行配置初始化

// config 配置
type config struct {
	AppCfg    *appConfig     `mapstructure:"App"`     // 应用配置
	MysqlCfg  *mysqlConfig   `mapstructure:"Mysql"`   // mysql 配置
	RedisCfg  *redisConfig   `mapstructure:"Redis"`   // redis 配置
	KafkaCfg  *kafkaConfig   `mapstructure:"Kafka"`   // kafka 配置
	ModelCfg  *modelConfig   `mapstructure:"Model"`   // 模型配置
	MilvusCfg *milvusConfig  `mapstructure:"Milvus"`  // milvus 配置
	EmailCfg  *emailConfig   `mapstructure:"Email"`   // 邮件配置
	MongoCfg  *mongodbConfig `mapstructure:"MongoDB"` // mongodb 配置
	JwtCfg    *jwtConfig     `mapstructure:"Jwt"`     // jwt 配置
}

// appConfig 应用配置
type appConfig struct {
	Name        string `mapstructure:"Name"`        // 应用名称
	Version     string `mapstructure:"Version"`     // 应用版本
	Port        string `mapstructure:"Port"`        // 应用端口
	Path        string `mapstructure:"Path"`        // 应用路径
	MaxFileSize int64  `mapstructure:"MaxFileSize"` // 最大文件大小
}

// mysqlConfig mysql 配置
type mysqlConfig struct {
	Host            string `mapstructure:"Host"`            // 数据库地址
	Port            string `mapstructure:"Port"`            // 数据库端口
	User            string `mapstructure:"User"`            // 数据库用户名
	Password        string `mapstructure:"Password"`        // 数据库密码
	DataBase        string `mapstructure:"DataBase"`        // 数据库名称
	Charset         string `mapstructure:"Charset"`         // 数据库字符集
	ParseTime       bool   `mapstructure:"ParseTime"`       // 解析时间
	Location        string `mapstructure:"Location"`        // 时区
	MaxOpenConns    int    `mapstructure:"MaxOpenConns"`    // 最大连接数
	MaxIdleConns    int    `mapstructure:"MaxIdleConns"`    // 最大空闲连接数
	ConnMaxIdleTime int    `mapstructure:"ConnMaxIdleTime"` // 连接空闲超时 mins
	ConnMaxLifeTime int    `mapstructure:"ConnMaxLifeTime"` // 连接超时 mins
}

// redisConfig redis 配置
type redisConfig struct {
	Host            string `mapstructure:"Host"`            // redis 地址
	Port            string `mapstructure:"Port"`            // redis 端口
	DataBase        int    `mapstructure:"DataBase"`        // redis 数据库
	User            string `mapstructure:"User"`            // redis 用户名
	Password        string `mapstructure:"Password"`        // redis 密码
	Timeout         int    `mapstructure:"Timeout"`         // redis 连接超时时间 seconds
	PoolSize        int    `mapstructure:"PoolSize"`        // redis 连接池大小
	MaxActiveConns  int    `mapstructure:"MaxActiveConns"`  // 最大连接数
	MaxIdleConns    int    `mapstructure:"MaxIdleConns"`    // 最大空闲连接数
	MinIdleConns    int    `mapstructure:"MinIdleConns"`    // 最小空闲连接数
	ConnMaxIdleTime int    `mapstructure:"ConnMaxIdleTime"` // 连接空闲超时 mins
	ConnMaxLifeTime int    `mapstructure:"ConnMaxLifeTime"` // 连接超时 mins
}

// kafkaConfig kafka 配置
type kafkaConfig struct {
	Brokers []string `mapstructure:"Brokers"` // kafka 地址
}

// modelConfig 模型配置
type modelConfig struct {
	ApiKey         string `mapstructure:"ApiKey"`         // 火山模型api key
	TextModel      string `mapstructure:"TextModel"`      // 文本模型
	EmbeddingModel string `mapstructure:"EmbeddingModel"` // 向量模型
	VoiceModel     string `mapstructure:"VoiceModel"`     // 语音模型
	VisionModel    string `mapstructure:"VisionModel"`    // 视觉模型
}

// milvusConfig milvus 配置
type milvusConfig struct {
	Host       string `mapstructure:"Host"`       // milvus 地址
	Port       string `mapstructure:"Port"`       // milvus 端口
	DataBase   string `mapstructure:"DataBase"`   // milvus 数据库
	Collection string `mapstructure:"Collection"` // milvus 集合
	User       string `mapstructure:"User"`       // milvus 用户名
	Password   string `mapstructure:"Password"`   // milvus 密码
	Dimension  int64  `mapstructure:"Dimension"`  // 向量维度 需要和EmbeddingModel模型一致
	TOPK       int    `mapstructure:"TOPK"`       // RAG 召回时选择最接近的答案数
	Timeout    int    `mapstructure:"Timeout"`    // 连接超时
}

// emailConfig 邮件配置
type emailConfig struct {
	Host           string `mapstructure:"Host"`           // 邮箱smtp地址
	Port           int    `mapstructure:"Port"`           // 邮箱smtp端口
	Sender         string `mapstructure:"Sender"`         // 发送者(发送邮件的邮箱)
	Authentication string `mapstructure:"AuthentiCation"` // 邮箱密码
	Expiration     int    `mapstructure:"Expiration"`     // 验证码有效时间
}

// mongodbConfig mongodb 配置
type mongodbConfig struct {
	Hosts           []string `mapstructure:"Hosts"`           // mongodb 集群地址
	User            string   `mapstructure:"User"`            // mongodb 用户名
	Password        string   `mapstructure:"Password"`        // mongodb 密码
	DataBase        string   `mapstructure:"DataBase"`        // mongodb 数据库
	Collection      string   `mapstructure:"Collection"`      // mongodb 集合
	MaxPoolSize     uint64   `mapstructure:"MaxPoolSize"`     // 配置连接池最大数
	MinPoolSize     uint64   `mapstructure:"MinPoolSize"`     // 配置连接池最小数
	MaxConnIdleTime int      `mapstructure:"MaxConnIdleTime"` // 配置连接最大空闲时间 mins
	MaxConnecting   uint64   `mapstructure:"MaxConnecting"`   // 配置单个连接池最大连接数
	ConnectTimeout  int      `mapstructure:"ConnectTimeout"`  // 配置连接mongodb超时时间 seconds
}

// jwtConfig jwt 配置
type jwtConfig struct {
	Secret            string `mapstructure:"Secret"`            // jwt 密钥
	Issuer            string `mapstructure:"Issuer"`            // jwt 签发者
	Subject           string `mapstructure:"Subject"`           // jwt 主题
	AccessExpiration  int64  `mapstructure:"AccessExpiration"`  // jwt 访问令牌过期时间
	RefreshExpiration int64  `mapstructure:"RefreshExpiration"` // jwt 刷新令牌过期时间
}

var (
	v   = viper.New() // viper 配置读取器
	Cfg = new(config) // 全局配置
)

// loadConfig 加载配置
func loadConfig() {
	_, path, _, ok := runtime.Caller(0)
	if !ok {
		log.Fatalf("Runtime Error: Failed to get config file path")
	}
	// 加载yaml文件
	viper.AddConfigPath(filepath.Dir(path)) // 告诉去这个路径下找
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	// 读取配置
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	// 默认值 TODO 设置更多的默认值
	v.SetDefault("App.Port", ":9999")
	v.SetDefault("App.Path", "/api/v1")

	// 将配置映射到结构体
	if err := viper.Unmarshal(Cfg); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}
}

func init() {
	loadConfig()
}
