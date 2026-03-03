package dao

import (
	"AIGO/internal/model"
	"AIGO/pkg/db"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// bson.D 是 MongoDB 驱动程序中用于构建查询过滤器的文档类型。
// 它允许以键值对的形式构建查询条件，与 bson.M 不同，bson.D 保持了键值对的顺序。
// 这在需要保持查询条件的顺序时非常有用，例如在构建复杂的查询时。

// CreateMessage 创建一条消息
func CreateMessage(ctx context.Context, msg *model.Message) (int64, error) {
	res, err := db.MongoDB.InsertOne(ctx, msg)
	if err != nil {
		return 0, fmt.Errorf("insert message failed: %v", err)
	}
	if _, ok := res.InsertedID.(string); ok {
		return 1, nil
	}
	return 0, fmt.Errorf("insert message failed, expected string, got: %T", res.InsertedID)
}

// GetMessagesBySessionID 根据会话ID获取消息历史记录（按创建时间升序排序）
func GetMessagesBySessionID(ctx context.Context, sessionID string) ([]*model.Message, error) {
	filter := bson.D{
		{Key: "session_id", Value: sessionID},
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	msgs := make([]*model.Message, 0)
	cur, err := db.MongoDB.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find messages failed: %v", err)
	}
	if err := cur.All(ctx, &msgs); err != nil {
		return nil, fmt.Errorf("decode messages failed: %v", err)
	}
	return msgs, nil
}

// GetMessagesBySessionIDs 根据会话ID列表获取消息历史记录（按创建时间升序排序）
func GetMessagesBySessionIDs(ctx context.Context, sessionIDs []string) ([]*model.Message, error) {
	filter := bson.D{
		{Key: "session_id", Value: bson.D{
			{Key: "$in", Value: sessionIDs},
		}},
	}
	msgs := make([]*model.Message, 0)
	if len(sessionIDs) == 0 {
		return msgs, fmt.Errorf("sessionIDs is empty")
	}
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	cur, err := db.MongoDB.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find messages failed: %v", err)
	}
	if err := cur.All(ctx, &msgs); err != nil {
		return nil, fmt.Errorf("decode messages failed: %v", err)
	}
	return msgs, nil
}

// GetAllMessages 获取所有消息（按创建时间升序排序）
func GetAllMessages(ctx context.Context) ([]*model.Message, error) {
	// 按创建时间升序排序
	filter := bson.D{
		{Key: "created_at", Value: 1},
	}
	opts := options.Find().SetSort(filter)
	msgs := make([]*model.Message, 0)
	cur, err := db.MongoDB.Find(ctx, bson.D{}, opts)
	if err != nil {
		return nil, fmt.Errorf("find messages failed: %v", err)
	}
	if err := cur.All(ctx, &msgs); err != nil {
		return nil, fmt.Errorf("decode messages failed: %v", err)
	}
	return msgs, nil
}
