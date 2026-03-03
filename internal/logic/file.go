package logic

import (
	"AIGO/internal/dao"
	"AIGO/pkg/log"
	"AIGO/pkg/utils/validate"
	"context"
	"fmt"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"
)

// SaveUpLoadFiles 保存用户上传的文件 不保存到本地 直接保存到milvus中减少io耗时提升接口响应速度
func SaveUpLoadFiles(ctx context.Context, username, sessionID string, fileHeaders []*multipart.FileHeader) ([]string, error) {
	fileIds := make([]string, 0, len(fileHeaders))
	// 1. 检验文件的合法性
	for _, fileHeader := range fileHeaders {
		if err := validate.CheckValidateFile(fileHeader); err != nil {
			return nil, err
		}
		// 定义唯一id 文件名+时间戳+扩展名 防止出现相同文件名但是内容不同的情况
		ext := filepath.Ext(fileHeader.Filename)
		name := strings.TrimSuffix(fileHeader.Filename, ext)
		flieId := fmt.Sprintf("%s_%d%s", name, time.Now().Unix(), ext)
		fileIds = append(fileIds, flieId)
	}
	go func(ctx context.Context, username, sessionID string, fileHeaders []*multipart.FileHeader, fileIds []string) {
		// 创建索引  一个用户一个会话一个索引器 避免重复创建索引
		indexer, err := dao.GetOrCreateRagIndexer(ctx, username, sessionID)
		if err != nil {
			log.Errorf("create rag indexer failed: %v", err)
			return
		}
		// 索引文件
		for idx, fileHeader := range fileHeaders {
			fileId := fileIds[idx]
			// 闭包处理每个文件 确保资源及时释放
			func() {
				// 打开文件流
				file, err := fileHeader.Open()
				if err != nil {
					log.Errorf("open file %s failed: %v", fileId, err)
					// 索引失败 继续索引下一个文件
					return
				}
				defer file.Close()
				// 索引文件
				if err := indexer.IndexFile(ctx, fileId, file); err != nil {
					log.Errorf("index file %s failed: %v", fileId, err)
				}
			}()
		}
	}(context.Background(), username, sessionID, fileHeaders, fileIds) // 异步索引文件 上下文独立 如果传入则必然是gin上下文那么后台没索引完时gin直接返回响应就会中断
	return fileIds, nil
}
