package logic

import (
	"AIGO/internal/dao"
	"AIGO/internal/model"
	"AIGO/pkg/log"
	"context"
	"encoding/json"
	"fmt"

	"github.com/segmentio/kafka-go"
)

// Consume 消费kafka消息并写入mongodb
func Consume(ctx context.Context, consumer *kafka.Reader) error {
	var msg model.Message
	for {
		kMsg, err := consumer.ReadMessage(ctx)
		if err != nil {
			return fmt.Errorf("read message failed: %w", err)
		}
		if err := json.Unmarshal(kMsg.Value, &msg); err != nil {
			return fmt.Errorf("unmarshal message failed: %w", err)
		}
		log.Info(msg.Content)
		// 写入mongodb
		if _, err := dao.CreateMessage(ctx, &msg); err != nil {
			return fmt.Errorf("insert message failed: %w", err)
		}
	}
}
