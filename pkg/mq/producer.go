package mq

import (
	"AIGO/config"

	"github.com/segmentio/kafka-go"
)

var KafkaProducer *kafka.Writer

func kafkaProducer() *kafka.Writer {
	return &kafka.Writer{
		Addr:                   kafka.TCP(config.Cfg.KafkaCfg.Brokers...),
		Topic:                  "test-topic",  // TODO 服务器创建并写入yaml配置
		Balancer:               &kafka.Hash{}, // 消息有序
		RequiredAcks:           kafka.RequireNone,
		AllowAutoTopicCreation: true,      // 允许自动创建 topic
		Compression:            kafka.Lz4, // 压缩算法
	}
}

func init() {
	KafkaProducer = kafkaProducer()
}
