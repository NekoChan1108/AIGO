package dao

import (
	"AIGO/internal/model"
	"AIGO/pkg/db"
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

// GetSessionByID 根据会话ID获取会话
func GetSessionByID(ctx context.Context, sessionID string) (*model.Session, error) {
	// 获取session
	session := &model.Session{}
	res := db.MongoDB.FindOne(ctx, bson.M{"_id": sessionID}).Decode(session)
	if len(res.Error()) > 0 {
		return nil, fmt.Errorf("get session by id failed: %v", res.Error())
	}
	return session, nil
}

// CreateSession 创建会话
func CreateSession(ctx context.Context, session *model.Session) (int64, error) {
	res, err := db.MongoDB.InsertOne(ctx, session)
	if err != nil {
		return 0, fmt.Errorf("insert session failed: %v", err)
	}
	if _, ok := res.InsertedID.(string); ok {
		return 1, nil
	}
	return 0, fmt.Errorf("insert session failed, expected string, got: %T", res.InsertedID)
}

// GetSessionByUsername 根据用户名获取会话
//
//	func GetSessionByUsername(ctx context.Context, username string) ([]model.Session, error) {
//		sessions := make([]model.Session, 0)
//		res, err := db.MongoDB.Find(ctx, bson.D{{Key: "username", Value: username}})
//		if err != nil {
//			return nil, fmt.Errorf("find sessions failed: %v", err)
//		}
//		if err := res.All(ctx, &sessions); err != nil {
//			return nil, fmt.Errorf("decode sessions failed: %v", err)
//		}
//		return sessions, nil
//	}
//
// GetSessionByUsername 根据用户名获取会话
func GetSessionByUsername(ctx context.Context, username string) ([]model.Session, error) {
	sessions := make([]model.Session, 0)
	// 过滤条件：匹配用户名 且 必须包含title字段(用于区分Session和Message)
	filter := bson.D{
		{Key: "username", Value: username},
		{Key: "title", Value: bson.D{{Key: "$exists", Value: true}}},
	}
	// 按创建时间倒序排序
	opts := options.Find().SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := db.MongoDB.Find(ctx, filter, opts)
	if err != nil {
		return nil, fmt.Errorf("find sessions failed: %v", err)
	}
	defer cursor.Close(ctx)

	if err := cursor.All(ctx, &sessions); err != nil {
		return nil, fmt.Errorf("decode sessions failed: %v", err)
	}

	return sessions, nil
}

// DeleteSession 删除会话
func DeleteSession(ctx context.Context, sessionID string) (int64, error) {
	res, err := db.MongoDB.DeleteOne(ctx, bson.M{"_id": sessionID})
	if err != nil || res.DeletedCount <= 0 {
		return 0, fmt.Errorf("delete session failed: %v", err)
	}
	return res.DeletedCount, nil
}
