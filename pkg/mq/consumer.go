package mq

import (
	"AIGO/config"

	"github.com/segmentio/kafka-go"
)

var KafkaConsumer *kafka.Reader

func kafkaConsumer() *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers:     config.Cfg.KafkaCfg.Brokers,
		Topic:       "test-topic",     // TODO 服务器创建并写入yaml配置
		GroupID:     "test-group",     // TODO 服务器创建并写入yaml配置
		StartOffset: kafka.LastOffset, // 从最新消息开始消费
	})
}

func init() {
	KafkaConsumer = kafkaConsumer()
}
